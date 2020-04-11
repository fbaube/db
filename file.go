package db

import (
	"database/sql"
	"fmt"
	"log"

	FU "github.com/fbaube/fileutils"
)

type File struct {
	Idx_File    int    //  `db:"idx_file"`
	Idx_Inb     int    // `db:"idx_inb"`
	Creatime    string // RFC 3339
	Contype     string
	RootTag     string
	MimeType    string
	Doctype     string
	Mtype       string
	RelFilePath string
	AbsFilePath FU.AbsFilePath // necessary ceremony
	Hed         string
	Bod         string
}

func (pMDB *MmmcDB) GetFileAll() (pp []*File) {
	pp = make([]*File, 0, 16)
	rows, err := pMDB.theSqlxDB.Queryx("SELECT * FROM FILE")
	if err != nil {
		panic("GetFileAll")
	}
	for rows.Next() {
		p := new(File)
		err := rows.StructScan(p)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	return pp
}

func (pMDB *MmmcDB) InsertFile(p *File) (idx int, e error) {
	// var e error
	var rslt sql.Result
	rslt, e = pMDB.TheSqlxTxn.NamedExec(
		// var rows *sqlx.Rows
		// rows, e = pMDB.theSqlxDB.NamedQuery(
		/*
						Idx         int    `db:"i_FILE"`
						InbatchIdx  int    `db:"i_INB"`
						Creatime    string // RFC 3339
						Contype     string
						RootTag     string
						MimeType    string
						Doctype     string
						Mtype       string
						RelFilePath FU.RelFilePath // necessary ceremony
						AbsFilePath FU.AbsFilePath // necessary ceremony


			  trackartist INTEGER,
			  FOREIGN KEY(trackartist) REFERENCES artist(artistid)
				per DB cnxn: PRAGMA foreign_keys = ON;
		*/
		"INSERT INTO FILE(contype, rootelm, mimetype, "+
			"doctype, mtype, idx_inb, relfilepath, absfilepath, hed, bod) "+
			// "doctype, mtype, relfilepath, absfilepath) "+
			"VALUES(:contype, :rootelm, :mimetype, :doctype, :mtype, "+
			":idx_inb, :relfilepath, :absfilepath, :hed, :bod)", p) // " RETURNING i_INB", p)

	// ":doctype, :mtype, :relfilepath, :absfilepath)", p) // " RETURNING i_INB", p)
	if e != nil {
		return -1, e
	}

	liid, _ := rslt.LastInsertId()
	naff, _ := rslt.RowsAffected()
	fmt.Printf("    DD:InsertFile: ID %d nR %d \n", liid, naff)
	return int(liid), nil
}
