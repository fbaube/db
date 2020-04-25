package db

import (
  "fmt"
  "io"
  "os"
  "database/sql"
  S "strings"
)

func (pDB *MmmcDB) CreateTable_sqlite(
	tableName string,
	forenKeys []string,
	intFields []string,
	intRanges []int,    // to DB[0] // -1, 0, 1
	strFields []string,
	strDescrs []string) { // to DB[0]
	/*
	var schemaINB string = `INB(
 idx_inb integer not null primary key, -- NOTE: "integer", not "int"
relfilepath text not null,
absfilepath text not null,
     creatime text not null default (datetime('now')), -- UTC ISO-8601
       desc text not null default "No load message", -- from CLI commit msg
     filect int  not null check (filect >= 0) default 0  )`
	*/
	var CTS string // the Create Table string
	var hasFKs, hasInts, hasStrs bool
	hasFKs  = (forenKeys != nil && len(forenKeys) > 0)
	hasInts = (intFields != nil && len(intFields) > 0)
	hasStrs = (strFields != nil && len(strFields) > 0)

	CTS = "CREATE TABLE " + tableName + "(\n"
	CTS += "idx_" + tableName + " integer not null primary key, "
	CTS += "-- NOTE: integer, not int \n"
	if hasFKs {
		for _, tbl := range forenKeys {
			// idx_inb integer not null references INB,
			CTS += "idx_" + tbl + " integer not null references " + tbl + ", \n"
		}
	}
	if hasInts {
		for i, fld := range intFields {
			// filect int not null check (filect >= 0) default 0
			CTS += fld + " int not null"
			switch intRanges[i] {
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
		for _, fld := range strFields {
			// filect int not null check (filect >= 0) default 0
			CTS += fld + " text not null, \n"
		}
	}
	if hasFKs {
		// FOREIGN KEY(idx_inb) REFERENCES INB(idx_inb)
		for _, tbl := range forenKeys {
			// idx_inb integer not null references INB,
			CTS += "foreign key(idx_" + tbl + ") references " + tbl + "(idx_" + tbl + "), \n"
		}
	}
	CTS = S.TrimSuffix(CTS, "\n")
	CTS = S.TrimSuffix(CTS, " ")
	CTS = S.TrimSuffix(CTS, ",")
	CTS += "\n)"
	println("= = = = = = = = \n", CTS, "= = = = = = = =")
	pDB.theSqlxDB.MustExec(CTS)
  pDB.DumpTableSchema_sqlite(tableName, os.Stdout)
	println("TODO: Insert record with IDX 0 and string descr's")
  println("TODO: Dump all table records (i.e. just one)")
}

func (pDB *MmmcDB) DumpTableSchema_sqlite(tableName string, w io.Writer) {
  println("TODO: Dump table structure:", tableName)
  var R *sql.Rows
  var e error
  var C []*sql.ColumnType
  // func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
  // func (rs *Rows) ColumnTypes() ([]*ColumnType, error)
  R,e = pDB.theSqlxDB.Query("select * from " + tableName)
  if e != nil { panic("DB: sqlite1: " + e.Error())}
  C,e = R.ColumnTypes()
  if e != nil { panic("DB: sqlite2: " + e.Error())}
  for i, c := range C {
    L,_ := c.Length()
    fmt.Printf("col[%d]: nm<%s> dbtp<%s> len<%d> gotp<%s> \n",
      i, c.Name(), c.DatabaseTypeName(), L, c.ScanType())
  }
}
