package repositories

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/moLIart/gomoku-backend/internal/infra"
	"github.com/moLIart/gomoku-backend/pkg/errorx"
)

type UnitOfWork struct {
	db *infra.Database
	tx *sqlx.Tx
}

func NewUnitOfWork(db *infra.Database) *UnitOfWork {
	return &UnitOfWork{
		db: db,
	}
}

func (uow *UnitOfWork) GetPlayerRepository() *PlayerRepository {
	return NewPlayerRepository(uow.tx)
}

func (uow *UnitOfWork) GetGameRepository() *GameRepository {
	return NewGameRepository(uow.tx)
}

func (uow *UnitOfWork) Begin(ctx context.Context) error {
	conn, err := uow.db.AcquireConn()
	if err != nil {
		return errorx.Wrap(err, "acquire db connection")
	}

	uow.tx, err = conn.BeginTxx(ctx, nil)
	if err != nil {
		return errorx.Wrap(err, "begin transaction")
	}

	return nil
}

func (uow *UnitOfWork) Complete(err error) error {
	if err != nil {
		if rbErr := uow.tx.Rollback(); rbErr != nil {
			if rbErr == sql.ErrTxDone {
				return errorx.Wrap(err, "detected error but transaction is already completed")
			}
			return errorx.Wrap(rbErr, "detected error but rollback failed for some reason")
		}
		return errorx.Wrap(err, "detected error and rolled transaction back")
	}

	if err := uow.tx.Commit(); err != nil {
		if err == sql.ErrTxDone {
			return nil
		}
		return errorx.Wrap(err, "commit failed due to error")
	}

	return nil

}
