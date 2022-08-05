package api

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/flywheel"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// documentDBCreate creates documentDB cluster and instances
func (o *docDBOrchestrator) documentDBCreate(ctx context.Context, req *DocDBCreateRequest) (*DocDBResponse, *flywheel.Task, error) {
	log.Infof("creating documentDB cluster %s with %d instance(s)", aws.StringValue(req.DBClusterIdentifier), aws.IntValue(req.InstanceCount))

	req.Tags = req.Tags.normalize(o.server.org)

	sgName := dbSubnetGroupName(o.server.org, req.SubnetIds)

	// check if a DBSubnetGroup exists, and create it if needed
	dbSubnetGroupFound, err := o.dbSubnetGroupExists(ctx, sgName)
	if err != nil {
		return nil, nil, err
	}

	if !dbSubnetGroupFound {
		if err := o.dbSubnetGroupCreate(ctx, sgName, req.SubnetIds); err != nil {
			return nil, nil, err
		}
	} else {
		log.Infof("subnet group %s already exists, will use it for this docdb cluster", sgName)
	}

	task := flywheel.NewTask()

	cluster, err := o.docdbClient.CreateDBCluster(ctx, &docdb.CreateDBClusterInput{
		BackupRetentionPeriod: req.BackupRetentionPeriod,
		DBClusterIdentifier:   req.DBClusterIdentifier,
		DBSubnetGroupName:     aws.String(sgName),
		Engine:                aws.String("docdb"),
		EngineVersion:         req.EngineVersion,
		MasterUsername:        req.MasterUsername,
		MasterUserPassword:    req.MasterUserPassword,
		StorageEncrypted:      aws.Bool(true),
		Tags:                  req.Tags.toDocDBTags(),
		VpcSecurityGroupIds:   req.VpcSecurityGroupIds,
	})
	if err != nil {
		return nil, nil, err
	}

	allDBInstances := []*docdb.DBInstance{}
	for i := 1; i <= aws.IntValue(req.InstanceCount); i++ {
		instanceName := fmt.Sprintf("%s-%d", aws.StringValue(req.DBClusterIdentifier), i)

		dbInstance, err := o.docdbClient.CreateDBInstance(ctx, &docdb.CreateDBInstanceInput{
			AutoMinorVersionUpgrade: aws.Bool(true),
			DBInstanceClass:         req.DBInstanceClass,
			DBClusterIdentifier:     req.DBClusterIdentifier,
			DBInstanceIdentifier:    aws.String(instanceName),
			Engine:                  aws.String("docdb"),
			Tags:                    req.Tags.toDocDBTags(),
		})
		if err != nil {
			// TODO: Rollback
			return nil, nil, err
		}

		allDBInstances = append(allDBInstances, dbInstance)
	}

	// start the async orchestration to wait for docdb cluster to become available
	go func() {
		cl := aws.StringValue(req.DBClusterIdentifier)

		taskCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		msgChan, errChan := o.startTask(taskCtx, task)

		msgChan <- fmt.Sprintf("requested creation of docdb cluster %s", cl)

		if err = retry(10, 3, 10*time.Second, func() error {
			msgChan <- fmt.Sprintf("checking if docdb cluster %s is available before continuing", cl)

			if err := o.refreshSession(taskCtx); err != nil {
				msgChan <- fmt.Sprintf("unable to refresh orchestrator session: %s", err)
				return err
			}

			// check cluster status
			cluster, err := o.docdbClient.GetDocDBDetails(taskCtx, cl)
			if err != nil {
				msgChan <- fmt.Sprintf("got error checking if docdb cluster %s is available: %s", cl, err)
				return err
			}

			if status := aws.StringValue(cluster.Status); status != "available" {
				msgChan <- fmt.Sprintf("docdb cluster %s is not yet available (%s)", cl, status)
				return fmt.Errorf("docdb cluster %s not yet available", cl)
			}

			// check instances
			instances, err := o.docdbClient.GetDocDBInstances(taskCtx, cl)
			if err != nil {
				msgChan <- fmt.Sprintf("got error describing docdb instances for %s: %s", cl, err)
				return err
			}

			if len(instances) == 0 {
				msgChan <- fmt.Sprintf("docdb cluster %s doesn't have any instances", cl)
				return fmt.Errorf("docdb cluster %s has no instances", cl)
			}

			for _, i := range instances {
				if status := aws.StringValue(i.DBInstanceStatus); status != "available" {
					msgChan <- fmt.Sprintf("not all docdb instances in cluster %s are available", cl)
					return fmt.Errorf("not all docdb instances in cluster %s are available", cl)
				}
			}

			msgChan <- fmt.Sprintf("docdb cluster %s is available", cl)
			return nil
		}); err != nil {
			errChan <- fmt.Errorf("failed to create docdb cluster %s, timeout waiting to become available: %s", cl, err.Error())
			return
		}
	}()

	return &DocDBResponse{
		Cluster:   cluster,
		Instances: allDBInstances,
	}, task, nil
}

// documentDBList lists all documentDB clusters
func (o *docDBOrchestrator) documentDBList(ctx context.Context) ([]string, error) {
	out, err := o.rgClient.GetResourcesInOrg(ctx, o.server.org, "database", "docdb")
	if err != nil {
		return nil, err
	}

	resources := make([]string, 0, len(out))
	for _, r := range out {
		a, err := arn.Parse(aws.StringValue(r.ResourceARN))
		if err != nil {
			return nil, apierror.New(apierror.ErrInternalError, "failed to parse ARN "+aws.StringValue(r.ResourceARN), err)
		}

		parts := strings.SplitN(a.Resource, ":", 2)
		if !strings.HasPrefix(parts[1], "cluster-") {
			// AWS DocumentDB creates 2 ARNs for each cluster: one with the name and one with a unique DbClusterResourceId
			// that we are excluding here (it looks like cluster-L3R4YRSBUYDP4GLMTJ2WF5GH5Q)
			resources = append(resources, parts[1])
		}
	}

	return resources, nil
}

// documentDBDetails returns details about a documentDB cluster
func (o *docDBOrchestrator) documentDBDetails(ctx context.Context, name string) (*DocDBResponse, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	cluster, err := o.docdbClient.GetDocDBDetails(ctx, name)
	if err != nil {
		return nil, err
	}

	t, err := o.docdbClient.GetDocDBTags(ctx, cluster.DBClusterArn)
	if err != nil {
		return nil, err
	}
	tags := fromDocDBTags(t)

	if !tags.inOrg(o.server.org) {
		return nil, apierror.New(apierror.ErrNotFound, "cluster not found in our org", nil)
	}

	return &DocDBResponse{
		Cluster: cluster,
		Tags:    tags,
	}, nil
}

// documentDBModify modifies documentDB cluster and instances
func (o *docDBOrchestrator) documentDBModify(ctx context.Context, name string, req *DocDBModifyRequest) (*DocDBResponse, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	documentDB, err := o.docdbClient.GetDocDBDetails(ctx, name)
	if err != nil {
		return nil, err
	}

	// modify cluster parameters
	cluster, err := o.docdbClient.ModifyDBCluster(ctx, &docdb.ModifyDBClusterInput{
		ApplyImmediately:       aws.Bool(true),
		BackupRetentionPeriod:  req.BackupRetentionPeriod,
		DBClusterIdentifier:    aws.String(name),
		EngineVersion:          req.EngineVersion,
		MasterUserPassword:     req.MasterUserPassword,
		NewDBClusterIdentifier: req.NewDBClusterIdentifier,
		VpcSecurityGroupIds:    req.VpcSecurityGroupIds,
	})
	if err != nil {
		return nil, err
	}

	allDBInstances := []*docdb.DBInstance{}

	// if needed, loop through all the cluster instances and modify them
	if req.DBInstanceClass != nil {
		for _, i := range documentDB.DBClusterMembers {
			dbInstance, err := o.docdbClient.ModifyDBInstance(ctx, &docdb.ModifyDBInstanceInput{
				ApplyImmediately:     aws.Bool(true),
				DBInstanceIdentifier: i.DBInstanceIdentifier,
				DBInstanceClass:      req.DBInstanceClass,
			})
			if err != nil {
				return nil, err
			}

			allDBInstances = append(allDBInstances, dbInstance)
		}
	}

	return &DocDBResponse{
		Cluster:   cluster,
		Instances: allDBInstances,
	}, nil
}

// documentDBDelete deletes documentDB cluster and associated instances
func (o *docDBOrchestrator) documentDBDelete(ctx context.Context, name string, snapshot bool) error {
	if name == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	documentDB, err := o.docdbClient.GetDocDBDetails(ctx, name)
	if err != nil {
		return err
	}

	log.Infof("deleting documentDB cluster %s (snapshot: %t)", name, snapshot)

	// first loop through all the cluster instances and delete them
	for _, i := range documentDB.DBClusterMembers {
		_, err := o.docdbClient.DeleteDBInstance(ctx, &docdb.DeleteDBInstanceInput{
			DBInstanceIdentifier: i.DBInstanceIdentifier,
		})
		if err != nil {
			return err
		}
	}

	input := docdb.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(name),
		SkipFinalSnapshot:   aws.Bool(true),
	}

	if snapshot {
		input.SkipFinalSnapshot = aws.Bool(false)
		input.FinalDBSnapshotIdentifier = aws.String("final-" + name)
	}

	if _, err = o.docdbClient.DeleteDBCluster(ctx, &input); err != nil {
		return err
	}

	return nil
}

// dbSubnetGroupName determines the DBSubnetGroup name based on the Org and subnet id's
// generate a 32-char MD5 hash based on all subnet id's
func dbSubnetGroupName(org string, subnetIds []string) string {
	return fmt.Sprintf("spinup-%s-docdb-sg-%x", org, md5.Sum([]byte(strings.Join(subnetIds, ""))))
}

// dbSubnetGroupExists checks if a DBSubnetGroup exists
func (o *docDBOrchestrator) dbSubnetGroupExists(ctx context.Context, name string) (bool, error) {
	result, err := o.docdbClient.GetDBSubnetGroup(ctx, name)
	if err != nil {
		if aerr, ok := errors.Cause(err).(apierror.Error); ok {
			if aerr.Code == apierror.ErrNotFound {
				log.Debugf("subnet group not found: %s", name)
				return false, nil
			} else {
				return false, err
			}
		}
	}

	if len(result) == 1 {
		return true, nil
	} else {
		return false, apierror.New(apierror.ErrInternalError, "unexpected number of matching subnet groups: "+fmt.Sprint(len(result)), nil)
	}
}

// dbSubnetGroupCreate creates a DBSubnetGroup
func (o *docDBOrchestrator) dbSubnetGroupCreate(ctx context.Context, name string, subnets []string) error {
	if subnets == nil {
		return apierror.New(apierror.ErrBadRequest, "no subnets specified", nil)
	}

	log.Infof("creating DBSubnetGroup %s with subnets: %v", name, subnets)

	_, err := o.docdbClient.CreateDBSubnetGroup(ctx, &docdb.CreateDBSubnetGroupInput{
		DBSubnetGroupDescription: aws.String(name),
		DBSubnetGroupName:        aws.String(name),
		SubnetIds:                aws.StringSlice(subnets),
	})
	if err != nil {
		return apierror.New(apierror.ErrBadRequest, "failed to create DBSubnetGroup", err)
	}

	return nil
}

// startTask starts the flywheel task and receives messages on the channels.  in the future, this
// functionality might be part of the flywheel library
func (o *docDBOrchestrator) startTask(ctx context.Context, task *flywheel.Task) (chan<- string, chan<- error) {
	msgChan := make(chan string)
	errChan := make(chan error)

	// track the task
	go func() {
		taskCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := o.server.flywheel.Start(taskCtx, task); err != nil {
			log.Errorf("failed to start flywheel task, won't be tracked: %s", err)
		}

		for {
			select {
			case msg := <-msgChan:
				log.Infof("task %s: %s", task.ID, msg)

				if ferr := o.server.flywheel.CheckIn(taskCtx, task.ID); ferr != nil {
					log.Errorf("failed to checkin task %s: %s", task.ID, ferr)
				}

				if ferr := o.server.flywheel.Log(taskCtx, task.ID, msg); ferr != nil {
					log.Errorf("failed to log flywheel message for %s: %s", task.ID, ferr)
				}
			case err := <-errChan:
				log.Error(err)

				if ferr := o.server.flywheel.Fail(taskCtx, task.ID, err.Error()); ferr != nil {
					log.Errorf("failed to fail flywheel task %s: %s", task.ID, ferr)
				}

				return
			case <-ctx.Done():
				log.Infof("marking task %s complete", task.ID)

				if ferr := o.server.flywheel.Complete(taskCtx, task.ID); ferr != nil {
					log.Errorf("failed to complete flywheel task %s: %s", task.ID, ferr)
				}

				return
			}
		}
	}()

	return msgChan, errChan
}

// docDBState is used to start and stop a given instance
func (o *docDBOrchestrator) docDBState(ctx context.Context, state string, name string) error {
	if state == "" || name == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	state = strings.ToLower(state)
	switch state {
	case "start":
		return o.docdbClient.StartDBCluster(ctx, name)
	case "stop":
		return o.docdbClient.StopDBCluster(ctx, name)
	default:
		msg := fmt.Sprintf("unknown power state %q", state)
		return apierror.New(apierror.ErrBadRequest, msg, nil)
	}
}
