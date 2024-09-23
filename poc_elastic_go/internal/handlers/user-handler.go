package handlers

import (
	"fmt"
	"math/rand"
	"time"

	"poc_elastic_go/internal/domain"
	"poc_elastic_go/internal/repository"
)

type UserHandler struct {
	repository repository.UserRepository
}

func NewUserHandler(repository repository.UserRepository) *UserHandler {
	return &UserHandler{repository: repository}
}

func (uh *UserHandler) CreateUser(user domain.User) error {
	return uh.repository.CreateUser(user)
}

func (uh *UserHandler) GetAllUsers(page int, size int) ([]domain.User, error) {
	return uh.repository.GetAllUsers(page, size)
}

func (uh *UserHandler) FindUsersByName(name string) ([]domain.User, error) {
	return uh.repository.FindUsersByName(name)
}

func (uh *UserHandler) FindUsersByCity(city string) ([]domain.User, error) {
	return uh.repository.FindUsersByCity(city)
}

func (uh *UserHandler) FindUserByNRC(nrc int) (*domain.User, error) {
	return uh.repository.FindUserByNRC(nrc)
}

func (uh *UserHandler) CreateRandomUser(nrc int) (*domain.User, error) {
	city := randomCity()
	state := randomState(city)
	randomUser := domain.User{
		Name:               randomName(),
		Age:                rand.Intn(100),
		NRC:                nrc,
		DateOfRegistration: time.Now(),
		Address: domain.Address{
			City:   city,
			State:  state,
			Street: randomStreet(),
			Number: rand.Intn(1000),
		},
	}

	err := uh.repository.CreateUser(randomUser)
	if err != nil {
		return nil, err
	}

	return &randomUser, nil
}

func (uh *UserHandler) RandomUpdate(nrc int) (*domain.User, error) {
	user, err := uh.repository.FindUserByNRC(nrc)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	city := randomCity()
	state := randomState(city)
	user.Name = randomName()
	user.Age = rand.Intn(100)
	user.Address.City = city
	user.Address.State = state
	user.Address.Street = randomStreet()
	user.Address.Number = rand.Intn(1000)

	err = uh.repository.UpdateUser(*user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func randomName() string {
	names := []string{"João", "Maria", "Pedro", "Ana", "Carlos", "Fernanda", "Lucas", "Juliana", "Rafael", "Camila"}
	return names[rand.Intn(len(names))]
}

func randomCity() string {
	cities := []string{"São Paulo", "Rio de Janeiro", "Belo Horizonte", "Brasília", "Salvador", "Fortaleza", "Curitiba", "Recife", "Porto Alegre", "Manaus"}
	return cities[rand.Intn(len(cities))]
}

func randomState(city string) string {
	cityStateMap := map[string]string{
		"São Paulo":      "SP",
		"Rio de Janeiro": "RJ",
		"Belo Horizonte": "MG",
		"Brasília":       "DF",
		"Salvador":       "BA",
		"Fortaleza":      "CE",
		"Curitiba":       "PR",
		"Recife":         "PE",
		"Porto Alegre":   "RS",
		"Manaus":         "AM",
	}

	return cityStateMap[city]
}

func randomStreet() string {
	streets := []string{"Rua das Flores", "Avenida Paulista", "Rua do Sol", "Avenida Atlântica", "Rua da Luz", "Avenida Brasil", "Rua da Paz", "Avenida das Américas", "Rua da Liberdade", "Avenida Central"}
	return streets[rand.Intn(len(streets))]
}
