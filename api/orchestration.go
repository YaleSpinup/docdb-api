package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	log "github.com/sirupsen/logrus"
)

var DefaultAvailabilityZones = []*string{
	aws.String("us-east-1a"),
	aws.String("us-east-1d"),
}

// createDocumentDB creates documentDB cluster and instances
func (o *docDBOrchestrator) createDocumentDB(ctx context.Context, name string, data *CreateDocDB) ([]byte, error) {
	if data == nil || name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("Creating documentDB instances and cluster: %s\n", data.DBClusterIdentifier)

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

	//
	// CreateDBSubnetGroup - spinup-orgname-docdb-subnet
	// Verify it exists, if not, create it
	dbSubnetGroupName := fmt.Sprintf("spinup-%s-docdb-subnetgroup", o.org)

	subnetIDSlice := []*string{aws.String("subnet-0707d40ddbb9d0818"), aws.String("subnet-02364e7e9fc4d8045")}
	//subnetIDSlice := []*string{aws.String("subnet-0707d40ddbb9d0818")}

	inputDBSubnetGroup := docdb.CreateDBSubnetGroupInput{
		DBSubnetGroupDescription: aws.String(dbSubnetGroupName),
		DBSubnetGroupName:        aws.String(dbSubnetGroupName),
		SubnetIds:                subnetIDSlice,
		Tags:                     newTags,
	}
	//SubnetIDs:                aws.StringValueSlice(subnetIDSlice),

	searchDBSubnetGroup, err := o.client.GetDBSubnetGroup(ctx, &docdb.DescribeDBSubnetGroupsInput{DBSubnetGroupName: aws.String(dbSubnetGroupName)})
	if err != nil {
		if searchDBSubnetGroup == nil {
			log.Infoln("searchDBSubnetGroup is nil")
		}
		log.Infof("Failed to get existing DBSubnetGroup: %s\n", err)
		//return nil, apierror.New(apierror.ErrNotFound, "Failed to get existing DBSubnetGroup", err)
	}

	log.Debugf("searchDBClusterresult: %s\n", searchDBSubnetGroup)

	if searchDBSubnetGroup == nil {
		DBSubnetGroupCreateStatus, err := o.client.CreateDBSubnetGroup(ctx, &inputDBSubnetGroup)
		if err != nil {
			return nil, apierror.New(apierror.ErrBadRequest, "GOOGLELY failed to create DBSubnetGroup", err)
		}

		log.Debugf("DBSubnetGroup create status: %+v\n", DBSubnetGroupCreateStatus)
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

	clusterCreateStatus, err := o.client.CreateDBCluster(ctx, data.DBClusterIdentifier, &input)
	if err != nil {
		// fixme - the err doesn't seem to be returning to the caller, to the user, and it
		// probably will be helpful
		return nil, apierror.New(apierror.ErrBadRequest, "failed to create db cluster", err)
	}

	log.Debugf("cluster create status: %+v\n", clusterCreateStatus)

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

		// jam these back in with structs to support them
		//DBSubnetGroup:         *instanceCreateStatus.DBSubnetGroup,
		//Endpoint:              Endpoint,

		allInstances = append(allInstances, &loopInstance)

		log.Debugf("instance create status: %s\n", instanceCreateStatus)

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

	log.Debugf("Deleting documentDB instances and cluster: %s, %s\n", data.InstanceNames, data.ClusterName)

	// Get instances in cluster, so we can delete them without user giving input
	getCluster := docdb.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(name),
	}

	documentDB, err := o.client.GetDB(ctx, &getCluster)
	if err != nil {
		return "", apierror.New(apierror.ErrBadRequest, "failed to get documentDB", err)
	}

	log.Debugf("getting documentDB DBClusterOuput: %s\n", documentDB)
	log.Debugf("getting documentDB.DBSubnetGroup: %s\n", aws.StringValue(documentDB.DBSubnetGroup))

	// Loop through the DBClusterMember instances and delete them
	for _, iName := range documentDB.DBClusterMembers {
		log.Debugf("instanceName: %s\n", *iName.DBInstanceIdentifier)

		instanceDeleteInput := docdb.DeleteDBInstanceInput{
			DBInstanceIdentifier: aws.String(*iName.DBInstanceIdentifier),
		}

		instanceDeleteStatus, err := o.client.DeleteDBInstance(ctx, &instanceDeleteInput)
		if err != nil {
			log.Infof("Failed to delete db instance: %s\n", err)
		}

		log.Debugf("instanceDeleteStatus: %s\n", instanceDeleteStatus)
	}

	input := docdb.DeleteDBClusterInput{
		DBClusterIdentifier:       aws.String(data.ClusterName),
		SkipFinalSnapshot:         &data.SkipFinalSnapshot,
		FinalDBSnapshotIdentifier: &data.FinalDBSnapshotIdentifier,
	}

	clusterDeleteStatus, err := o.client.DeleteDBCluster(ctx, &input)
	if err != nil {
		return "failed to delete db cluster: ", err
	}

	log.Debugf("clusterDeleteStatus: %s\n", clusterDeleteStatus)

	return "{OK}", nil

}
