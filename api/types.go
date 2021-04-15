package api

// DeleteDocDB is data used to delete a documentDB
type DeleteDocDB struct {
	FinalDBSnapshotIdentifier string
	SkipFinalSnapshot         bool
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

// CreateDocDB is data used to create a documentDB
type CreateDocDB struct {
	AvailabilityZones   []string
	DBClusterIdentifier string
	DBSubnetGroupName   string
	Engine              string
	MasterUsername      string
	MasterUserPassword  string
	Tags                map[string]Tags
}
