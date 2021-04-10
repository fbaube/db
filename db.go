package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	WU "github.com/fbaube/wasmutils"
	"github.com/jmoiron/sqlx"

	// to get init()
	_ "github.com/mattn/go-sqlite3"
)

// Times has (create, import, last edit) and uses only ISO-8601 / RFC 3339.
type Times struct {
	T_Cre string
	T_Imp string
	T_Edt string
}

// At the CLI:
// sqlite3 my_database.sq3 ".backup 'backup_file.sq3'"
// sqlite3 m_database.sq3 ".backup m_database.sq3.bak"
// sqlite3 my_database .backup > my_database.back

// https://github.com/golang/go/wiki/SQLInterface
// ExecContext is used when no rows are returned;
// QueryContext is used for retrieving rows.
// QueryRowContext is used where only a single row is expected.
// If a DB column is nullable, pass a type supporting null values to Scan.
// Only NullBool, NullFloat64, NullInt64, NullString are in database/sql.
// Implementations of DB-specific null types are left to the DB driver.

// DBNAME defaults to "mmmc.db"
const DBNAME = "mmmc.db"

// MmmcDB stores DB filepaths, DB cnxns, DB txns.
type MmmcDB struct {
	FU.PathProps
	// Connection
	*sqlx.DB
	// Session-level open Txn (if non-nil). Relevant API:
	// func (db *sqlx.DB)    Beginx() (*sqlx.Tx, error)
	// func (db *sqlx.DB) MustBegin()  *sqlx.Tx
	// func (tx *sql.Tx)     Commit()   error
	// func (tx *sql.Tx)   Rollback()   error
	*sqlx.Tx
}

// theDB IS A SINGLETON. If it is non-nil then assume that all
// three of its path-related fields are properly initialized.
var theDB *MmmcDB

// NOTE These are probably disastrous in multithreaded use.
var stmt *sql.Stmt
var rslt *sql.Result
var e error

// NewMmmcDB does not open a DB, it merely checks that the given path is OK.
// That is to say, it initializes path-related variables but does not do more.
// argpath can be a relative path passed to the CLI; if it is "", the DB path
// is set to the CWD (current working directory).
func NewMmmcDB(argpath string) (*MmmcDB, error) {
	var relFP = argpath
	if argpath != "" {
		// If the DB name was unnecessarily provided,
		// trim it off to prevent problems.
		relFP = S.TrimSuffix(argpath, DBNAME)
	} else {
		var e error
		relFP, e = os.Getwd() // "."
		if e != nil {
			if WU.IsWasm() {
				L.L.Warning("FIXME: Where is DB in browser WASM ?")
			}
			L.L.Error("DB: can't get CWD: %w", e)
			os.Exit(1)
		}
	}
	pDB := new(MmmcDB)
	dp := FU.NewPathProps(FP.Dir(pDB.PathProps.AbsFP.S()))
	if !dp.IsOkayDir() {
		retErr := "DB dir not exist or not a dir: " + dp.String()
		L.L.Error(retErr)
		return nil, errors.New(retErr)
	}
	pDB.PathProps = *FU.NewPathProps(FP.Join(relFP, "mmmc.db"))
	theDB = pDB
	return pDB, nil
}

// ForceExistDBandTables creates a new empty DB with the proper schema.
func (p *MmmcDB) ForceExistDBandTables() {
	if theDB == nil {
		L.L.Panic("theDB does not exist yet")
	}
	var dest string = p.PathProps.AbsFP.S()
	var e error
	var theSqlDB *sql.DB

	theSqlDB, e = sql.Open("sqlite3", dest)
	checkerr(e)
	e = theSqlDB.Ping()
	checkerr(e)
	e = theSqlDB.PingContext(context.Background())
	checkerr(e)
	L.L.Okay("New DB created at: " + FU.Tildotted(dest))
	drivers := sql.Drivers()
	L.L.Info("DB driver(s): %+v", drivers)
	theDB.DB = sqlx.NewDb(theSqlDB, "sqlite3")

	for _, cfg := range AllTableConfigs {
		p.CreateTable_sqlite(cfg)
	}
	// It may seem odd that this is necessary,
	// but for some retro compatibility, SQLite does
	// not by default enforce foreign key constraints.
	mustExecStmt("PRAGMA foreign_keys = ON;")
}

func (p *MmmcDB) Verify() {
	mustExecStmt("PRAGMA integrity_check;")
	mustExecStmt("PRAGMA foreign_key_check;")
}
