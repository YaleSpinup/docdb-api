package api

import (
	"github.com/aws/aws-sdk-go/service/docdb"
)

// DocDBCreateRequest is data used to create a documentDB
type DocDBCreateRequest struct {
	BackupRetentionPeriod *int64
	InstanceCount         *int
	DBClusterIdentifier   *string
	DBInstanceClass       *string
	EngineVersion         *string
	MasterUsername        *string
	MasterUserPassword    *string
	SubnetIds             []string
	Tags                  Tags
	VpcSecurityGroupIds   []*string
}

// DocDBModifyRequest is data used to modify a documentDB
type DocDBModifyRequest struct {
	BackupRetentionPeriod  *int64
	DBInstanceClass        *string
	EngineVersion          *string
	MasterUserPassword     *string
	NewDBClusterIdentifier *string
	Tags                   Tags
	VpcSecurityGroupIds    []*string
}

// DocDBResponse is the output from documentDB operations
type DocDBResponse struct {
	// https://docs.aws.amazon.com/sdk-for-go/api/service/docdb/#DBCluster
	Cluster *docdb.DBCluster
	// https://docs.aws.amazon.com/sdk-for-go/api/service/docdb/#DBInstance
	Instances []*docdb.DBInstance `json:",omitempty"`
	Tags      Tags                `json:",omitempty"`
}

type docDBInstanceStateChangeRequest struct {
	State string `json:"state"`
}
