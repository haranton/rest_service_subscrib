package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool

func init() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Note: .env file not found: %v\n", err)
	}
}

func ConnectDB() (*pgxpool.Pool, error) {
	// Получаем параметры подключения
	host := getEnv("DB_HOST", "localhost")
	portStr := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "postgres")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Преобразуем порт в число
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid DB_PORT value '%s', using default 5432\n", portStr)
		port = 5432
	}

	// Формируем DSN строку
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode,
	)

	// Настройка пула соединений
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DB config: %w", err)
	}

	// Создаем пул соединений
	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Printf("Successfully connected to database %s@%s:%d\n", dbname, host, port)
	return pool, nil
}

func CloseDB() {
	if pool != nil {
		pool.Close()
		log.Println("Database connection closed")
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
