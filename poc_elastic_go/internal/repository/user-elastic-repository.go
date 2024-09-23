package repository

import (
	"context"
	"log"
	"poc_elastic_go/internal/domain"
	"reflect"
	"time"

	"github.com/olivere/elastic/v7"
)

type UserRepository interface {
	CreateUser(domain.User) error
	UpdateUser(domain.User) error
	GetAllUsers(int, int) ([]domain.User, error)
	FindUsersByName(string) ([]domain.User, error)
	FindUsersByCity(string) ([]domain.User, error)
	FindUserByNRC(int) (*domain.User, error)
}

type UserElasticRepository struct {
	client *elastic.Client
}

func NewUserElasticRepository(client *elastic.Client) *UserElasticRepository {
	return &UserElasticRepository{client: client}
}

func (ur *UserElasticRepository) CreateUser(user domain.User) error {
	user.DateOfRegistration = time.Now()

	_, err := ur.client.Index().
		Index("users").
		BodyJson(user).
		Do(context.Background())
	if err != nil {
		log.Printf("Error saving user to Elasticsearch: %v", err)
		return err
	}
	return nil
}

func (ur *UserElasticRepository) UpdateUser(user domain.User) error {
	_, err := ur.client.Update().
		Index("users").
		Id(user.ID).
		Doc(user).
		Do(context.Background())
	if err != nil {
		log.Printf("Error updating user in Elasticsearch: %v", err)
		return err
	}
	return nil
}

func (ur *UserElasticRepository) GetAllUsers(page int, size int) ([]domain.User, error) {
	searchResult, err := ur.client.Search().
		Index("users").
		Sort("date_of_registration", false).
		From((page - 1) * size).
		Size(size).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		log.Printf("Error fetching users from Elasticsearch: %v", err)
		return nil, err
	}

	var users []domain.User
	for idx, item := range searchResult.Each(reflect.TypeOf(domain.User{})) {
		user := item.(domain.User)
		user.ID = searchResult.Hits.Hits[idx].Id // Assuming the order is maintained
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserElasticRepository) FindUsersByName(name string) ([]domain.User, error) {
	searchResult, err := ur.client.Search().
		Index("users").
		Query(elastic.NewPrefixQuery("name", name)).
		Do(context.Background())
	if err != nil {
		log.Printf("Error fetching users from Elasticsearch: %v", err)
		return nil, err
	}

	var users []domain.User
	for idx, item := range searchResult.Each(reflect.TypeOf(domain.User{})) {
		user := item.(domain.User)
		user.ID = searchResult.Hits.Hits[idx].Id
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserElasticRepository) FindUsersByCity(city string) ([]domain.User, error) {
	searchResult, err := ur.client.Search().
		Index("users").
		Query(elastic.NewPrefixQuery("address.city", city)).
		Do(context.Background())
	if err != nil {
		log.Printf("Error fetching users from Elasticsearch: %v", err)
		return nil, err
	}

	var users []domain.User
	for idx, item := range searchResult.Each(reflect.TypeOf(domain.User{})) {
		user := item.(domain.User)
		user.ID = searchResult.Hits.Hits[idx].Id
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserElasticRepository) FindUserByNRC(nrc int) (*domain.User, error) {
	searchResult, err := ur.client.Search().
		Index("users").
		Query(elastic.NewTermQuery("nrc", nrc)).
		Do(context.Background())
	if err != nil {
		log.Printf("Error fetching user from Elasticsearch: %v", err)
		return nil, err
	}

	if searchResult.Hits.TotalHits.Value == 0 {
		return nil, nil
	}

	var user domain.User
	for idx, item := range searchResult.Each(reflect.TypeOf(domain.User{})) {
		user = item.(domain.User)
		user.ID = searchResult.Hits.Hits[idx].Id
		break
	}

	return &user, nil
}
