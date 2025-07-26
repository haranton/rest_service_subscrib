package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() error {
	// 1. Получаем абсолютный путь
	absPath := "C:\\Users\\anton\\Documents\\godev\\rest_service\\internal\\db\\migrations"

	// 2. Форматируем путь для Windows
	migrationsPath := "file://" + strings.ReplaceAll(absPath, "\\", "/")

	// 3. Формируем DSN для PostgreSQL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 4. Логируем для отладки
	log.Printf("Migrations path: %s", migrationsPath)

	// 5. Инициализируем мигратор
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}
	defer m.Close()

	// 6. Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
