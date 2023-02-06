package datastore

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type PopulationTypeMetadata struct {
	// This ID refers to the Population Type
	ID               string `bson:"id"`
	DefaultDatasetID string `bson:"default-dataset-id"`
}

// GetMetadataRecord gets the metadata stored against a particular population type
// currently only just the default dataset
func (c *MongoClient) GetMetadataRecord(ctx context.Context, populationType string) (*PopulationTypeMetadata, error) {
	var metadata PopulationTypeMetadata

	if err := c.conn.Collection(c.MetadataCollection).FindOne(ctx, bson.M{"id": populationType}, &metadata); err != nil {
		return nil, &er{
			err:      errors.Wrap(err, "failed to find metadata"),
			notFound: errors.Is(err, mongodb.ErrNoDocumentFound),
		}

	}
	return &metadata, nil
}

func (c *MongoClient) PutMetadataRecord(ctx context.Context, metadata PopulationTypeMetadata) error {
	var err error

	if _, err = c.conn.Collection(c.MetadataCollection).UpsertById(ctx, metadata.ID, bson.M{"$set": metadata}); err != nil {
		fmt.Printf("%+v", err)
		return err
	}

	return nil
}
