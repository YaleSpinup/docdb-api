package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
)

type Tag struct {
	Key   string
	Value string
}

type Tags []Tag

// inOrg returns true if there is a spinup:org tag that matches our org
func (tags *Tags) inOrg(org string) bool {
	for _, t := range *tags {
		if t.Key == "spinup:org" && t.Value == org {
			return true
		}
	}
	return false
}

// normalize sets required tags
func (tags *Tags) normalize(org string) Tags {
	normalizedTags := Tags{
		{
			Key:   "spinup:org",
			Value: org,
		},
		{
			Key:   "spinup:type",
			Value: "database",
		},
		{
			Key:   "spinup:flavor",
			Value: "docdb",
		},
	}

	for _, t := range *tags {
		switch t.Key {
		case "yale:org", "spinup:org", "spinup:type", "spinup:flavor":
			continue
		default:
			normalizedTags = append(normalizedTags, t)
		}
	}

	return normalizedTags
}

// toDocDBTags converts from api Tags to RDS tags
func (tags *Tags) toDocDBTags() []*docdb.Tag {
	docdbTags := make([]*docdb.Tag, 0, len(*tags))
	for _, t := range *tags {
		docdbTags = append(docdbTags, &docdb.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}
	return docdbTags
}

// fromDocDBTags converts from DocDB tags to api Tags
func fromDocDBTags(docdbTags []*docdb.Tag) Tags {
	tags := make(Tags, 0, len(docdbTags))
	for _, t := range docdbTags {
		tags = append(tags, Tag{
			Key:   aws.StringValue(t.Key),
			Value: aws.StringValue(t.Value),
		})
	}
	return tags
}
