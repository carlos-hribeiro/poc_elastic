package repository

import (
	"context"
	"log"
	"poc_elastic_go/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserMongoRepository struct {
	client *mongo.Client
}

func NewUserMongoRepository(client *mongo.Client) *UserMongoRepository {
	return &UserMongoRepository{client: client}
}

func (ur *UserMongoRepository) CreateUser(user domain.User) error {
	user.DateOfRegistration = time.Now()
	_, err := ur.client.Database("test").Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Error saving user to MongoDB: %v", err)
		return err
	}

	return nil
}

func (ur *UserMongoRepository) UpdateUser(user domain.User) error {
	_, err := ur.client.Database("test").Collection("users").UpdateOne(context.Background(), bson.M{"id": user.ID}, bson.M{"$set": user})
	if err != nil {
		log.Printf("Error updating user in MongoDB: %v", err)
		return err
	}
	return nil
}

func (ur *UserMongoRepository) GetAllUsers(page int, size int) ([]domain.User, error) {
	opts := options.Find().SetSkip(int64((page - 1) * size)).SetLimit(int64(size)).SetSort(bson.D{{"date_of_registration", -1}})
	cursor, err := ur.client.Database("test").Collection("users").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		log.Printf("Error fetching users from MongoDB: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []domain.User
	if err := cursor.All(context.Background(), &users); err != nil {
		log.Printf("Error fetching users from MongoDB: %v", err)
		return nil, err
	}
	return users, nil
}

func (ur *UserMongoRepository) FindUsersByName(name string) ([]domain.User, error) {
	cursor, err := ur.client.Database("test").Collection("users").Find(context.Background(), bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	if err != nil {
		log.Printf("Error fetching users from MongoDB: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []domain.User
	if err := cursor.All(context.Background(), &users); err != nil {
		log.Printf("Error fetching users from MongoDB: %v", err)
		return nil, err
	}
	return users, nil
}

func (ur *UserMongoRepository) FindUsersByCity(city string) ([]domain.User, error) {
	cursor, err := ur.client.Database("test").Collection("users").Find(context.Background(), bson.M{"address.city": bson.M{"$regex": city, "$options": "i"}})
	if err != nil {
		log.Printf("Error fetching users from MongoDB: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []domain.User
	if err := cursor.All(context.Background(), &users); err != nil {
		log.Printf("Error fetching users from MongoDB: %v", err)
		return nil, err
	}
	return users, nil
}

func (ur *UserMongoRepository) FindUserByNRC(nrc int) (*domain.User, error) {
	var user domain.User
	err := ur.client.Database("test").Collection("users").FindOne(context.Background(), bson.M{"nrc": nrc}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Error fetching user from MongoDB: %v", err)
		return nil, err
	}
	return &user, nil
}
