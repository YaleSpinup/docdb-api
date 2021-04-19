package api

import (
	"github.com/YaleSpinup/docdb-api/docdb"
)

type docDBOrchestrator struct {
	client docdb.DocDB
	org string
}

func newDocDBOrchestrator(client docdb.DocDB, org string) *docDBOrchestrator {
	return &docDBOrchestrator{
		client: client,
		org:    org,
	}
}