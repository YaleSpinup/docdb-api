package docdb

import (
	"context"

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
func (d *DocDB) GetDBSubnetGroup(ctx context.Context, input *docdb.DescribeDBSubnetGroupsInput) (*docdb.DescribeDBSubnetGroupsOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("Listing documentDB dbsubnetgroups: %+v\n", input.DBSubnetGroupName)

	out, err := d.Service.DescribeDBSubnetGroups(input)
	if err != nil {
		return nil, err
	}

	log.Debugf("GOOGLE dbsubnetgroups: %+v\n", out)

	return out, nil
}

// ListDB lists documentdb clusters
func (d *DocDB) ListDB(ctx context.Context, input *docdb.DescribeDBClustersInput) (*docdb.DescribeDBClustersOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("listing documentDB clusters with input %+v", input)

	out, err := d.Service.DescribeDBClusters(input)
	if err != nil {
		return nil, err
	}

	filter := []*docdb.DBCluster{}

	for _, cluster := range out.DBClusters {
		if aws.StringValue(cluster.Engine) == "docdb" {
			log.Debugf("docbd clusters name, engine: %s, %v\n", aws.StringValue(cluster.DBClusterIdentifier), aws.StringValue(cluster.Engine))
			filter = append(filter, cluster)
		}
	}

	filterOut := &docdb.DescribeDBClustersOutput{
		DBClusters: filter,
	}

	return filterOut, err
}

// GetDB gets information on a documentDB cluster+instance
func (d *DocDB) GetDB(ctx context.Context, input *docdb.DescribeDBClustersInput) (*docdb.DBCluster, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("getting documentDB cluster with input %+v", input)

	out, err := d.Service.DescribeDBClusters(input)
	if err != nil {
		return nil, err
	}

	if len(out.DBClusters) > 1 {
		return nil, apierror.New(apierror.ErrInternalError, "GetDB received more than one DBcluster", nil)
	}

	return out.DBClusters[0], err

}

// CreateDBCluster creates a documentDB cluster
func (d *DocDB) CreateDBCluster(ctx context.Context, name string, input *docdb.CreateDBClusterInput) (*docdb.DBCluster, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("creating documentDB cluster with input %+v", input)

	out, err := d.Service.CreateDBCluster(input)
	if err != nil {
		return nil, err
	}

	return out.DBCluster, nil
}

// CreateDBInstance creates a documentDB instance
func (d *DocDB) CreateDBInstance(ctx context.Context, input *docdb.CreateDBInstanceInput) (*docdb.DBInstance, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("creating documentDB instance with input %+v", input)

	out, err := d.Service.CreateDBInstance(input)
	if err != nil {
		return nil, err
	}

	return out.DBInstance, nil
}

// CreateDBSubnetGroup creates a documentDB DBSubnetGroup
func (d *DocDB) CreateDBSubnetGroup(ctx context.Context, input *docdb.CreateDBSubnetGroupInput) (*docdb.DBSubnetGroup, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("creating documentDB DBSubnetGroup with input %+v", input)

	out, err := d.Service.CreateDBSubnetGroup(input)
	if err != nil {
		return nil, err
	}

	return out.DBSubnetGroup, nil
}

// DeleteDBCluster deletes a documentDB cluster
func (d *DocDB) DeleteDBCluster(ctx context.Context, input *docdb.DeleteDBClusterInput) (*docdb.DeleteDBClusterOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("deleting documentDB cluster with input %+v", input)

	out, err := d.Service.DeleteDBCluster(input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteDBInstance deletes a documentDB instance
func (d *DocDB) DeleteDBInstance(ctx context.Context, input *docdb.DeleteDBInstanceInput) (*docdb.DeleteDBInstanceOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("deleting documentDB instance with input %+v", input)

	out, err := d.Service.DeleteDBInstance(input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteDBSubnetGroup deletes a documentDB DBSubnetGroup
func (d *DocDB) DeleteDBSubnetGroup(ctx context.Context, input *docdb.DeleteDBSubnetGroupInput) (string, error) {
	if input == nil {
		return "", apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Debugf("deleting documentDB DBSubnetGroup with input %+v", input)

	_, err := d.Service.DeleteDBSubnetGroup(input)
	if err != nil {
		return "", err
	}

	return "{\"OK\"}", nil
}
