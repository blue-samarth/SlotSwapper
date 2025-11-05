package utils

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TransactionOptions struct {
	Timeout        time.Duration
	IsolationLevel string
	// MaxRetries      int
	// RetryDelay     time.Duration
}

var DefaultTransactionOptions = TransactionOptions{
	Timeout:        5 * time.Second,
	IsolationLevel: "READ COMMITTED",
	// MaxRetries:     3,
	// RetryDelay:     100 * time.Millisecond,
}

func WithTransaction(db *gorm.DB, options TransactionOptions, fn func(*gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	if err := tx.Exec("SET LOCAL lock_timeout = '5s'").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set lock timeout: %w", err)
	}

	if options.IsolationLevel != "" {
		if err := tx.Exec(fmt.Sprintf("SET LOCAL transaction_isolation = '%s'", options.IsolationLevel)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to set isolation level: %w", err)
		}
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func WithDefaultTransaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	return WithTransaction(db, DefaultTransactionOptions, fn)
}
