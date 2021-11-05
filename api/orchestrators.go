package api

import (
	"context"
	"strconv"
	"time"

	"github.com/YaleSpinup/docdb-api/common"
	"github.com/YaleSpinup/docdb-api/docdb"
	"github.com/YaleSpinup/docdb-api/resourcegroupstaggingapi"
	"github.com/YaleSpinup/flywheel"
	log "github.com/sirupsen/logrus"
)

type docDBOrchestrator struct {
	server      *server
	sp          *sessionParams
	docdbClient docdb.DocDB
	rgClient    *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
}

// sessionParams stores all required parameters to initialize the connection session
type sessionParams struct {
	role         string
	inlinePolicy string
	policyArns   []string
}

// newDocDBOrchestrator creates a new session and initializes all clients
func (s *server) newDocDBOrchestrator(ctx context.Context, sp *sessionParams) (*docDBOrchestrator, error) {
	log.Debug("initializing docDBOrchestrator")

	sess, err := s.assumeRole(
		ctx,
		s.session.ExternalID,
		sp.role,
		sp.inlinePolicy,
		sp.policyArns...,
	)
	if err != nil {
		return nil, err
	}

	return &docDBOrchestrator{
		server:      s,
		sp:          sp,
		docdbClient: docdb.New(docdb.WithSession(sess.Session)),
		rgClient:    resourcegroupstaggingapi.New(resourcegroupstaggingapi.WithSession(sess.Session)),
	}, nil
}

// refreshSession refreshes the session for all client connections
func (o *docDBOrchestrator) refreshSession(ctx context.Context) error {
	log.Debug("refreshing docDBOrchestrator session")

	sess, err := o.server.assumeRole(
		ctx,
		o.server.session.ExternalID,
		o.sp.role,
		o.sp.inlinePolicy,
		o.sp.policyArns...,
	)
	if err != nil {
		return err
	}

	o.docdbClient = docdb.New(docdb.WithSession(sess.Session))
	o.rgClient = resourcegroupstaggingapi.New(resourcegroupstaggingapi.WithSession(sess.Session))

	return nil
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
