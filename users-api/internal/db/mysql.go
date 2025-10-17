// internal/db/mysql.go
package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config desacoplado del resto del microservicio.
// Adaptá tu internal/config para construir este struct y pasarlo a New().
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// New crea una conexión gorm a MySQL usando la configuración provista.
// Devuelve *gorm.DB y el error si ocurre.
func New(cfg Config) (*gorm.DB, error) {
	// Construir DSN con parámetros razonables por defecto
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Ajustes del pool de conexiones
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	// valores por defecto razonables; pueden exponerse vía Config si se desea
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// Parámetros extra del DSN (se agregan al final).
