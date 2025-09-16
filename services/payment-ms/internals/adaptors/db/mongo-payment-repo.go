package db

import (
	"context"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/domain"
	//"payment-ms/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPaymentRepository struct {
	collection *mongo.Collection
}

func NewMongoPaymentRepository(db *mongo.Database) *MongoPaymentRepository {
	return &MongoPaymentRepository{collection: db.Collection("payments")}
}

func (r *MongoPaymentRepository) Create(ctx context.Context, p *domain.Payment) (*domain.Payment, error) {
	 res, err := r.collection.InsertOne(ctx, p)
    if err != nil {
        return nil, err
    }

	 if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
        p.ID = oid.Hex()
    }

    return p, nil
}

func (r *MongoPaymentRepository) FindByID(ctx context.Context, id string) (*domain.Payment, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var payment domain.Payment
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&payment)
	return &payment, err
}

func (r *MongoPaymentRepository) List(ctx context.Context) ([]*domain.Payment, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var payments []*domain.Payment
	for cur.Next(ctx) {
		var p domain.Payment
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		payments = append(payments, &p)
	}
	return payments, nil
}

func (r *MongoPaymentRepository) UpdateStatus(ctx context.Context, id string, status string) (*domain.Payment, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    _, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": status}})
    if err != nil {
        return nil, err
    }

    var updated domain.Payment
    err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updated)
    if err != nil {
        return nil, err
    }

    return &updated, nil
}


func (r *MongoPaymentRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
