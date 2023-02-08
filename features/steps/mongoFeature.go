package steps

import (
	"context"
	"encoding/json"
	"fmt"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/datastore"
	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoVersion = "4.4.8"
)

type MongoFeature struct {
	componenttest.ErrorFeature
	*componenttest.MongoFeature

	cfg *config.Config
}

func NewMongoFeature(ef componenttest.ErrorFeature, cfg *config.Config) *MongoFeature {
	mf := &MongoFeature{
		ErrorFeature: ef,
		MongoFeature: componenttest.NewMongoFeature(componenttest.MongoOptions{
			MongoVersion: mongoVersion,
		}),
		cfg: cfg,
	}

	mf.cfg.Mongo.ClusterEndpoint = mf.MongoFeature.Server.URI()

	return mf
}

func (mf *MongoFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	mf.MongoFeature.RegisterSteps(ctx)

	ctx.Step(
		`^I have this metadata:$`,
		mf.iHaveThisMetadata,
	)

	ctx.Step(
		`^a document in collection "([^"]*)" with key "([^"]*)" value "([^"]*)" should match:$`,
		mf.aDocumentInCollectionWithKeyValueShouldMatch,
	)

}

func (mf *MongoFeature) iHaveThisMetadata(docs *godog.DocString) error {
	var populationTypeMetadata []datastore.DefaultDatasetMetadata

	err := json.Unmarshal([]byte(docs.Content), &populationTypeMetadata)
	if err != nil {
		return fmt.Errorf("failed to unmarshal filter: %w", err)
	}

	if err := mf.insertPopulationTypeMetadata(populationTypeMetadata); err != nil {
		return fmt.Errorf("error inserting filters: %w", err)
	}

	return nil
}

func (mf *MongoFeature) insertPopulationTypeMetadata(metadata []datastore.DefaultDatasetMetadata) error {
	ctx := context.Background()
	db := mf.cfg.Mongo.Database
	col := mf.cfg.MetadataCollection

	upsert := true
	for _, metadataItem := range metadata {
		if _, err := mf.Client.Database(db).Collection(col).UpdateByID(ctx, metadataItem.ID, bson.M{"$set": metadataItem}, &options.UpdateOptions{Upsert: &upsert}); err != nil {
			return fmt.Errorf("failed to upsert filter: %w", err)
		}
	}
	return nil
}

func (mf *MongoFeature) aDocumentInCollectionWithKeyValueShouldMatch(col, key, val string, doc *godog.DocString) error {
	var expected, result interface{}
	if err := json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return fmt.Errorf("failed to unmarshal document: %w", err)
	}

	var bdoc primitive.D
	if err := mf.Client.Database(mf.cfg.Mongo.Database).Collection(col).FindOne(context.Background(), bson.M{key: val}).Decode(&bdoc); err != nil {
		return fmt.Errorf("failed to retrieve document: %w", err)
	}

	b, err := bson.MarshalExtJSON(bdoc, true, true)
	if err != nil {
		return fmt.Errorf("failed to marshal bson document: %w", err)
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		return fmt.Errorf("-expected +got)\n%s\n", diff)
	}

	return nil
}
