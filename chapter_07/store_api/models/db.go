package models

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
	LogMode     bool
	MaxIdleConn int
	MaxOpenConn int
	MaxLifetime int
}

type User struct {
	gorm.Model
	Orders []Order
	Data   string `sql:"type:jsonb not null default '{}'::jsonb" json:"-"`
}

type Order struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignkey:UserID"`
	Data   string `sql:"type:jsonb not null default '{}'::jsonb"`
}

func (User) TableName() string {
	return "user"
}

func (Order) TableName() string {
	return "order"
}

func LoadConfig() (*DBConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Note: .env file not found, using environment variables: %v", err)
	}

	config := &DBConfig{
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        getEnv("DB_PORT", "5432"),
		User:        getEnv("DB_USER", "postgre"),
		Password:    getEnv("DB_PASSWORD", "passw000rd"),
		DBName:      getEnv("DB_NAME", "gowebapp"),
		SSLMode:     getEnv("DB_SSL_MODE", "disable"),
		LogMode:     getEnvBool("GORM_LOG_MODE", true),
		MaxIdleConn: getEnvInt("GORM_MAX_IDLE_CONNS", 10),
		MaxOpenConn: getEnvInt("GORM_MAX_OPEN_CONNS", 100),
		MaxLifetime: getEnvInt("GORM_CONN_MAX_LIFETIME", 3600),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *DBConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func (c *DBConfig) ConnectionURI() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func InitDB() (*gorm.DB, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	db, err := gorm.Open("postgres", config.ConnectionURI())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db.LogMode(config.LogMode)
	
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Connected to database: %s@%s:%s/%s",
		config.User, config.Host, config.Port, config.DBName)

	err = setupTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup tables: %v", err)
	}

	return db, nil
}

func setupTables(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	tables := []interface{}{&User{}, &Order{}}
	
	for _, table := range tables {
		if !db.HasTable(table) {
			err := db.CreateTable(table).Error
			if err != nil {
				return fmt.Errorf("failed to create table: %v", err)
			}
			log.Printf("Table created: %T", table)
		}
	}

	err := db.AutoMigrate(tables...).Error
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %v", err)
	}

	if !db.HasTable(&Order{}) {
		err := db.Model(&Order{}).AddForeignKey("user_id", "\"user\"(id)", "RESTRICT", "RESTRICT").Error
		if err != nil {
			log.Printf("Note: Could not add foreign key (might already exist): %v", err)
		}
	}

	log.Println("Database tables setup completed")
	return nil
}

func CloseDB(db *gorm.DB) {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}