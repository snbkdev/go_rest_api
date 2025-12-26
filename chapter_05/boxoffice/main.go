package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Структуры данных
type BoxOffice struct {
    Budget uint64 `bson:"budget"`
    Gross  uint64 `bson:"gross"`
}

type Movie struct {
    Name       string   `bson:"name"`
    Year       string   `bson:"year"`
    Directors  []string `bson:"directors"`
    Writers    []string `bson:"writers"`
    BoxOffice  BoxOffice `bson:"boxOffice"`
}

func main() {
    // 1. Загружаем переменные окружения
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found: %v", err)
    }

    // 2. Получаем настройки из .env или используем значения по умолчанию
    mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
    databaseName := getEnv("MONGODB_DATABASE", "appdb")
    collectionName := getEnv("MONGODB_COLLECTION", "movies")

    // 3. Подключаемся к MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer func() {
        if err := client.Disconnect(ctx); err != nil {
            log.Printf("Warning: error disconnecting from MongoDB: %v", err)
        }
    }()

    // 4. Проверяем соединение
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }
    log.Println("Connected to MongoDB successfully!")

    // 5. Получаем коллекцию
    collection := client.Database(databaseName).Collection(collectionName)

    // 6. Создаем документ для вставки
    darkKnight := Movie{
        Name:      "The Dark Knight",
        Year:      "2008",
        Directors: []string{"Christopher Nolan"},
        Writers:   []string{"Jonathan Nolan", "Christopher Nolan"},
        BoxOffice: BoxOffice{
            Budget: 185000000,
            Gross:  1087900204,
        },
    }

    // 7. Вставляем документ
    insertCtx, insertCancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer insertCancel()

    insertResult, err := collection.InsertOne(insertCtx, darkKnight)
    if err != nil {
        log.Fatalf("Failed to insert document: %v", err)
    }
    log.Printf("Inserted document with ID: %v", insertResult.InsertedID)

    // 8. Ищем документ (бюджет > 150,000,000)
    filter := bson.M{"boxOffice.budget": bson.M{"$gt": 150000000}}

    findCtx, findCancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer findCancel()

    var result Movie
    err = collection.FindOne(findCtx, filter).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            log.Println("No movies found with budget > 150,000,000")
        } else {
            log.Fatalf("Failed to find document: %v", err)
        }
        return
    }

    // 9. Выводим результат
    fmt.Printf("Movie found: %s (%s)\n", result.Name, result.Year)
    fmt.Printf("Budget: $%d, Gross: $%d\n", result.BoxOffice.Budget, result.BoxOffice.Gross)
    fmt.Printf("Directors: %v\n", result.Directors)
    fmt.Printf("Writers: %v\n", result.Writers)
}

// Вспомогательная функция для получения переменных окружения
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}