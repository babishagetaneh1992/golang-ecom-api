package db

import (
	"context"
	"fmt"

	//"order-microservice/internals/domain"

	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOrderRepository struct {
	collection *mongo.Collection
}

func NewMongoOrderRepository(db *mongo.Database) *MongoOrderRepository {
	return &MongoOrderRepository{
		collection: db.Collection("orders"),
	}
}

func (r *MongoOrderRepository) Create(ctx context.Context, o *domain.Order) (*domain.Order, error) {
	oid := primitive.NewObjectID()

	// build Mongo doc manually so _id is ObjectID
	doc := bson.M{
		"_id":     oid,
		"user_id": o.UserID,
		"items":   o.Items,
		"total":   o.Total,
		"status":  o.Status,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	// assign back as string for domain
	o.ID = oid.Hex()
	return o, nil
}

func (r *MongoOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result struct {
		ID     primitive.ObjectID `bson:"_id"`
		UserID string             `bson:"user_id"`
		Items  []domain.OrderItem `bson:"items"`
		Total  float64            `bson:"total"`
		Status string             `bson:"status"`
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:     result.ID.Hex(),
		UserID: result.UserID,
		Items:  result.Items,
		Total:  result.Total,
		Status: result.Status,
	}, nil
}

func (r *MongoOrderRepository) List(ctx context.Context) ([]*domain.Order, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var orders []*domain.Order
	for cur.Next(ctx) {
		var result struct {
			ID     primitive.ObjectID `bson:"_id"`
			UserID string             `bson:"user_id"`
			Items  []domain.OrderItem `bson:"items"`
			Total  float64            `bson:"total"`
			Status string             `bson:"status"`
		}

		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		orders = append(orders, &domain.Order{
			ID:     result.ID.Hex(),
			UserID: result.UserID,
			Items:  result.Items,
			Total:  result.Total,
			Status: result.Status,
		})
	}

	return orders, nil
}

func (r *MongoOrderRepository) UpdateOrderStatus(ctx context.Context, id string, status string) (*domain.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID %s: %w", id, err)
	}

	res, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		return nil, fmt.Errorf("mongo update error: %w", err)
	}

	// 🔍 Debug logging
	fmt.Printf("🔄 UpdateOrderStatus: id=%s status=%s matched=%d modified=%d\n",
		id, status, res.MatchedCount, res.ModifiedCount)

	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("no order found with id=%s", id)
	}

	// fetch updated doc
	return r.FindByID(ctx, id)
}


func (r *MongoOrderRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
