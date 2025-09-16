package db

import (
	"context"
	"time"

	//"cart-microservice/internal/domain"
	//"cart-microservice/internal/ports"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCartRepo struct {
	col *mongo.Collection
}

func NewMongoCartRepo(db *mongo.Database) ports.CartRepository {
	return &MongoCartRepo{
		col: db.Collection("carts"),
	}
}

func (r *MongoCartRepo) GetCart(userID string) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.col.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return &domain.Cart{UserID: userID, Items: []domain.CartItem{}}, nil
	}
	return &cart, err
}

func (r *MongoCartRepo) AddItem(userID string, item domain.CartItem) error {
	item.AddedAt = time.Now()
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$push": bson.M{"items": item},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func (r *MongoCartRepo) RemoveItem(userID, productID string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$pull": bson.M{"items": bson.M{"product_id": productID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.col.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoCartRepo) ClearCart(userID string) error {
	_, err := r.col.DeleteOne(context.TODO(), bson.M{"user_id": userID})
	return err
}
