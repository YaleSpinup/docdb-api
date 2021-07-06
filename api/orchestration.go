package api

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	log "github.com/sirupsen/logrus"
)

var DefaultAvailabilityZones = []*string{
	aws.String("us-east-1a"),
	aws.String("us-east-1d"),
}

var SubnetIDSlice = []*string{
	aws.String("subnet-0707d40ddbb9d0818"),
	aws.String("subnet-02364e7e9fc4d8045"),
}

/*
// subnetAzs returns a map of subnets to availability zone names used by EFS
// this may change to support getting the list of subnets as well, currently it uses
// the defaults from the EFS service
func (o *docDBOrchestrator) subnetAzs(ctx context.Context, account string) (map[string]string, error) {
	log.Infof("determining availability zone for account %s using efs onezone", account)

	// 		DBSubnetGroupCreateStatus, err := o.client.CreateDBSubnetGroup(ctx, &inputDBSubnetGroup)

	if !ok {
		return nil, apierror.New(apierror.ErrNotFound, "account doesnt exist", nil)
	}

	subnets := make(map[string]string)
	for _, s := range efsService.DefaultSubnets {
		subnet, err := ec2Service.GetSubnet(ctx, s)
		if err != nil {
			return nil, err
		}

		log.Debugf("got details about subnet %s: %+v", s, subnet)

		subnets[s] = aws.StringValue(subnet.AvailabilityZone)
	}


	if len(subnets) == 0 {
		return nil, apierror.New(apierror.ErrBadRequest, "failed to determine usable availability zone", nil)
	}

	return subnets, nil
}
*/

// checkDBSubnetGroup checks if DBsubnetgroup exists
func (o *docDBOrchestrator) checkDBSubnetGroup(ctx context.Context, name string) bool {
	if name == "" {
		return false
	}

	search := &docdb.DescribeDBSubnetGroupsInput{DBSubnetGroupName: aws.String(name)}

	searchDBSubnetGroup, err := o.client.GetDBSubnetGroup(ctx, search)
	if err != nil {
		if searchDBSubnetGroup == nil {
			// check that err returned contains 404
			// failed to create db cluster (DBSubnetGroupNotFoundFault: DB subnet group
			//'<search-name>' does not exist.
			// status code: 404, request id: some request id)
			regex := regexp.MustCompile("404")
			match := regex.Match([]byte(fmt.Sprintln(err)))
			if match == true {
				log.Infof("DBSubnetGroup %s does not exist\n", aws.StringValue(search.DBSubnetGroupName))
				log.Debugf("DBSubnetGroup %s does not exist: %v\n", aws.StringValue(search.DBSubnetGroupName), err)
				return false
			}
		}
	}
	return true
}

// createDBSubnetGroup creates a dbSubnetGroup
func (o *docDBOrchestrator) createDBSubnetGroup(ctx context.Context, name string, tags []*docdb.Tag) (string, error) {
	if name == "" || tags == nil {
		return "", apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating DBSubnetGroup: %s\n", name)

	inputDBSubnetGroup := docdb.CreateDBSubnetGroupInput{
		DBSubnetGroupDescription: aws.String(name),
		DBSubnetGroupName:        aws.String(name),
		SubnetIds:                SubnetIDSlice,
		Tags:                     tags,
	}

	DBSubnetGroupCreateStatus, err := o.client.CreateDBSubnetGroup(ctx, &inputDBSubnetGroup)
	if err != nil {
		return "", apierror.New(apierror.ErrBadRequest, "failed to create DBSubnetGroup", err)
	}
	log.Debugf("DBSubnetGroup create status: %+v\n", DBSubnetGroupCreateStatus)

	return "{OK}", nil

}

// createDocumentDB creates documentDB cluster and instances
func (o *docDBOrchestrator) createDocumentDB(ctx context.Context, name string, data *CreateDocDB) ([]byte, error) {
	if data == nil || name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating documentDB cluster: %s, and instance count: %d\n", data.DBClusterIdentifier, data.InstanceCount)

	// normalize tags
	tags := make([]*docdb.Tag, 0, len(data.Tags))

	for _, t := range data.Tags {
		tags = append(tags, &docdb.Tag{
			Key:   (t.Key),
			Value: (t.Value),
		})
	}

	newTags := []*docdb.Tag{
		{
			Key:   aws.String("spinup:org"),
			Value: aws.String(o.org),
		},
	}

	for _, t := range tags {
		if aws.StringValue(t.Key) != "spinup:org" && aws.StringValue(t.Key) != "yale:org" {
			newTags = append(newTags, t)
		}
	}

	// check if DBSubnetGroup exists, if not create it
	// create a DBSubnetGroup for each AWS account and spinup environment
	// e.g.: spinup-localdev-docdb-subnetgroup, spinup-spinup-docdb-subnetgroup
	dbSubnetGroupName := fmt.Sprintf("spinup-%s-docdb-subnetgroup", o.org)

	foundDbSubnetGroup := o.checkDBSubnetGroup(ctx, dbSubnetGroupName)

	if foundDbSubnetGroup == true {
		log.Infof("found dbSubnetGroupName: %s\n", dbSubnetGroupName)
	} else {
		_, err := o.createDBSubnetGroup(ctx, dbSubnetGroupName, newTags)
		if err != nil {
			log.Debugf("failed to create dbSubnetGroupName: %s, %s \n", dbSubnetGroupName, err)
		}
	}

	input := docdb.CreateDBClusterInput{
		AvailabilityZones:   DefaultAvailabilityZones,
		DBClusterIdentifier: &data.DBClusterIdentifier,
		DBSubnetGroupName:   aws.String(dbSubnetGroupName),
		Engine:              &data.Engine,
		MasterUsername:      &data.MasterUsername,
		MasterUserPassword:  &data.MasterUserPassword,
		StorageEncrypted:    aws.Bool(true),
		Tags:                newTags,
	}

	clusterCreateStatus, err := o.client.CreateDBCluster(ctx, &input)
	if err != nil {
		// fixme - the err doesn't seem to be returning to the caller, to the user, and it
		// probably will be helpful
		return nil, apierror.New(apierror.ErrBadRequest, "failed to create db cluster", err)
	}

	allInstances := []*DBInstance{}

	for i := 1; i <= data.InstanceCount; i++ {
		// normalize instanceName
		instanceName := fmt.Sprintf("%s-%d", data.DBClusterIdentifier, i)

		instanceData := docdb.CreateDBInstanceInput{
			AutoMinorVersionUpgrade: aws.Bool(true),
			// FIXME - use client input here
			AvailabilityZone:           aws.String("us-east-1a"),
			DBInstanceClass:            &data.DBInstanceClass,
			DBClusterIdentifier:        &data.DBClusterIdentifier,
			DBInstanceIdentifier:       aws.String(instanceName),
			Engine:                     &data.Engine,
			PreferredMaintenanceWindow: &data.MaintenanceWindow,
			PromotionTier:              &data.PromotionTier,
		}

		instanceCreateStatus, err := o.client.CreateDBInstance(ctx, &instanceData)
		if err != nil {
			//return []string{"failed to create db instance: "}, err
			return nil, apierror.New(apierror.ErrBadRequest, "failed to create db instances", err)
		}

		loopInstance := *&DBInstance{
			AvailabilityZone:      *instanceCreateStatus.AvailabilityZone,
			BackupRetentionPeriod: *instanceCreateStatus.BackupRetentionPeriod,
			DBInstanceArn:         *instanceCreateStatus.DBInstanceArn,
			DBInstanceClass:       *instanceCreateStatus.DBInstanceClass,
			DBInstanceStatus:      *instanceCreateStatus.DBInstanceStatus,
			DBInstanceIdentifier:  *instanceCreateStatus.DBInstanceIdentifier,
			Engine:                *instanceCreateStatus.Engine,
			EngineVersion:         *instanceCreateStatus.EngineVersion,
			KmsKeyId:              *instanceCreateStatus.KmsKeyId,
		}

		/*
			DBSubnetGroup: &DBSubnetGroup{
				DBSubnetGroupARN: *instanceCreateStatus.DBSubnetGroup.DBSubnetGroupArn,
			},

			Endpoint: Endpoint{
				Port:         *&instanceCreateStatus.Endpoint.Port,
				Address:      *&instanceCreateStatus.Endpoint.Address,
				HostedZoneId: *&instanceCreateStatus.Endpoint.HostedZoneId,
			},

					// DBSubnetGroup lists a DBSubnetGroup configuration
			type DBSubnetGroup struct {
				DBSubnetGroupARN         string
				DBSubnetGroupDescription string
				DBSubnetGroupName        string
				SubnetGroupStatus        string
				Subnets                  []*Subnet
				VpcID                    string
			}
		*/

		// jam these back in with structs to support them
		//DBSubnetGroup:         *instanceCreateStatus.DBSubnetGroup,
		//Endpoint:              Endpoint,

		allInstances = append(allInstances, &loopInstance)

	}

	createOut := Cluster{
		DBClusters: DBCluster{
			DBClusterArn:        *clusterCreateStatus.DBClusterArn,
			DBClusterIdentifier: *clusterCreateStatus.DBClusterIdentifier,
			Endpoint:            *clusterCreateStatus.Endpoint,
			ReaderEndpoint:      *clusterCreateStatus.ReaderEndpoint,
			StorageEncrypted:    *clusterCreateStatus.StorageEncrypted,
			DBSubnetGroup:       *clusterCreateStatus.DBSubnetGroup,
			DBInstances:         allInstances,
		},
	}

	marshaledData, err := json.Marshal(createOut)
	if err != nil {
		return nil, apierror.New(apierror.ErrInternalError, "failed to marshal docdb output", err)
	}

	return marshaledData, nil

}

// listDocumentDB lists all documentDB clusters and instances on account
func (o *docDBOrchestrator) listDocumentDB(ctx context.Context) (*docdb.DescribeDBClustersOutput, error) {

	input := docdb.DescribeDBClustersInput{}

	dbList, err := o.client.ListDB(ctx, &input)
	if err != nil {
		return nil, apierror.New(apierror.ErrBadRequest, "failed to list documentDBs", err)
	}

	return dbList, nil

}

// getDocumentDB gets data on a documentDB cluster+instance
func (o *docDBOrchestrator) getDocumentDB(ctx context.Context, name string) (*docdb.DBCluster, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	input := docdb.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(name),
	}

	documentDB, err := o.client.GetDB(ctx, &input)
	if err != nil {
		return nil, apierror.New(apierror.ErrBadRequest, "failed to get documentDB", err)
	}

	return documentDB, nil
}

// deleteDocumentDB deletes documentDB instances and cluster
func (o *docDBOrchestrator) deleteDocumentDB(ctx context.Context, name string, data *DeleteDocDB) (string, error) {
	if name == "" {
		return "", apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("deleting documentBD instance(s) and cluster: %s\n", data.ClusterName)

	// Get instances in cluster, so we can delete them without user giving input
	getCluster := docdb.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(name),
	}

	documentDB, err := o.client.GetDB(ctx, &getCluster)
	if err != nil {
		return "", apierror.New(apierror.ErrBadRequest, "failed to get documentDB", err)
	}

	log.Infof("getting documentDB cluster info: %s\n", aws.StringValue(documentDB.DBClusterIdentifier))

	// Loop through the DBClusterMember instances and delete them
	for _, iName := range documentDB.DBClusterMembers {
		instanceDeleteInput := docdb.DeleteDBInstanceInput{
			DBInstanceIdentifier: aws.String(*iName.DBInstanceIdentifier),
		}

		_, err := o.client.DeleteDBInstance(ctx, &instanceDeleteInput)
		if err != nil {
			log.Infof("Failed to delete db instance: %s\n", err)
			return "failed to delete db instance: ", err
		}
	}

	input := docdb.DeleteDBClusterInput{
		DBClusterIdentifier:       aws.String(data.ClusterName),
		SkipFinalSnapshot:         &data.SkipFinalSnapshot,
		FinalDBSnapshotIdentifier: &data.FinalDBSnapshotIdentifier,
	}

	_, err = o.client.DeleteDBCluster(ctx, &input)
	if err != nil {
		return "failed to delete db cluster: ", err
	}

	return "{OK}", nil
}
