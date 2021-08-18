package api

import (
	"github.com/YaleSpinup/docdb-api/docdb"
	"github.com/YaleSpinup/docdb-api/resourcegroupstaggingapi"
)

type docDBOrchestrator struct {
	client   docdb.DocDB
	rgClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	org      string
}

func newDocDBOrchestrator(client docdb.DocDB, rgclient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI, org string) *docDBOrchestrator {
	return &docDBOrchestrator{
		client:   client,
		rgClient: rgclient,
		org:      org,
	}
}
