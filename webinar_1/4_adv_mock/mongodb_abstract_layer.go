package main

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
}

type IMongoSingleResult interface {
	Decode(v interface{}) error
}

type IMongoCursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
}

type MongoCollection struct {
	Сoll *mongo.Collection
}

type MongoSingleResult struct {
	sr *mongo.SingleResult
}

type MongoCursor struct {
	cur *mongo.Cursor
}

func (msr *MongoSingleResult) Decode(v interface{}) error {
	return msr.sr.Decode(v)
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

func (mc *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error) {
	cursorResult, err := mc.Сoll.Find(ctx, filter, opts...)
	return &MongoCursor{cur: cursorResult}, err
}

func (mc *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult {
	singleResult := mc.Сoll.FindOne(ctx, filter, opts...)
	return &MongoSingleResult{sr: singleResult}
}
