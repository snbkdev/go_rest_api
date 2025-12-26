package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type BoxOffice struct {
    Budget uint64 `json:"budget" bson:"budget"`
    Gross  uint64 `json:"gross" bson:"gross"`
}

type Movie struct {
    ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name       string             `json:"name" bson:"name"`
    Year       string             `json:"year" bson:"year"`
    Directors  []string           `json:"directors" bson:"directors"`
    Writers    []string           `json:"writers" bson:"writers"`
    BoxOffice  BoxOffice          `json:"boxOffice" bson:"boxOffice"`
}

type DB struct {
    client     *mongo.Client
    collection *mongo.Collection
}

func NewDB() (*DB, error) {
    // Загружаем .env
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found: %v", err)
    }

    mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
    databaseName := getEnv("MONGODB_DATABASE", "appdb")
    collectionName := getEnv("MONGODB_COLLECTION", "movies")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        return nil, err
    }

    if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }

    collection := client.Database(databaseName).Collection(collectionName)

    return &DB{
        client:     client,
        collection: collection,
    }, nil
}

func (db *DB) Close() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    db.client.Disconnect(ctx)
}

func (db *DB) GetMovie(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        writeJSONError(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    var movie Movie
    err = db.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&movie)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            writeJSONError(w, "Movie not found", http.StatusNotFound)
        } else {
            writeJSONError(w, "Database error", http.StatusInternalServerError)
        }
        return
    }

    writeJSONResponse(w, movie, http.StatusOK)
}

func (db *DB) PostMovie(w http.ResponseWriter, r *http.Request) {
    var movie Movie

    if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
        writeJSONError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    movie.ID = primitive.NewObjectID()

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    _, err := db.collection.InsertOne(ctx, movie)
    if err != nil {
        writeJSONError(w, "Failed to create movie", http.StatusInternalServerError)
        return
    }

    writeJSONResponse(w, movie, http.StatusCreated)
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func (db *DB) GetAllMovies(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    cursor, err := db.collection.Find(ctx, bson.M{})
    if err != nil {
        writeJSONError(w, "Failed to fetch movies", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var movies []Movie
    if err := cursor.All(ctx, &movies); err != nil {
        writeJSONError(w, "Failed to decode movies", http.StatusInternalServerError)
        return
    }

    writeJSONResponse(w, movies, http.StatusOK)
}

func (db *DB) DeleteMovie(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        writeJSONError(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    result, err := db.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        writeJSONError(w, "Failed to delete movie", http.StatusInternalServerError)
        return
    }

    if result.DeletedCount == 0 {
        writeJSONError(w, "Movie not found", http.StatusNotFound)
        return
    }

    writeJSONResponse(w, map[string]interface{}{
        "message": "Movie deleted successfully",
        "deletedCount": result.DeletedCount,
    }, http.StatusOK)
}

func main() {
    db, err := NewDB()
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer db.Close()

    log.Println("Connected to MongoDB successfully!")

    r := mux.NewRouter()

    api := r.PathPrefix("/v1").Subrouter()
    api.HandleFunc("/movies", db.GetAllMovies).Methods("GET")
    api.HandleFunc("/movies/{id}", db.GetMovie).Methods("GET")
    api.HandleFunc("/movies", db.PostMovie).Methods("POST")
    api.HandleFunc("/movies/{id}", db.DeleteMovie).Methods("DELETE")

    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        writeJSONResponse(w, map[string]string{"status": "healthy"}, http.StatusOK)
    }).Methods("GET")

    serverAddr := getEnv("SERVER_ADDRESS", "127.0.0.1:8000")

    srv := &http.Server{
        Handler:      r,
        Addr:         serverAddr,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Printf("Server starting on http://%s", serverAddr)
    log.Printf("API endpoints:")
    log.Printf("  GET    http://%s/v1/movies", serverAddr)
    log.Printf("  GET    http://%s/v1/movies/{id}", serverAddr)
    log.Printf("  POST   http://%s/v1/movies", serverAddr)
    log.Printf("  DELETE http://%s/v1/movies/{id}", serverAddr)
    log.Printf("  GET    http://%s/health", serverAddr)

    if err := srv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}