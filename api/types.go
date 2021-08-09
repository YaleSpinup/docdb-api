package api

import (
	"github.com/aws/aws-sdk-go/aws"
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
	SubnetIds             []*string
	Tags                  []*Tag
	VpcSecurityGroupIds   []*string
}

// DocDBResponse is the output from documentDB operations
type DocDBResponse struct {
	// https://docs.aws.amazon.com/sdk-for-go/api/service/docdb/#DBCluster
	Cluster *docdb.DBCluster
	// https://docs.aws.amazon.com/sdk-for-go/api/service/docdb/#DBInstance
	Instances []*docdb.DBInstance
}

type Tag struct {
	Key   *string
	Value *string
}

// normalizeTags strips the org from the given tags and ensures it is set to the API org
func normalizeTags(org string, tags []*Tag) []*Tag {
	normalizedTags := []*Tag{
		{
			Key:   aws.String("spinup:org"),
			Value: aws.String(org),
		},
	}
	for _, t := range tags {
		if aws.StringValue(t.Key) == "spinup:org" || aws.StringValue(t.Key) == "yale:org" {
			continue
		}
		normalizedTags = append(normalizedTags, t)
	}

	return normalizedTags
}

// fromDocDBTags converts from RDS tags to api Tags
func fromDocDBTags(ecrTags []*docdb.Tag) []*Tag {
	tags := make([]*Tag, 0, len(ecrTags))
	for _, t := range ecrTags {
		tags = append(tags, &Tag{
			Key:   t.Key,
			Value: t.Value,
		})
	}
	return tags
}

// toDocDBTags converts from api Tags to RDS tags
func toDocDBTags(tags []*Tag) []*docdb.Tag {
	docdbTags := make([]*docdb.Tag, 0, len(tags))
	for _, t := range tags {
		docdbTags = append(docdbTags, &docdb.Tag{
			Key:   t.Key,
			Value: t.Value,
		})
	}
	return docdbTags
}
