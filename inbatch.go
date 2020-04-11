package db

import (
	"database/sql"
	"fmt"
	"log"

	FU "github.com/fbaube/fileutils"
)

type Inbatch struct {
	Idx         int `db:"idx_inb"`
	Desc        string
	FileCt      int
	Creatime    string // RFC 3339
	RelFilePath string
	AbsFilePath FU.AbsFilePath // necessary ceremony
}

func (pMDB *MmmcDB) GetInbatchAll() (pp []*Inbatch) {
	pp = make([]*Inbatch, 0, 16)
	rows, err := pMDB.theSqlxDB.Queryx("SELECT * FROM INB")
	if err != nil {
		panic("GetInbatchAll")
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

func (pMDB *MmmcDB) InsertInbatch(p *Inbatch) (idx int, e error) {
	// var e error
	var rslt sql.Result
	// rslt, e = pMDB.theSqlxDB.NamedExec(
	rslt, e = pMDB.TheSqlxTxn.NamedExec(
		// var rows *sqlx.Rows
		// rows, e = pMDB.theSqlxDB.NamedQuery(
		"INSERT INTO INB(relfilepath, absfilepath, desc) "+
			"VALUES(:relfilepath, :absfilepath, :desc)", p) // " RETURNING i_INB", p)
	if e != nil {
		return -1, e
	}

	liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertInbatch: ID=%d (nR=%d) \n", liid, naff)
	return int(liid), nil
}
