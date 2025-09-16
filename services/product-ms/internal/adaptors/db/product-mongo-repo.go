package db

import (
	"context"
	"fmt"

	//"product-microservice/internal/domain"
	//"product-microservice/internal/ports"

	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(db *mongo.Database) ports.ProductRepository {
	return &MongoProductRepository{collection: db.Collection("products")}
}


func (r *MongoProductRepository) CreateProduct(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	// Generate ObjectID and set as hex string
	objID := primitive.NewObjectID()
	p.ID = objID.Hex()

	// Store in Mongo with ObjectID, not string
	_, err := r.collection.InsertOne(ctx, bson.M{
		"_id":         objID,
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"stock":       p.Stock,
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}


func (r *MongoProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
       return nil, fmt.Errorf("invalid id: %v", err)
	}

	var product domain.Product
	 err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return  nil, err
	}

	product.ID = objectID.Hex()
	return &product, nil
}


func (r *MongoProductRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var products []domain.Product
    for cursor.Next(ctx) {
        var p domain.Product
        if err := cursor.Decode(&p); err != nil {
            return nil, err
        }
        products = append(products, p)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return products, nil
}


func (r *MongoProductRepository) Update(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"stock":       p.Stock,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, err
	}

	// fetch updated product
	var updated domain.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updated)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}



func (r *MongoProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
		return  fmt.Errorf("invalid id: %v", err)
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return   err
	}
	return  nil
}

