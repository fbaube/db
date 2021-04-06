package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	S "strings"

	SB "github.com/fbaube/semblox"
)

// CreateTable_sqlite creates a table for our simplified SQLite DB model where
// - Every column is either string ("TEXT") or int ("INTEGER"),
// - Every column is NOT NULL,
// - Every column has type checking (TBS), and
// - Every table has a primary index field, and
// - Every index (primary and foreign) includes the full name of the table,
// which simplifies column creation and cross-referencing (including JOINs).
//
func (pDB *MmmcDB) CreateTable_sqlite(ts TableConfig) {
	var CTS string // the Create Table SQL string
	var hasFKs bool
	hasFKs = (ts.ForenKeys != nil && len(ts.ForenKeys) > 0)

	// === CREATE TABLE
	CTS = "CREATE TABLE " + ts.TableName + "(\n"
	// == PRIMARY KEY
	CTS += "idx_" + ts.TableName + " integer not null primary key autoincrement, "
	CTS += "-- NOTE: integer, not int. \n"
	if hasFKs {
		// === FOREIGN KEYS
		for _, tbl := range ts.ForenKeys {
			// idx_inb integer not null references INB,
			CTS += "idx_" + tbl + " integer " + /* not null */ "references " + tbl + ", \n"
		}
	}
	for _, fld := range ts.Columns {
		switch fld.TxtIntKeyEtc {
		case SB.D_INT:
			// e.g.: filect int not null check (filect >= 0) default 0
			// also: `Col1 INTEGER CHECK (typeof(Col1) == 'integer')`
			//
			CTS += fld.Code + " int not null"
			// CTS += fld.Code + " int not null check (typeof(" + fld.Code + ") == 'int')"
			/*
				switch ts.intRanges[i] {
				case 1:
					// check (filect >= 0)
					CTS += " check (" + fld + " > 0), \n"
				case 0:
					CTS += " check (" + fld + " >= 0), \n"
				default: // case -1:
					CTS += ", \n"
				}
			*/
			CTS += ", \n"
		case SB.D_TXT:
			CTS += fld.Code + " text not null check (typeof(" + fld.Code + ") == 'text'), \n"
		default:
			panic("Unhandled: " + fld.TxtIntKeyEtc)
		}
	}
	if hasFKs {
		// FOREIGN KEY(idx_inb) REFERENCES INB(idx_inb)
		for _, tbl := range ts.ForenKeys {
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
	pDB.DumpTableSchema_sqlite(ts.TableName, os.Stdout)
	println("TODO: Insert record with IDX 0 and string descr's")
	println("TODO: Dump all table records (i.e. just one)")
}

func (pDB *MmmcDB) DbTblColsIRL(tableName string) []*DbColIRL {
	if tableName == "" {
		return nil
	}
	var e error
	var Rs *sql.Rows
	var CTs []*sql.ColumnType
	var retval []*DbColIRL

	Rs, e = pDB.DB.Query("select * from " + tableName + " limit 1")
	if e != nil {
		println("DB select * failed on table", tableName, ":", e.Error())
		return nil
	}
	CTs, e = Rs.ColumnTypes()
	if e != nil {
		println("DB.ColumnTypes failed on table", tableName, ":", e.Error())
	}
	for _, ct := range CTs {
		dci := new(DbColIRL)
		dci.TxtIntKeyEtc = SB.TxtIntKeyEtc(ct.DatabaseTypeName())
		dci.Code = ct.Name()
		retval = append(retval, dci)
	}
	return retval
}

func (pDB *MmmcDB) DumpTableSchema_sqlite(tableName string, w io.Writer) {

	var theCols []*DbColIRL
	theCols = pDB.DbTblColsIRL(tableName)

	for i, c := range theCols {
		fmt.Fprintf(w, "col[%d]: %s \n", i, *c)
	}
}
