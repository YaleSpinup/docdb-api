package api

import "time"

// DeleteDocDB is data used to delete a documentDB
type DeleteDocDB struct {
	FinalDBSnapshotIdentifier string
	SkipFinalSnapshot         bool
	ClusterName               string
	InstanceNames             []string
}

// Tag provides metadata and billing information
type Tag struct {
	Key   *string
	Value *string
}

// CreateDocDB is data used to create a documentDB
type CreateDocDB struct {
	AvailabilityZones   []string
	InstanceCount       int
	DBClusterIdentifier string
	DBSubnetGroupName   string
	DBInstanceClass     string
	Engine              string
	MasterUsername      string
	MasterUserPassword  string
	MaintenanceWindow   string
	PromotionTier       int64
	StorageEncrypted    bool
	Tags                []*Tag
}

// ClusterMember is a map of cluster member info
type ClusterMember struct {
	DBClusterParameterGroupStatus string
	DBInstanceIdentifier          string
	IsClusterWriter               bool
	PromotionTier                 int64
}

// Something is a bunch of list returned things
type Something struct {
	DBClusterIdentifier string
	DBClusterMembers    map[string]ClusterMember
}

// ListDBReturn shows a subset of returned data
type ListDBReturn struct {
	DBClusters map[string]Something
}

//DBInstance helps us collect useful data from the upstream instance create call output
type DBInstance struct {
	AvailabilityZone      string
	BackupRetentionPeriod string
	DBInstanceArn         string
	DBInstanceClass       string
	DBInstanceStatus      string
	DBInstanceIdentifier  string
	DBSubnetGroup         string
	Endpoint              string
	Engine                string
	EngineVersion         string
	InstanceCreateTime    time.Time
	KmsKeyId              string
	ReaderEndpoint        string
	StorageEncrypted      bool
}

//DBCluster helps us collect useful data from the upstream Cluster create call output
type DBCluster struct {
	DBClusterArn        string
	DBClusterIdentifier string
	Endpoint            string
	ReaderEndpoint      string
	StorageEncrypted    bool
	DBSubnetGroup       string
	DBInstances         []*DBInstance
}

// Cluster is the DBCluster outer JSON Key
type Cluster struct {
	DBClusters DBCluster
}
