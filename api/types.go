package api

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
	//Key   *string `type:"string"`
	//Value *string `type:"string"`
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
