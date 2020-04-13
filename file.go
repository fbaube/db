package db

import (
	"database/sql"
	"fmt"
	"log"

	FU "github.com/fbaube/fileutils"
)

// File handles an MMC file.
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

// GetFilesAll gets all files in the DB.
func (p *MmmcDB) GetFilesAll() (pp []*File) {
	pp = make([]*File, 0, 16)
	rows, err := p.theSqlxDB.Queryx("SELECT * FROM FILE")
	if err != nil {
		panic("GetFilesAll")
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

// InsertFile adds a file to the DB.
func (p *MmmcDB) InsertFile(pF *File) (idx int, e error) {
	// var e error
	var rslt sql.Result
	rslt, e = p.TheSqlxTxn.NamedExec(
		// var rows *sqlx.Rows
		// rows, e = p.theSqlxDB.NamedQuery(
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
