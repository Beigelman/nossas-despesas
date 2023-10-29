package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type IsolationLevel int

const (
	DefaultIsolationLevel IsolationLevel = iota
	ReadUncommittedIsolationLevel
	ReadCommittedIsolationLevel
	WriteCommittedIsolationLevel
	RepeatableReadIsolationLevel
	SnapshotIsolationLevel
	SerializableIsolationLevel
	LinearizableIsolationLevel
)

type TxOptions struct {
	ReadOnly  bool
	Isolation IsolationLevel
}

type TransactionFunction func(ctx context.Context, tx *sqlx.Tx) error

type TransactionManager func(ctx context.Context, txFn TransactionFunction, ops ...TxOptions) error

func (sql *SQLDatabase) NewTransactionManager() TransactionManager {
	return func(ctx context.Context, txFn TransactionFunction, opts ...TxOptions) error {
		var txOptions *TxOptions
		if len(opts) > 0 {
			txOptions = &opts[0]
		} else {
			txOptions = nil
		}

		tx, err := sql.db.BeginTxx(ctx, convertTxOptions(txOptions))
		if err != nil {
			return fmt.Errorf("BeginTxx: %w", err)
		}

		defer func() {
			r := recover()
			if r == nil && err == nil {
				return
			}

			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("tx.Rollback: %w, err: %w", rollbackErr, err)
			}

			if r != nil {
				err = fmt.Errorf("panic: %v, err: %w", r, err)
			}
		}()

		err = txFn(ctx, tx)
		if err != nil {
			return fmt.Errorf("UnitOfWorkFn: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("tx.Commit: %w", err)
		}

		return nil
	}
}

func convertTxOptions(txOptions *TxOptions) *sql.TxOptions {
	if txOptions == nil {
		return nil
	}

	var level sql.IsolationLevel

	switch txOptions.Isolation {
	case ReadUncommittedIsolationLevel:
		level = sql.LevelReadUncommitted
	case ReadCommittedIsolationLevel:
		level = sql.LevelReadCommitted
	case WriteCommittedIsolationLevel:
		level = sql.LevelWriteCommitted
	case RepeatableReadIsolationLevel:
		level = sql.LevelRepeatableRead
	case SnapshotIsolationLevel:
		level = sql.LevelSnapshot
	case SerializableIsolationLevel:
		level = sql.LevelSerializable
	case LinearizableIsolationLevel:
		level = sql.LevelLinearizable
	case DefaultIsolationLevel:
	default:
		level = sql.LevelDefault
	}

	return &sql.TxOptions{
		Isolation: level,
		ReadOnly:  txOptions.ReadOnly,
	}
}
