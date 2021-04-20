package api

// DeleteDocDB is data used to delete a documentDB
type DeleteDocDB struct {
	FinalDBSnapshotIdentifier string
	SkipFinalSnapshot         bool
	ClusterName               string
	InstanceNames             []string
}

// Tags provides metadata and billing information
type Tags struct {
	ChargingAccount     string
	OwnedBy             string
	OwnedByDepartment   string
	CreatedBy           string
	CreatedByDeptarment string
	Application         string
}

// Tag provides metadata and billing information
type Tag struct {
	Key   *string `type:"string"`
	Value *string `type:"string"`
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
