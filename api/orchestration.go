package api

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	log "github.com/sirupsen/logrus"
)

// createDocumentDB creates documentDB cluster and instances
func (o *docDBOrchestrator) createDocumentDB(ctx context.Context, data *CreateDocDB) ([]string, error) {
	if &data.InstanceCount == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid data: missing InstanceCount", nil)
	}

	output := []string{}

	log.Debugf("Creating documentDB instances and cluster: %s\n", data.DBClusterIdentifier)

	log.Debugf("GOOGLEY data: %s\n", data)
	log.Debugf("GOOGLEY tags: %s\n", data.Tags)

	tags := make([]*docdb.Tag, 0, len(data.Tags))

	for _, t := range data.Tags {
		tags = append(tags, &docdb.Tag{
			Key:   (t.Key),
			Value: (t.Value),
		})
	}

	log.Debugf("GOOGLEY data.Tags: %s\n", data.Tags)
	log.Debugf("GOOGLEY tags: %s\n", tags)

	input := docdb.CreateDBClusterInput{
		AvailabilityZones: []*string{
			aws.String("us-east-1a"),
			aws.String("us-east-1d"),
		},
		DBClusterIdentifier: &data.DBClusterIdentifier,
		DBSubnetGroupName:   &data.DBSubnetGroupName,
		Engine:              &data.Engine,
		MasterUsername:      &data.MasterUsername,
		MasterUserPassword:  &data.MasterUserPassword,
		Tags:                tags,
	}

	log.Debugf("tags input: %s\n", input)
	clusterCreateStatus, err := o.client.CreateDBCluster(ctx, data.DBClusterIdentifier, &input)
	if err != nil {
		return []string{"failed to create db cluster"}, err
	}

	output = append(output, fmt.Sprint(clusterCreateStatus))

	// create instances based on InstanceCount sent in
	// don't begin a 0 instance
	data.InstanceCount++
	for i := 1; i < data.InstanceCount; i++ {
		// normalize instanceName
		instanceName := fmt.Sprintf("%s-%v", data.DBClusterIdentifier, i)

		instanceData := docdb.CreateDBInstanceInput{
			AutoMinorVersionUpgrade:    aws.Bool(true),
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
			return []string{"failed to create db instance: "}, err
		}

		output = append(output, fmt.Sprint(instanceCreateStatus))

	}

	return output, nil

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
func (o *docDBOrchestrator) getDocumentDB(ctx context.Context, name string) (*docdb.DescribeDBClustersOutput, error) {

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
func (o *docDBOrchestrator) deleteDocumentDB(ctx context.Context, data *DeleteDocDB) (string, error) {

	log.Debugf("Deleting documentDB instances and cluster: %s, %s\n", data.InstanceNames, data.ClusterName)

	for _, iName := range data.InstanceNames {
		log.Debugf("instanceName: %s\n", iName)

		instanceDeleteInput := docdb.DeleteDBInstanceInput{
			DBInstanceIdentifier: aws.String(iName),
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

	return "", nil

}
