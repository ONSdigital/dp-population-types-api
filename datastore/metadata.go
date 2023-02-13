package datastore

import (
	"context"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type DefaultDatasetMetadata struct {
	// This ID refers to the Population Type
	ID               string `json:"id" bson:"id"`
	DefaultDatasetID string `json:"default_dataset_id" bson:"default_dataset_id"`
	Edition          string `json:"edition" bson:"edition"`
	Version          int    `json:"version" bson:"version"`
}

// GetMetadataRecord gets the metadata stored against a particular population type
// currently only just the default dataset
func (c *MongoClient) GetDefaultDatasetMetadata(ctx context.Context, populationType string) (*DefaultDatasetMetadata, error) {
	var metadata DefaultDatasetMetadata

	if err := c.conn.Collection(c.MetadataCollection).FindOne(ctx, bson.M{"id": populationType}, &metadata); err != nil {
		return nil, &er{
			err:      errors.Wrap(err, "failed to find metadata"),
			notFound: errors.Is(err, mongodb.ErrNoDocumentFound),
		}

	}
	return &metadata, nil
}

func (c *MongoClient) PutDefaultDatasetMetadata(ctx context.Context, metadata DefaultDatasetMetadata) error {
	var err error

	if _, err = c.conn.Collection(c.MetadataCollection).UpsertById(ctx, metadata.ID, bson.M{"$set": metadata}); err != nil {
		return err
	}

	return nil
}
