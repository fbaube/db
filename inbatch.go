package db

import (
	"database/sql"
	// "github.com/jmoiron/sqlx"
	"fmt"
	"log"
	S "strings"
	FU "github.com/fbaube/fileutils"
)

// Inbatch describes a single import batch at the CLI.
type Inbatch struct {
	Idx         int `db:"idx_inb"`
	Descr       string
	FileCt      int
	Creatime    string // RFC 3339
	RelFilePath string
	AbsFilePath FU.AbsFilePath `db:"absfilepath"` // necessary ceremony
}

func (pDB *MmmcDB) CreateTable_Inbatch_sqlite() {
	pDB.CreateTable_sqlite("inbatch",
    nil,     // no foreign keys
    []string { "filect" },
    []int    {  1  },  // >=1
    []string { "relfilepath", "absfilepath", "creatime", "descr" },
    []string { "Rel.FP (from CLI)", "Absolute filepath",
               "Creation date+time", "Import description" },
							 )
		}

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
		panic("FKs-1")
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
		panic("FKs-2")
	}
	CTS = S.TrimSuffix(CTS, ", \n") + "\n)"
	println("= = = = = = = = \n", CTS, "= = = = = = = =")
	pDB.theSqlxDB.MustExec(CTS)
	// TODO: Insert record with IDX 0 and string descriptions
	}

/*
func (b Inbatch) Enmap() (map[string]int, map[string]string) {
	iMap := make(map[string]int)
	sMap := make(map[string]string)
	iMap["inbatch_idx"] = b.Idx
	iMap["filect"] = b.FileCt
	sMap["descr"]  = b.Descr
	sMap["creatime"] = b.Creatime
	sMap["relfilepath"] = b.RelFilePath
	sMap["absfilepath"] = b.AbsFilePath.S()
	return iMap, sMap
}
*/

// GetInbatchesAll gets all input batches in the system.
func (p *MmmcDB) GetInbatchesAll() (pp []*Inbatch) {
	pp = make([]*Inbatch, 0, 16)
	rows, err := p.theSqlxDB.Queryx("SELECT * FROM INBATCH")
	if err != nil {
		panic("GetInbatchesAll")
	}
	for rows.Next() {
		p := new(Inbatch)
		err := rows.StructScan(p)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	return pp
}

// InsertInbatch adds an input batch to the system.
func (p *MmmcDB) InsertInbatch(pIB *Inbatch) (idx int, e error) {
	// var e error
	var rslt sql.Result
	// rslt, e = p.theSqlxDB.NamedExec(
	rslt, e = p.TheSqlxTxn.NamedExec(
		// var rows *sqlx.Rows
		// rows, e = p.theSqlxDB.NamedQuery(
		"INSERT INTO INBATCH(relfilepath, absfilepath, descr) "+
			"VALUES(:relfilepath, :absfilepath, :descr)", p) // " RETURNING i_INB", p)
	if e != nil {
		return -1, e
	}

	liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertInbatch: ID=%d (nR=%d) \n", liid, naff)
	return int(liid), nil
}
