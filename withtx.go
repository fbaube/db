package dbutils

import "database/sql"

// TxFunc is a function that can be executed inside a transaction.
type TxFunc func(tx *sql.Tx) error

// WithTx wraps a fTxFunc unction "f" in an SQL transaction. After the function returns,
// the transaction is committed if there's no error, or rolled back if there is one.
func WithTx(db *sql.DB, f TxFunc) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	return f(tx)
}
