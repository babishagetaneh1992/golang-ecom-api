package db

import (
	//"user-microservice/internal/ports"

	"context"
	"fmt"

	//"user-microservice/internal/domain"
	//"user-microservice/internal/ports"

	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) ports.UserRepository {
	return  &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) Create(user *domain.User) error {
	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *MongoUserRepository) GetById(id string) (*domain.User, error) {
	var u domain.User

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return  nil, fmt.Errorf("invalid id: %v", err)
	}

	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}


func (r *MongoUserRepository) GetAll() ([]domain.User, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return  nil, err
	}

	var users []domain.User
    if err := cursor.All(context.Background(), & users); err != nil {
		return  nil, err
	}

	return  users, nil
}

func (r *MongoUserRepository) Exists(id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return  false, err
	}

	count, err := r.collection.CountDocuments(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return  false, err
	}

	return  count > 0, nil

}

func (r *MongoUserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return nil, nil
		}
		
		return nil, err
	}
	return &user, nil
}

