package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConectarDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Faltan variables de entorno para la conexión a la base de datos")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Error,
			Colorful: true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error al obtener la conexion de la base de datos: %v", err)
	}

	sqlDB.SetMaxOpenConns(5)                  // Número máximo de conexiones abiertas
	sqlDB.SetMaxIdleConns(2)                  // Núme	ro máximo de conexiones inactivas
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // Tiempo máximo de vida de una conexión

	return db, nil
}
