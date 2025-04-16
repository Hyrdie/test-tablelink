package app

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Service Configuration
	ServiceName string
	Env         string
	LogLevel    int

	// Database Configuration
	DBConnURI string

	// Redis Configuration
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Connection Pool Configuration
	DBMaxPoolSize        int
	DBMaxIdleConnections int
	DBMaxIdleTime        time.Duration
	DBMaxLifeTime        time.Duration

	// Server Configuration
	BindAddress string
}

var cfg *Config

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Debug: Print all environment variables
	log.Println("Loading environment variables...")
	for _, env := range os.Environ() {
		log.Println(env)
	}

	// Parse numeric values
	logLevel, err := strconv.Atoi(getEnv("LOG_LEVEL", "1"))
	if err != nil {
		return nil, fmt.Errorf("invalid LOG_LEVEL: %w", err)
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	maxPoolSize, err := strconv.Atoi(getEnv("PG_MAX_POOL_SZE", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_POOL_SZE: %w", err)
	}

	maxIdleConnections, err := strconv.Atoi(getEnv("PG_MAX_IDLE_CONNECTIONS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_IDLE_CONNECTIONS: %w", err)
	}

	maxIdleTime, err := time.ParseDuration(getEnv("PG_MAX_IDLE_TIME", "5m"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_IDLE_TIME: %w", err)
	}

	maxLifeTime, err := time.ParseDuration(getEnv("PG_MAX_LIFE_TIME", "1h"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_LIFE_TIME: %w", err)
	}

	// Construct database connection URI from individual environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "test-tablelink")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")

	dbConnURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	cfg = &Config{
		// Service Configuration
		ServiceName: getEnv("SERVICE_NAME", "test-tablelink"),
		Env:         getEnv("ENV", "development"),
		LogLevel:    logLevel,

		// Database Configuration
		DBConnURI: dbConnURI,

		// Redis Configuration
		RedisAddr:     getEnv("REDIS_HOST", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		// Connection Pool Configuration
		DBMaxPoolSize:        maxPoolSize,
		DBMaxIdleConnections: maxIdleConnections,
		DBMaxIdleTime:        maxIdleTime,
		DBMaxLifeTime:        maxLifeTime,

		// Server Configuration
		BindAddress: getEnv("BIND_ADDRESS", ":8080"),
	}

	// Debug: Print loaded configuration
	log.Printf("Loaded configuration: %+v\n", cfg)

	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
