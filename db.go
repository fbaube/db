package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	FP "path/filepath"
	"strconv"
	S "strings"
	"time"
	MU "github.com/fbaube/miscutils"
	FU "github.com/fbaube/fileutils"
	"github.com/jmoiron/sqlx"
	// to get init()
	_ "github.com/mattn/go-sqlite3"
)

// sqlite3 my_database.sq3 ".backup 'backup_file.sq3'"
// sqlite3 m_database.sq3 ".backup m_database.sq3.bak"
// sqlite3 my_database .backup > my_database.back

// https://github.com/golang/go/wiki/SQLInterface
// ExecContext is used for queries where no rows are returned:
// QueryContext is used for retrieving rows.
// QueryRowContext is used where only a single row is expected.
// If a DB column is nullable, pass a type supporting null values to Scan.
// Only NullBool, NullFloat64, NullInt64, NullString are in database/sql.
// Implementations of DB-specific null types are left to the DB driver.

const DBNAME = "mmmc.db"

type MmmcDB struct {
	DirrPath FU.BasicPath
	FilePath FU.BasicPath
	// Connection
	// theSqlDB *sql.DB
	// Connection
	theSqlxDB *sqlx.DB
	// Session-level open Txn (if non-nil)
	TheSqlxTxn *sqlx.Tx
}

func (p *MmmcDB) DBX() *sqlx.DB {
	return p.theSqlxDB
}

// theDB IS A SINGLETON. If it is non-nil then assume that all
// three of its path-related fields are properly initialized.
var theDB *MmmcDB

var stmt *sql.Stmt
var rslt *sql.Result
var e error

// NewlyConfiguredMmmcDB initializes path-related variables but does not do more.
// argpath can be a relative path passed to the CLI.
// Unlike how some other filepath arguments are handled,
// if the DB path is "", it is set to the CWD (current
// working directory).
func NewlyConfiguredMmmcDB(argpath string) (*MmmcDB, error) {
	pDB := new(MmmcDB)
	if argpath == "" {
		argpath = "."
	}
	// If the DB name was accidentally provided, trim it off to prevent problems.
	var relFP string
	relFP = S.TrimSuffix(argpath, DBNAME)

	pDB.DirrPath = *FU.NewBasicPath(relFP)
	if !pDB.DirrPath.IsOkayDir() { // PathType() != "DIR" {
		retErr := MU.TracedError(fmt.Errorf("DB dir not exist or not a dir: %s",
		  pDB.DirrPath.AbsFilePath.S()))
		return nil, retErr
	}
	pDB.FilePath = *FU.NewBasicPath(FP.Join(relFP, "mmmc.db"))
	theDB = pDB
	return pDB, nil
}

func checkerr(e error) {
	if e == nil {
		return
	}
	panic("Sqlite3 FAILURE: " + e.Error())
}

func mustExecStmt(s string) {
	stmt, e = theDB.theSqlxDB.Prepare(s)
	checkerr(e)
	/* rslt, */ _, e := stmt.Exec()
	checkerr(e)
	// liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("DD:mustExecStmt: ID %d nR %d \n", liid, naff)
}

// MustCreateTable makes sure it exists but
// does NOT drop an already-existing table.
func MustExistTable(s string) {
	mustExecStmt(s)
}

func Make09azStringLen1(i int) string {
	if i <= 9 {
		return strconv.Itoa(i)
	}
	var bb = make([]byte, 1, 1)
	bb[0] = byte(i - 10 + 'a')
	return string(bb)
}

func ComprestNowString() string {
	var now = time.Now()
	// year = last digit
	var y string = fmt.Sprintf("%d", now.Year())[3:]
	var m string = Make09azStringLen1(int(now.Month()))
	var d string = Make09azStringLen1(now.Day())
	var h string = Make09azStringLen1(now.Hour())
	var n string = Make09azStringLen1(now.Minute() / 2)
	// fmt.Printf("%s-%s-%s-%s-%s", y, m, d, h, n)
	return fmt.Sprintf("%s%s%s%s%s", y, m, d, h, n)
}

func ComprestYYYYMMstring() string {
	var now = time.Now()
	// year = last digit
	var y string = fmt.Sprintf("%d", now.Year())[3:]
	var m string = Make09azStringLen1(int(now.Month()))
	// fmt.Printf("%s-%s-%s-%s-%s", y, m, d, h, n)
	return fmt.Sprintf("%s%s", y, m)
}

// MoveCurrentToBackup makes a best effort but can fail if the
// backup destination is a directory or has a permissions problem.
// The current DB is renamed and therefore "disappears".
func (p *MmmcDB) MoveCurrentToBackup() error {
	if !p.FilePath.Exists {
		println("    --> No current DB to move to backup")
		return nil
	}
	// func Rename(oldpath, newpath string) error
	var cns = ComprestNowString()
	var fromFP string = p.FilePath.AbsFilePath.S()
	var toFP string = p.FilePath.AbsFilePathParts.BaseName + "-" + cns + ".db"
	e := os.Rename(fromFP, toFP)
	if e != nil {
		return fmt.Errorf("Can't move current DB to <%s>: %w: ", toFP, e)
	}
	println("    --> Old DB moved to:", toFP)
	return nil
}

// DupeCurrentToBackup makes a best effort but can fail if the
// backup destination is a directory or has a permissions problem.
// The current DB is not affected.
func (p *MmmcDB) DupeCurrentToBackup() error {
	if !p.FilePath.Exists {
		println("    --> No current DB to duplicate to backup")
		return nil
	}
	var cns = ComprestNowString()
	var fromFP string = p.FilePath.AbsFilePath.S()
	var toFP string = p.FilePath.BaseName + "-" + cns + ".db"

	e := FU.CopyFromTo(fromFP, toFP)
	if e != nil {
		return fmt.Errorf("Can't copy current DB to <%s>: %w: ", toFP, e)
	}
	println("    --> Old DB copied to backup at:", toFP)
	return nil
}

// ForceEmpty is a convenience function.
func (p *MmmcDB) ForceEmpty() {
	if theDB == nil {
		panic("db.forcempty.unitialized.L176")
	}
	p.MoveCurrentToBackup()
	p.ForceExist()
}

// ForceExist creates a new empty DB with the proper schema.
func (p *MmmcDB) ForceExist() {
	if theDB == nil {
		panic("db.forcexist.uninitd.L185")
	}
	var dest string = p.FilePath.AbsFilePath.S()
	// println("    --> Creating new DB at:", dest)
	var e error
	var theSqlDB *sql.DB

	theSqlDB, e = sql.Open("sqlite3", dest)
	// loggerAdapter := zerologadapter.New(zerolog.New(zerolog.NewConsoleWriter()))
	// theSqlDB = sqldblogger.OpenDriver(dest, &sqlite3.SQLiteDriver{}, loggerAdapter /*, ...options */)

	checkerr(e)
	e = theSqlDB.Ping()
	checkerr(e)
	e = theSqlDB.PingContext(context.Background())
	checkerr(e)
	println("    --> New DB created at:", dest)
	drivers := sql.Drivers()
	fmt.Printf("    --> DB driver(s): %+v \n", drivers)
	theDB.theSqlxDB = sqlx.NewDb(theSqlDB, "sqlite3")

	for _, s := range schemasALL {
		mustExecStmt("CREATE TABLE IF NOT EXISTS " + s)
	}
	mustExecStmt("PRAGMA foreign_keys = ON;")
}
