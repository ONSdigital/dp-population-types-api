package datastore

import (
	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"
)

// Config holds the config for the mongodb store
type Config struct {
	mongo.MongoDriverConfig
	MetadataDatabase   string
	MetadataCollection string
}
