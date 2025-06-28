package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// IsolationLevel defines the isolation levels for database transactions.
type IsolationLevel int

const (
	// DefaultIsolationLevel uses the default isolation level.
	DefaultIsolationLevel IsolationLevel = iota
	// ReadUncommittedIsolationLevel allows dirty reads.
	ReadUncommittedIsolationLevel
	// ReadCommittedIsolationLevel ensures a transaction only reads committed data.
	ReadCommittedIsolationLevel
	// WriteCommittedIsolationLevel is typically not used in most databases.
	WriteCommittedIsolationLevel
	// RepeatableReadIsolationLevel ensures consistent reads during a transaction.
	RepeatableReadIsolationLevel
	// SnapshotIsolationLevel provides snapshot isolation.
	SnapshotIsolationLevel
	// SerializableIsolationLevel is the strictest isolation level.
	SerializableIsolationLevel
	// LinearizableIsolationLevel is the strongest isolation level with total order of all operations.
	LinearizableIsolationLevel
)

// TxOptions represents the options for a transaction.
type TxOptions struct {
	ReadOnly  bool
	Isolation IsolationLevel
}

// TransactionFunction defines the function type for operations within a transaction.
// The function receives a context and a transaction object, and returns an error if any operation fails.
type TransactionFunction func(ctx context.Context, tx *sqlx.Tx) error

// NewTransactionManager creates a new TransactionManager for the Client.
// The manager handles the lifecycle of a transaction, including committing or rolling back
// depending on whether the transaction function returns an error or a panic occurs.
//
// Returns:
//   - TransactionManager: A function that manages a transaction, executing the given transaction function.
func (sql *Client) Transaction(ctx context.Context, txFn TransactionFunction, opts ...TxOptions) (err error) {
	var txOptions *TxOptions
	if len(opts) > 0 {
		txOptions = &opts[0]
	} else {
		txOptions = nil
	}

	tx, err := sql.conn.BeginTxx(ctx, convertTxOptions(txOptions))
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		r := recover()
		if r == nil && err == nil {
			return
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("roll back transaction: %w, err: %w", rollbackErr, err)
		}

		if r != nil {
			err = fmt.Errorf("panic: %v, err: %w", r, err)
		}
	}()

	err = txFn(ctx, tx)
	if err != nil {
		return fmt.Errorf("run transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func convertTxOptions(txOptions *TxOptions) *sql.TxOptions {
	if txOptions == nil {
		return &sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		}
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
		fallthrough
	default:
		level = sql.LevelDefault
	}

	return &sql.TxOptions{
		Isolation: level,
		ReadOnly:  txOptions.ReadOnly,
	}
}
