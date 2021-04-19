package api

import (
	"context"

	//"github.com/YaleSpinup/docdb-api/docdb"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	log "github.com/sirupsen/logrus"
)

// createDocumentDB creates documentDB cluster and instances
func (o *docDBOrchestrator) createDocumentDB(ctx context.Context, name string, data *docdb.CreateDBClusterInput) (string, error) {

	log.Debugf("Creating documentDB instances and cluster: %s\n", name)
	log.Debugf("GOOGLEY ctx: %s\n", ctx)
	log.Debugf("GOOGLEY data: %s\n", data)

	/*
		input := docdb.CreateDBClusterInput{
			AvailabilityZones: []*string{
				aws.String("us-east-1a"),
				aws.String("us-east-1d"),
			},
			DBClusterIdentifier: aws.String(name),
			DBSubnetGroupName:   aws.String("default-vpc-0e7363e700630fab5"),
			Engine:              aws.String("docdb"),
			MasterUsername:      aws.String("foobarusername"),
			MasterUserPassword:  aws.String("foobarbazboo"),
			//Tags: "[]",
		}
	*/

	clusterCreateStatus, err := o.client.CreateDBCluster(ctx, name, data)
	if err != nil {
		return "blah", err
	}

	log.Debugf("GOOGLEY mydata: %s\n", clusterCreateStatus)

	instanceData := docdb.CreateDBInstanceInput{
		AutoMinorVersionUpgrade:    aws.Bool(true),
		AvailabilityZone:           aws.String("us-east-1a"),
		DBInstanceClass:            aws.String("db.t3.medium"),
		DBClusterIdentifier:        aws.String(name),
		DBInstanceIdentifier:       aws.String("instance1"),
		Engine:                     aws.String(*data.Engine),
		PreferredMaintenanceWindow: aws.String("Sun:04:00-Sun:04:30"),
		PromotionTier:              aws.Int64(1),
		//Tags:
	}

	log.Debugf("GOOGLEY instanceData: %s\n", instanceData)

	instanceCreateStatus, err := o.client.CreateDBInstance(ctx, &instanceData)
	if err != nil {
		return "failed to create db instance: ", err
	}

	log.Debugf("GOOGLEY instanceCreateStatus: %s\n", instanceCreateStatus)

	return "", nil

}

// deleteDocumentDB deletes documentDB instances and cluster
func (o *docDBOrchestrator) deleteDocumentDB(ctx context.Context, name string, data *DeleteDocDB) (string, error) {

	log.Debugf("Deleting documentDB instances and cluster %s\n", name)

	input := docdb.DeleteDBClusterInput{
		DBClusterIdentifier:       aws.String(name),
		SkipFinalSnapshot:         aws.Bool(data.SkipFinalSnapshot),
		FinalDBSnapshotIdentifier: aws.String(data.FinalDBSnapshotIdentifier),
	}

	clusterDeleteStatus, err := o.client.DeleteDBCluster(ctx, &input)
	if err != nil {
		return "failed to delete db cluster: ", err
	}

	log.Debugf("GOOGLEY clusterDeleteStatus: %s\n", clusterDeleteStatus)

	//instanceCreateStatus, err := o.client.DeleteDBInstance(ctx, &instanceData)

	return "", nil

}
