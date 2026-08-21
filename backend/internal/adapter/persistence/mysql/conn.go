// Package mysql is the MySQL persistence adapter. It implements the
// item.Repository port using sqlx.
package mysql

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"
)

// Config holds MySQL connection settings.
type Config struct {
	Host     string `mapstructure:"MYSQL_HOST"`
	Port     int    `mapstructure:"MYSQL_PORT"`
	User     string `mapstructure:"MYSQL_USER"`
	Password string `mapstructure:"MYSQL_PASSWORD"`
	DB       string `mapstructure:"MYSQL_DB"`
}

// Connect opens a sqlx.DB and verifies it with a Ping.
func Connect(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql: open: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql: ping: %w", err)
	}
	return db, nil
}
