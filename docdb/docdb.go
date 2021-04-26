package docdb

import (
	"context"

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

// ListDB lists documentdb clusters
func (d *DocDB) ListDB(ctx context.Context, input *docdb.DescribeDBClustersInput) (*docdb.DescribeDBClustersOutput, error) {
	log.Info("listing documentDB clusters")

	out, err := d.Service.DescribeDBClusters(input)
	if err != nil {
		return nil, err
	}

	return out, err
}

// GetDB gets information on a documentDB cluster+instance
func (d *DocDB) GetDB(ctx context.Context, input *docdb.DescribeDBClustersInput) (*docdb.DescribeDBClustersOutput, error) {

	out, err := d.Service.DescribeDBClusters(input)
	if err != nil {
		return nil, err
	}

	return out, err

}

// CreateDBCluster creates a documentDB cluster
func (d *DocDB) CreateDBCluster(ctx context.Context, name string, input *docdb.CreateDBClusterInput) (*docdb.CreateDBClusterOutput, error) {

	out, err := d.Service.CreateDBCluster(input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// CreateDBInstance creates a documentDB instance
func (d *DocDB) CreateDBInstance(ctx context.Context, input *docdb.CreateDBInstanceInput) (*docdb.CreateDBInstanceOutput, error) {

	out, err := d.Service.CreateDBInstance(input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteDBCluster deletes a documentDB cluster
func (d *DocDB) DeleteDBCluster(ctx context.Context, input *docdb.DeleteDBClusterInput) (*docdb.DeleteDBClusterOutput, error) {

	out, err := d.Service.DeleteDBCluster(input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteDBInstance deletes a documentDB instance
func (d *DocDB) DeleteDBInstance(ctx context.Context, input *docdb.DeleteDBInstanceInput) (*docdb.DeleteDBInstanceOutput, error) {

	out, err := d.Service.DeleteDBInstance(input)
	if err != nil {
		return nil, err
	}

	return out, nil
}
