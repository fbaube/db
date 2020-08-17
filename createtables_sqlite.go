package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	S "strings"
)

// CreateTable_sqlite creates a table for our simplified DB model where
// all columns can be either string or int. Note that all arguments except
// the first can be nil.
func (pDB *MmmcDB) CreateTable_sqlite(ts TableSpec) {
	var CTS string // the Create Table string
	var hasFKs, hasInts, hasStrs bool
	hasFKs = (ts.forenKeys != nil && len(ts.forenKeys) > 0)
	hasInts = (ts.intFields != nil && len(ts.intFields) > 0)
	hasStrs = (ts.strFields != nil && len(ts.strFields) > 0)

	CTS = "CREATE TABLE " + ts.tableName + "(\n"
	CTS += "idx_" + ts.tableName + " integer not null primary key, "
	CTS += "-- NOTE: integer, not int \n"
	if hasFKs {
		for _, tbl := range ts.forenKeys {
			// idx_inb integer not null references INB,
			CTS += "idx_" + tbl + " integer " + /* not null */ "references " + tbl + ", \n"
		}
	}
	if hasInts {
		for i, fld := range ts.intFields {
			// filect int not null check (filect >= 0) default 0
			CTS += fld + " int not null"
			switch ts.intRanges[i] {
			case 1:
				// check (filect >= 0)
				CTS += " check (" + fld + " > 0), \n"
			case 0:
				CTS += " check (" + fld + " >= 0), \n"
			default: // case -1:
				CTS += ", \n"
			}
		}
	}
	if hasStrs {
		for _, fld := range ts.strFields {
			// filect int not null check (filect >= 0) default 0
			CTS += fld + " text not null, \n"
		}
	}
	if hasFKs {
		// FOREIGN KEY(idx_inb) REFERENCES INB(idx_inb)
		for _, tbl := range ts.forenKeys {
			// idx_inb integer not null references INB,
			// TMP := "foreign key(idx_" + tbl + ") references " + tbl + "(idx_" + tbl + "), \n"
			// println("TMP:", TMP)
			CTS += "foreign key(idx_" + tbl + ") references " + tbl + "(idx_" + tbl + "), \n"
		}
	}
	CTS = S.TrimSuffix(CTS, "\n")
	CTS = S.TrimSuffix(CTS, " ")
	CTS = S.TrimSuffix(CTS, ",")
	CTS += "\n);"
	println("= = = = = = = = \n", CTS, "= = = = = = = =")
	pDB.DB.MustExec(CTS)
	pDB.DumpTableSchema_sqlite(ts.tableName, os.Stdout)
	println("TODO: Insert record with IDX 0 and string descr's")
	println("TODO: Dump all table records (i.e. just one)")
}

func (pDB *MmmcDB) DumpTableSchema_sqlite(tableName string, w io.Writer) {
	println("TODO: Dump table structure:", tableName)
	var R *sql.Rows
	var e error
	var CT []*sql.ColumnType
	// func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	// func (rs *Rows) ColumnTypes() ([]*ColumnType, error)
	R, e = pDB.DB.Query("select * from " + tableName)
	if e != nil {
		panic("DB: sqlite1: " + e.Error())
	}
	CT, e = R.ColumnTypes()
	if e != nil {
		panic("DB: sqlite2: " + e.Error())
	}
	for i, ct := range CT {
		L, _ := ct.Length()
		st := ct.ScanType()
		sst := ""
		if st != nil {
			sst = fmt.Sprintf("scantp<%+v> ", st)
		}
		fmt.Printf("col[%d]: %s \t dbtpnm<%s> len<%d> %s \n",
			i, ct.Name(), ct.DatabaseTypeName(), L, sst)
	}
}
