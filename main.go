package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task represents a task item
type Task struct {
	ID           string   `json:"id" bson:"_id"`
	Description  string   `json:"description" bson:"description"`
	Deadline     string   `json:"deadline" bson:"deadline"`
	TimeRequired string   `json:"timeRequired" bson:"timeRequired"`
	Priority     string   `json:"priority" bson:"priority"`
	Urgency      int      `json:"urgency" bson:"urgency"`
	Dependencies []string `json:"dependencies" bson:"dependencies"`
	Resources    []string `json:"resources" bson:"resources"`
	Subtasks     []string `json:"subtasks" bson:"subtasks"`
	Group        *string  `json:"group,omitempty" bson:"group,omitempty"`
}

const (
	mongoURI       = "mongodb://localhost:27017"
	databaseName   = "taskmanager"
	collectionName = "tasks"
)

var (
	collection *mongo.Collection
	ctx        context.Context
	client     *mongo.Client
)

func initDatabase() error {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Successfully connected to MongoDB!")

	// Get database and collection
	db := client.Database(databaseName)
	collection = db.Collection(collectionName)

	log.Printf("Database '%s' and collection '%s' initialized successfully!", databaseName, collectionName)

	// Insert a sample task if collection is empty
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count == 0 {
		sampleTask := Task{
			ID:           uuid.New().String(),
			Description:  "Sample Task",
			Deadline:     "2024-12-31",
			TimeRequired: "2h",
			Priority:     "Medium",
			Urgency:      3,
			Dependencies: []string{},
			Resources:    []string{"Computer"},
			Subtasks:     []string{"Step 1", "Step 2"},
		}

		_, err = collection.InsertOne(ctx, sampleTask)
		if err != nil {
			return err
		}
		log.Println("Sample task created successfully!")
	}

	return nil
}

func main() {
	// Initialize database
	if err := initDatabase(); err != nil {
		log.Fatal("Could not initialize database:", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/api/tasks", getTasks).Methods("GET")
	router.HandleFunc("/api/tasks", createTask).Methods("POST")
	router.HandleFunc("/api/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/api/tasks/{id}", deleteTask).Methods("DELETE")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start server
	log.Println("Server starting on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", c.Handler(router)))
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find all tasks in the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Decode the results
	var tasks []Task
	if err = cursor.All(context.Background(), &tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()

	// Insert the task into MongoDB
	_, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var task Task
	err := collection.FindOne(context.Background(), bson.M{"_id": params["id"]}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = params["id"]

	// Update the task in MongoDB
	result, err := collection.ReplaceOne(
		context.Background(),
		bson.M{"_id": params["id"]},
		task,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": params["id"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
