package docdb

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/aws/aws-sdk-go/service/docdb/docdbiface"
	log "github.com/sirupsen/logrus"
)

// https://github.com/aws/aws-sdk-go/blob/main/service/docdb/docdbiface/interface.go

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

/*
func WithCredentials(key, secret, token, region string) ECROption {
	return func(e *ECR) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		e.session = sess
	}
}

func WithDefaultKMSKeyId(keyId string) ECROption {
	return func(e *ECR) {
		log.Debugf("using default kms keyid %s", keyId)
		e.DefaultKMSKeyId = keyId
	}
}
*/

// ListDB lists documentdb clusters
func (d *DocDB) ListDB(ctx context.Context, input *docdb.DescribeDBClustersInput) (*docdb.DescribeDBClustersOutput, error) {
	log.Info("listing documentDB clusters")

	// List clusters
	out, err := d.Service.DescribeDBClusters(input)
	if err != nil {
		//msg := fmt.Sprint("failed to get documentDB list")
		// return nil, ErrCode("failed to list documentDBs", err)
		return nil, err
	}

	return out, err
}

// CreateDB creates a documentDB cluster
func (d *DocDB) CreateDB(ctx context.Context, input *docdb.CreateDBClusterInput) (*docdb.CreateDBClusterOutput, error) {

	out, err := d.Service.CreateDBCluster(input)
	if err != nil {
		return nil, err
	}

	return out, nil
}
