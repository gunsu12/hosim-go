package database

import (
	"context"

	"gorm.io/gorm"
)

type txContextKey struct{}

var transactionContextKey = txContextKey{}

// TransactionManager adalah abstraksi boundary transaksi database untuk application layer.
type TransactionManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type gormTransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager membuat implementasi TransactionManager berbasis GORM.
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &gormTransactionManager{db: db}
}

// WithTx mengeksekusi fungsi fn dalam transaksi database atomik.
// Jika fn mengembalikan error, transaksi di-rollback secara otomatis.
// Jika fn selesai tanpa error, transaksi di-commit.
func (tm *gormTransactionManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, transactionContextKey, tx)
		return fn(txCtx)
	})
}

// ContextWithTx memasukkan handle transaksi GORM ke dalam context.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionContextKey, tx)
}

// GetTxFromContext mengambil handle transaksi GORM dari context jika ada.
func GetTxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return nil
}

// DBFromContextOr mengembalikan DB dari context transaksi bila ada, atau fallback ke defaultDB.
func DBFromContextOr(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx := GetTxFromContext(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return defaultDB.WithContext(ctx)
}
