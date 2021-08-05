package api

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/docdb"
	log "github.com/sirupsen/logrus"
)

// documentDBCreate creates documentDB cluster and instances
func (o *docDBOrchestrator) documentDBCreate(ctx context.Context, req *DocDBCreateRequest) (*DocDBResponse, error) {
	log.Infof("creating documentDB cluster %s with %d instance(s)", aws.StringValue(req.DBClusterIdentifier), aws.IntValue(req.InstanceCount))

	req.Tags = normalizeTags(o.org, req.Tags)

	// check if a DBSubnetGroup exists, and create it if needed
	dbSubnetGroupFound, err := o.dbSubnetGroupExists(ctx, dbSubnetGroupName(o.org))
	if err != nil {
		return nil, err
	} else if !dbSubnetGroupFound {
		if err := o.dbSubnetGroupCreate(ctx, dbSubnetGroupName(o.org), req.SubnetIds); err != nil {
			return nil, err
		}
	} else {
		log.Infof("subnet group %s already exists, will use it for this docdb cluster", dbSubnetGroupName(o.org))
	}

	cluster, err := o.client.CreateDBCluster(ctx, &docdb.CreateDBClusterInput{
		BackupRetentionPeriod: req.BackupRetentionPeriod,
		DBClusterIdentifier:   req.DBClusterIdentifier,
		DBSubnetGroupName:     aws.String(dbSubnetGroupName(o.org)),
		Engine:                aws.String("docdb"),
		EngineVersion:         req.EngineVersion,
		MasterUsername:        req.MasterUsername,
		MasterUserPassword:    req.MasterUserPassword,
		StorageEncrypted:      aws.Bool(true),
		Tags:                  toDocDBTags(req.Tags),
		VpcSecurityGroupIds:   req.VpcSecurityGroupIds,
	})
	if err != nil {
		return nil, err
	}

	allDBInstances := []*docdb.DBInstance{}
	for i := 1; i <= aws.IntValue(req.InstanceCount); i++ {
		instanceName := fmt.Sprintf("%s-%d", aws.StringValue(req.DBClusterIdentifier), i)

		dbInstance, err := o.client.CreateDBInstance(ctx, &docdb.CreateDBInstanceInput{
			AutoMinorVersionUpgrade: aws.Bool(true),
			DBInstanceClass:         req.DBInstanceClass,
			DBClusterIdentifier:     req.DBClusterIdentifier,
			DBInstanceIdentifier:    aws.String(instanceName),
			Engine:                  aws.String("docdb"),
			Tags:                    toDocDBTags(req.Tags),
		})
		if err != nil {
			// TODO: Rollback
			return nil, apierror.New(apierror.ErrBadRequest, "failed to create docdb instance "+instanceName, err)
		}

		allDBInstances = append(allDBInstances, dbInstance)
	}

	return &DocDBResponse{
		Cluster:   cluster,
		Instances: allDBInstances,
	}, nil
}

// documentDBList lists all documentDB clusters
func (o *docDBOrchestrator) documentDBList(ctx context.Context) ([]string, error) {
	output, err := o.client.ListDocDBClusters(ctx)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// documentDBDetails returns details about a documentDB cluster
func (o *docDBOrchestrator) documentDBDetails(ctx context.Context, name string) (*docdb.DBCluster, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	output, err := o.client.GetDocDB(ctx, name)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == docdb.ErrCodeDBClusterNotFoundFault {
				return nil, apierror.New(apierror.ErrNotFound, "cluster not found", nil)
			}
		}
		return nil, err
	}

	return output, nil
}

// documentDBDelete deletes documentDB cluster and associated instances
func (o *docDBOrchestrator) documentDBDelete(ctx context.Context, name string, snapshot bool) error {
	if name == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	documentDB, err := o.client.GetDocDB(ctx, name)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == docdb.ErrCodeDBClusterNotFoundFault {
				return apierror.New(apierror.ErrNotFound, "cluster not found", nil)
			}
		}
		return err
	}

	log.Infof("deleting documentDB cluster %s (snapshot: %t)", name, snapshot)

	// first loop through all the cluster instances and delete them
	for _, i := range documentDB.DBClusterMembers {
		_, err := o.client.DeleteDBInstance(ctx, &docdb.DeleteDBInstanceInput{
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

	_, err = o.client.DeleteDBCluster(ctx, &input)
	if err != nil {
		return err
	}

	return nil
}

// dbSubnetGroupName determines the DBSubnetGroup based on the Org
func dbSubnetGroupName(org string) string {
	return fmt.Sprintf("spinup-%s-docdb-subnetgroup", org)
}

// dbSubnetGroupExists checks if a DBSubnetGroup exists
func (o *docDBOrchestrator) dbSubnetGroupExists(ctx context.Context, name string) (bool, error) {
	result, err := o.client.GetDBSubnetGroup(ctx, name)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == docdb.ErrCodeDBSubnetGroupNotFoundFault {
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
		return false, apierror.New(apierror.ErrInternalError, "unexpected number of matching subnet groups", nil)
	}
}

// dbSubnetGroupCreate creates a DBSubnetGroup
func (o *docDBOrchestrator) dbSubnetGroupCreate(ctx context.Context, name string, subnets []*string) error {
	if subnets == nil {
		return apierror.New(apierror.ErrBadRequest, "no subnets specified", nil)
	}

	log.Infof("creating DBSubnetGroup %s with subnets: %v", name, aws.StringValueSlice(subnets))

	_, err := o.client.CreateDBSubnetGroup(ctx, &docdb.CreateDBSubnetGroupInput{
		DBSubnetGroupDescription: aws.String(name),
		DBSubnetGroupName:        aws.String(name),
		SubnetIds:                subnets,
	})
	if err != nil {
		return apierror.New(apierror.ErrBadRequest, "failed to create DBSubnetGroup", err)
	}

	return nil
}
