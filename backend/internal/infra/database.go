package infra

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/moLIart/gomoku-backend/pkg/errorx"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	dataSource string
	conn       *sqlx.DB
}

func NewDatabase(dataSource string) *Database {
	return &Database{
		dataSource: dataSource,
		conn:       nil,
	}
}

func (s *Database) AcquireConn() (*sqlx.DB, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("database connection closed")
	}
	return s.conn, nil
}

func (s *Database) Start(ctx context.Context) error {
	log.Info("Starting database...")

	// Open a new database connection
	conn, err := sqlx.Open("postgres", s.dataSource)
	if err != nil {
		return errorx.Wrap(err, "database connection open")
	}

	s.conn = conn

	// Ping the database to ensure the connection is valid
	// TODO: Use a more robust health check
	err = s.conn.PingContext(ctx)
	if err != nil {
		return errorx.Wrap(err, "ping database")
	}

	return nil
}

func (s *Database) Stop() {
	log.Info("Stopping database...")

	if err := s.conn.Close(); err != nil {
		log.Error("error closing database connection")
	}

	s.conn = nil
	log.Info("database connection closed")
}
