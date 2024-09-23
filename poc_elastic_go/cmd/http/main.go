package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"poc_elastic_go/internal/domain"
	"poc_elastic_go/internal/handlers"
	"poc_elastic_go/internal/repository"
	"strconv"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	var err error
	var client interface{}
	var userRepository repository.UserRepository

	backend := flag.String("backend", "elastic", "Choose the backend (elastic or mongo)")
	flag.Parse()

	fmt.Printf("Backend: %s\n", *backend)
	if *backend == "elastic" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient := &http.Client{Transport: tr}
		client, err = elastic.NewClient(elastic.SetURL("https://localhost:9200"), elastic.SetBasicAuth("elastic", "G+Rehd00aiBg8KpPKHNf"), elastic.SetHttpClient(httpClient), elastic.SetSniff(false))
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}
		userRepository = repository.NewUserElasticRepository(client.(*elastic.Client))
	} else if *backend == "mongo" {
		client, err = mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}
		userRepository = repository.NewUserMongoRepository(client.(*mongo.Client))
	} else {
		log.Fatalf("Invalid backend: %s", *backend)
	}

	userHandler := handlers.NewUserHandler(userRepository)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var user domain.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			fmt.Printf("Error decoding request body: %v\n", err)
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		err = userHandler.CreateUser(user)
		if err != nil {
			http.Error(w, "Error saving user to database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Client created successfully")
	})

	http.HandleFunc("/users/all", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		page, size := 1, 10
		pageParam := r.URL.Query().Get("page")
		if pageParam != "" {
			page, _ = strconv.Atoi(pageParam)
		}
		sizeParam := r.URL.Query().Get("size")
		if sizeParam != "" {
			size, _ = strconv.Atoi(sizeParam)
		}

		users, err := userHandler.GetAllUsers(page, size)
		if err != nil {
			http.Error(w, "Error fetching users from database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	http.HandleFunc("/users/findByName", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Name parameter is required", http.StatusBadRequest)
			return
		}

		users, err := userHandler.FindUsersByName(name)
		if err != nil {
			http.Error(w, "Error fetching users from database", http.StatusInternalServerError)
			return
		}

		fmt.Printf("users: %v\n", users)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	http.HandleFunc("/users/findByCity", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		city := r.URL.Query().Get("city")
		if city == "" {
			http.Error(w, "City parameter is required", http.StatusBadRequest)
			return
		}

		users, err := userHandler.FindUsersByCity(city)
		if err != nil {
			http.Error(w, "Error fetching users from database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	http.HandleFunc("/users/findByNRC", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		nrcParam := r.URL.Query().Get("nrc")
		if nrcParam == "" {
			http.Error(w, "NRC parameter is required", http.StatusBadRequest)
			return
		}

		nrc, err := strconv.Atoi(nrcParam)
		if err != nil {
			http.Error(w, "Invalid NRC parameter", http.StatusBadRequest)
			return
		}

		user, err := userHandler.FindUserByNRC(nrc)
		if err != nil {
			http.Error(w, "Error fetching user from database", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	http.HandleFunc("/users/random", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			NRC int `json:"nrc"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		user, err := userHandler.CreateRandomUser(data.NRC)
		if err != nil {
			http.Error(w, "Error creating random user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	http.HandleFunc("/users/random-update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			NRC int `json:"nrc"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		user, err := userHandler.RandomUpdate(data.NRC)
		if err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
