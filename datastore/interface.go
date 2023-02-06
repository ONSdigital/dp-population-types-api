package datastore

import (
	"context"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-mongodb/v3/health"
	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/pkg/errors"
)

type MongoClient struct {
	conn               *mongo.MongoConnection
	health             *health.CheckMongoClient
	cfg                Config
	MetadataCollection string
}

func NewClient(ctx context.Context, cfg Config) (*MongoClient, error) {
	c := MongoClient{
		cfg: cfg,
	}
	var err error
	if c.conn, err = mongo.Open(&cfg.MongoDriverConfig); err != nil {
		return nil, errors.Wrap(err, "failed to open mongodb connection: %w")
	}

	collectionBuilder := map[health.Database][]health.Collection{
		health.Database(cfg.Database): {
			health.Collection(cfg.MetadataCollection),
		},
	}

	c.health = health.NewClientWithCollections(c.conn, collectionBuilder)

	return &c, nil
}

// Conn returns the underlying mongodb connection.
func (c *MongoClient) Conn() *mongo.MongoConnection {
	return c.conn
}

// Close represents mongo session closing within the context deadline
func (c *MongoClient) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (c *MongoClient) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return c.health.Checker(ctx, state)
}
