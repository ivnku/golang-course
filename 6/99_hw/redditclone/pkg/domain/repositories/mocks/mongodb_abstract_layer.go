package mocks

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockgen command:
// mockgen -source=mongodb_abstract_layer.go -destination=mongodb_abstract_layer_mock.go -package=mongodbabstractlayer IMongoDatabase

type IMongoDatabase interface {
	Collection(name string) IMongoCollection
}

type IMongoCollection interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (IMongoCursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (IMongoUpdateResult, error)
}

// Mongo SingleResult
type IMongoSingleResult interface {
	Decode(v interface{}) error
}

type MongoSingleResult struct {
	sr *mongo.SingleResult
}

func (msr *MongoSingleResult) Decode(v interface{}) error {
	return msr.sr.Decode(v)
}

// Mongo UpdateResult
type IMongoUpdateResult interface {
	UnmarshalBSON(b []byte) error
}

type MongoUpdateResult struct {
	ur *mongo.UpdateResult
}

func (mur *MongoUpdateResult) UnmarshalBSON(b []byte) error {
	return mur.ur.UnmarshalBSON(b)
}

// Mongo Cursor
type IMongoCursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
	All(ctx context.Context, results interface{}) error
}

type MongoCursor struct {
	cur *mongo.Cursor
}

func (mc *MongoCursor) Close(ctx context.Context) error {
	return mc.cur.Close(ctx)
}

func (mc *MongoCursor) Next(ctx context.Context) bool {
	return mc.cur.Next(ctx)
}

func (mc *MongoCursor) Decode(val interface{}) error {
	return mc.cur.Decode(val)
}

func (mc *MongoCursor) All(ctx context.Context, results interface{}) error {
	return mc.cur.All(ctx, results)
}

// Mongo Collection
type MongoCollection struct {
	Collection *mongo.Collection
}

func (mc *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error) {
	cursorResult, err := mc.Collection.Find(ctx, filter, opts...)
	return &MongoCursor{cur: cursorResult}, err
}

func (mc *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult {
	singleResult := mc.Collection.FindOne(ctx, filter, opts...)
	return &MongoSingleResult{sr: singleResult}
}

func (mc *MongoCollection) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*options.InsertOneOptions,
) (*mongo.InsertOneResult, error) {
	return mc.Collection.InsertOne(ctx, document, opts...)
}

func (mc *MongoCollection) DeleteOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {

	return mc.Collection.DeleteOne(ctx, filter, opts...)
}

func (mc *MongoCollection) Aggregate(
	ctx context.Context,
	pipeline interface{},
	opts ...*options.AggregateOptions,
) (IMongoCursor, error) {
	curr, err := mc.Collection.Aggregate(ctx, pipeline, opts...)
	return &MongoCursor{cur: curr}, err
}

func (mc *MongoCollection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (IMongoUpdateResult, error) {

	updateResult, err := mc.Collection.UpdateOne(ctx, filter, update, opts...)

	return &MongoUpdateResult{ur: updateResult}, err
}
