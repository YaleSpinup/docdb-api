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
func (d *DocDB) GetDBSubnetGroup(ctx context.Context, input *docdb.DescribeDBSubnetGroupsInput) ([]*docdb.DBSubnetGroup, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("searching for documentDB dbsubnetgroup: %s\n", aws.StringValue(input.DBSubnetGroupName))

	out, err := d.Service.DescribeDBSubnetGroups(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("search output for documentDB dbsubnetgroup: %v\n", out.DBSubnetGroups)

	return out.DBSubnetGroups, nil
}

// ListDB lists documentdb clusters
func (d *DocDB) ListDB(ctx context.Context, input *docdb.DescribeDBClustersInput) (*docdb.DescribeDBClustersOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infoln("listing documentDB clusters and instance(s)")

	filters := []*docdb.Filter{
		{
			Name:   aws.String("engine"),
			Values: aws.StringSlice([]string{"docdb"}),
		},
	}

	out, err := d.Service.DescribeDBClusters(&docdb.DescribeDBClustersInput{Filters: filters})
	if err != nil {
		return nil, err
	}

	log.Debugf("listing documentDB clusters and instance(s) with output: %+v\n", out)

	return out, err
}

// GetDB gets information on a documentDB cluster+instance
func (d *DocDB) GetDB(ctx context.Context, input *docdb.DescribeDBClustersInput) (*docdb.DBCluster, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting documentDB cluster and instance(s): %+v", aws.StringValue(input.DBClusterIdentifier))

	out, err := d.Service.DescribeDBClusters(input)
	if err != nil {
		return nil, err
	}

	if len(out.DBClusters) == 0 {
		msg := fmt.Sprintf("%s not found", input)
		return nil, apierror.New(apierror.ErrNotFound, msg, nil)
	}

	if num := len(out.DBClusters); num > 1 {
		msg := fmt.Sprintf("unexpected number of DBClusters found for input %s (%d)", input, num)
		return nil, apierror.New(apierror.ErrInternalError, msg, nil)
	}

	log.Debugf("getting documentDB cluster and instance(s) with output: %+v", out)

	return out.DBClusters[0], err
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

	log.Debugf("created documentDB cluster with output: %+v\n", out.DBCluster)

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

	log.Debugf("created documentDB instance with output: %+v\n", out.DBInstance)

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
