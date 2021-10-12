package docdb

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/aws/aws-sdk-go/service/docdb/docdbiface"
	log "github.com/sirupsen/logrus"
)

// DocDB is a wrapper around the aws docdb service
type DocDB struct {
	session         *session.Session
	Service         docdbiface.DocDBAPI
	DefaultKMSKeyId string
}

type DocDBOption func(*DocDB)

func New(opts ...DocDBOption) DocDB {
	e := DocDB{}

	for _, opt := range opts {
		opt(&e)
	}

	if e.session != nil {
		e.Service = docdb.New(e.session)
	}

	return e
}

func WithSession(sess *session.Session) DocDBOption {
	return func(e *DocDB) {
		log.Debug("using aws session")
		e.session = sess
	}
}

func WithCredentials(key, secret, token, region string) DocDBOption {
	return func(e *DocDB) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		e.session = sess
	}
}

func WithDefaultKMSKeyId(keyId string) DocDBOption {
	return func(e *DocDB) {
		log.Debugf("using default kms keyid %s", keyId)
		e.DefaultKMSKeyId = keyId
	}
}

// GetDBSubnetGroup gets documentDB DBSubnetGroup by name
func (d *DocDB) GetDBSubnetGroup(ctx context.Context, name string) ([]*docdb.DBSubnetGroup, error) {
	log.Debugf("getting details for documentDB subnet group: %s", name)

	out, err := d.Service.DescribeDBSubnetGroups(&docdb.DescribeDBSubnetGroupsInput{DBSubnetGroupName: aws.String(name)})
	if err != nil {
		return nil, err
	}

	log.Debugf("search output for documentDB db subnet group: %+v", out.DBSubnetGroups)

	return out.DBSubnetGroups, nil
}

// ListDocDBs lists all documentDB clusters
func (d *DocDB) ListDocDBClusters(ctx context.Context) ([]string, error) {
	log.Debug("listing documentDB clusters")

	filters := []*docdb.Filter{
		{
			Name:   aws.String("engine"),
			Values: aws.StringSlice([]string{"docdb"}),
		},
	}

	clusters := []string{}
	if err := d.Service.DescribeDBClustersPagesWithContext(ctx,
		&docdb.DescribeDBClustersInput{Filters: filters},
		func(page *docdb.DescribeDBClustersOutput, lastPage bool) bool {
			for _, c := range page.DBClusters {
				clusters = append(clusters, aws.StringValue(c.DBClusterIdentifier))
			}

			return true
		}); err != nil {
		return nil, err
	}

	log.Debugf("listing documentDB clusters output: %+v", clusters)

	return clusters, nil
}

// GetDocDBDetails gets information about a documentDB cluster
func (d *DocDB) GetDocDBDetails(ctx context.Context, name string) (*docdb.DBCluster, error) {
	log.Debugf("getting information about documentDB cluster %s", name)

	out, err := d.Service.DescribeDBClustersWithContext(ctx, &docdb.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	if len(out.DBClusters) == 0 {
		msg := fmt.Sprintf("%s not found", name)
		return nil, apierror.New(apierror.ErrNotFound, msg, nil)
	}

	if num := len(out.DBClusters); num > 1 {
		msg := fmt.Sprintf("unexpected number of DBClusters found for %s (%d)", name, num)
		return nil, apierror.New(apierror.ErrInternalError, msg, nil)
	}

	log.Debugf("getting documentDB cluster and instance(s) output: %+v", out)

	return out.DBClusters[0], err
}

// GetDocDBInstances gets information about all instances in a documentDB cluster
func (d *DocDB) GetDocDBInstances(ctx context.Context, name string) ([]*docdb.DBInstance, error) {
	log.Debugf("getting information about documentDB instances in cluster %s", name)

	filters := []*docdb.Filter{
		{
			Name:   aws.String("db-cluster-id"),
			Values: aws.StringSlice([]string{name}),
		},
	}

	out, err := d.Service.DescribeDBInstancesWithContext(ctx, &docdb.DescribeDBInstancesInput{
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("getting documentDB instances output: %+v", out)

	return out.DBInstances, err
}

// GetDocDBTags gets the tags for a documentDB cluster
func (d *DocDB) GetDocDBTags(ctx context.Context, arn *string) ([]*docdb.Tag, error) {
	log.Debugf("getting tags for documentDB cluster %s", aws.StringValue(arn))

	out, err := d.Service.ListTagsForResourceWithContext(ctx, &docdb.ListTagsForResourceInput{
		ResourceName: arn,
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("getting documentDB tags output: %+v", out)

	return out.TagList, err
}

// CreateDBCluster creates a documentDB cluster
func (d *DocDB) CreateDBCluster(ctx context.Context, input *docdb.CreateDBClusterInput) (*docdb.DBCluster, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating documentDB cluster: %s", aws.StringValue(input.DBClusterIdentifier))

	out, err := d.Service.CreateDBCluster(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("created documentDB cluster with output: %+v", out.DBCluster)

	return out.DBCluster, nil
}

// CreateDBInstance creates a documentDB instance
func (d *DocDB) CreateDBInstance(ctx context.Context, input *docdb.CreateDBInstanceInput) (*docdb.DBInstance, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating documentDB instance: %s", aws.StringValue(input.DBInstanceIdentifier))

	out, err := d.Service.CreateDBInstance(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("created documentDB instance with output: %+v", out.DBInstance)

	return out.DBInstance, nil
}

// CreateDBSubnetGroup creates a documentDB DBSubnetGroup
func (d *DocDB) CreateDBSubnetGroup(ctx context.Context, input *docdb.CreateDBSubnetGroupInput) (*docdb.DBSubnetGroup, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating documentDB DBSubnetGroup: %s", aws.StringValue(input.DBSubnetGroupName))

	out, err := d.Service.CreateDBSubnetGroup(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("created documentDB DBSubnetGroup with output: %+v", out.DBSubnetGroup)

	return out.DBSubnetGroup, nil
}

// DeleteDBCluster deletes a documentDB cluster
func (d *DocDB) DeleteDBCluster(ctx context.Context, input *docdb.DeleteDBClusterInput) (*docdb.DeleteDBClusterOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("deleting documentDB cluster: %s", aws.StringValue(input.DBClusterIdentifier))

	out, err := d.Service.DeleteDBCluster(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("deleted documentDB cluster with ouput: %+v", out)

	return out, nil
}

// DeleteDBInstance deletes a documentDB instance
func (d *DocDB) DeleteDBInstance(ctx context.Context, input *docdb.DeleteDBInstanceInput) (*docdb.DeleteDBInstanceOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("deleting documentDB instance: %s", aws.StringValue(input.DBInstanceIdentifier))

	out, err := d.Service.DeleteDBInstance(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("deleted documentDB instance with output: %+v", out)

	return out, nil
}

// ModifyDBCluster modifies a documentDB cluster
func (d *DocDB) ModifyDBCluster(ctx context.Context, input *docdb.ModifyDBClusterInput) (*docdb.DBCluster, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("modifying documentDB cluster: %s", aws.StringValue(input.DBClusterIdentifier))

	out, err := d.Service.ModifyDBCluster(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("modified documentDB cluster with output: %+v", out.DBCluster)

	return out.DBCluster, nil
}

// ModifyDBInstance modifies a documentDB instance
func (d *DocDB) ModifyDBInstance(ctx context.Context, input *docdb.ModifyDBInstanceInput) (*docdb.DBInstance, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("modifying documentDB instance: %s", aws.StringValue(input.DBInstanceIdentifier))

	out, err := d.Service.ModifyDBInstance(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("modified documentDB instance with output: %+v", out.DBInstance)

	return out.DBInstance, nil
}
