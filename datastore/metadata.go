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

// GetDefaultDatasetPopulationTypes returns a list of population types for which there is default dataset information
func (c *MongoClient) GetDefaultDatasetPopulationTypes(ctx context.Context) ([]string, error) {
	var datasetMetadata []DefaultDatasetMetadata

	// assumption that number of default datasets will be low, which seems very reasonable
	_, err := c.conn.Collection(c.MetadataCollection).Find(ctx, bson.M{}, &datasetMetadata)
	if err != nil {
		return nil, &er{
			err:      errors.Wrap(err, "failed to find default dataset metadata"),
			notFound: errors.Is(err, mongodb.ErrNoDocumentFound),
		}
	}

	populationTypes := make([]string, len(datasetMetadata))

	for _, metadata := range datasetMetadata {
		populationTypes = append(populationTypes, metadata.ID)
	}

	return populationTypes, nil
}

// GetDefaultDatasetMetadata gets the metadata stored against a particular population type
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

// PutDefaultDatasetMetadata puts the metadata stored against a particular population type
// by upserting the PutDefaultDatasetMetadata struct
func (c *MongoClient) PutDefaultDatasetMetadata(ctx context.Context, metadata DefaultDatasetMetadata) error {
	var err error

	if _, err = c.conn.Collection(c.MetadataCollection).UpsertById(ctx, metadata.ID, bson.M{"$set": metadata}); err != nil {
		return err
	}

	return nil
}
