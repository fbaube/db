package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	FP "path/filepath"
	S "strings"
	MU "github.com/fbaube/miscutils"
	FU "github.com/fbaube/fileutils"
	"github.com/jmoiron/sqlx"
	// to get init()
	_ "github.com/mattn/go-sqlite3"
)

// Times has (create, import, last edit) and uses only ISO-8601 / RFC 3339.
type Times struct {
	Created  string
	Imported string
	Edited   string
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
	FU.BasicPath
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

// NewMmmcDB initializes path-related variables but does not do more.
// argpath can be a relative path passed to the CLI;
// if it is "", the DB path is set to the CWD (current
// working directory).
func NewMmmcDB(argpath string) (*MmmcDB, error) {
	pDB := new(MmmcDB)
	if argpath == "" {
		var e error
		argpath, e = os.Getwd() // "."
		if e != nil {
			panic(e)
		}
	}
	// If the DB name was accidentally provided, trim it off to prevent problems.
	var relFP string
	relFP = S.TrimSuffix(argpath, DBNAME)

	// pDB.DirrPath = *FU.NewBasicPath(relFP)
	// if !pDB.DirrPath.IsOkayDir() { // PathType() != "DIR" {
	// dp := FU.NewBasicPath(pDB.BasicPath.AbsFilePathParts.DirPath.S())
	dp := FU.NewBasicPath(FP.Dir(pDB.BasicPath.S()))
	if !dp.IsOkayDir() {
		retErr := MU.TracedError(fmt.Errorf("DB dir not exist or not a dir: %s", dp))
		return nil, retErr
	}
	pDB.BasicPath = *FU.NewBasicPath(FP.Join(relFP, "mmmc.db"))
	theDB = pDB
	return pDB, nil
}

// ForceExistDBandTables creates a new empty DB with the proper schema.
func (p *MmmcDB) ForceExistDBandTables() {
	if theDB == nil {
		panic("db.forcexist.uninitd.L95")
	}
	var dest string = p.BasicPath.AbsFilePath.S()
	// println("    --> Creating new DB at:", dest)
	var e error
	var theSqlDB *sql.DB

	theSqlDB, e = sql.Open("sqlite3", dest)
	// loggerAdapter := zerologadapter.New(zerolog.New(zerolog.NewConsoleWriter()))

	checkerr(e)
	e = theSqlDB.Ping()
	checkerr(e)
	e = theSqlDB.PingContext(context.Background())
	checkerr(e)
	println("    --> New DB created at:", dest)
	drivers := sql.Drivers()
	fmt.Printf("    --> DB driver(s): %+v \n", drivers)
	theDB.DB = sqlx.NewDb(theSqlDB, "sqlite3")

	for _, s := range schemasALL {
		mustExecStmt("CREATE TABLE IF NOT EXISTS " + s)
	}
	// p.CreateTable_Inbatch_sqlite()
	// p.CreateTable_Content_sqlite()
	p.CreateTable_sqlite(TableSpec_Inbatch)
	p.CreateTable_sqlite(TableSpec_Content)

	// It seems weird that this is necessary, but cos of some retro compatibility,
	// SQLite does not by default enforce foreign key constraints. 
	mustExecStmt("PRAGMA foreign_keys = ON;")
}
