package api

import (
	"strconv"
	"time"

	"github.com/YaleSpinup/docdb-api/common"
	"github.com/YaleSpinup/docdb-api/docdb"
	"github.com/YaleSpinup/docdb-api/resourcegroupstaggingapi"
	"github.com/YaleSpinup/flywheel"
)

type docDBOrchestrator struct {
	client   docdb.DocDB
	rgClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	flywheel *flywheel.Manager
	org      string
}

func newDocDBOrchestrator(client docdb.DocDB, rgclient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI, flywheel *flywheel.Manager, org string) *docDBOrchestrator {
	return &docDBOrchestrator{
		client:   client,
		rgClient: rgclient,
		flywheel: flywheel,
		org:      org,
	}
}

func newFlywheelManager(config common.Flywheel) (*flywheel.Manager, error) {
	opts := []flywheel.ManagerOption{}

	if config.RedisAddress != "" {
		opts = append(opts, flywheel.WithRedisAddress(config.RedisAddress))
	}

	if config.RedisUsername != "" {
		opts = append(opts, flywheel.WithRedisAddress(config.RedisUsername))
	}

	if config.RedisPassword != "" {
		opts = append(opts, flywheel.WithRedisAddress(config.RedisPassword))
	}

	if config.RedisDatabase != "" {
		db, err := strconv.Atoi(config.RedisDatabase)
		if err != nil {
			return nil, err
		}
		opts = append(opts, flywheel.WithRedisDatabase(db))
	}

	if config.TTL != "" {
		ttl, err := time.ParseDuration(config.TTL)
		if err != nil {
			return nil, err
		}
		opts = append(opts, flywheel.WithTTL(ttl))
	}

	manager, err := flywheel.NewManager(config.Namespace, opts...)
	if err != nil {
		return nil, err
	}

	return manager, nil
}
