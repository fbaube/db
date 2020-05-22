package db

import (
	"database/sql"
	"fmt"
	"log"

	FU "github.com/fbaube/fileutils"
)

type Content struct {
	Idx int         // `db:"idx_content"`
	Idx_Inbatch int // `db:"idx_inbatch"`
	BaseInfo
	Meta_raw string
	Text_raw string
	Analysis
	// For these next two fields, instead put the refs & defs
	//   into another table that FKEY's into this table.
	// ExtlLinkRefs // links that point outside this File
	// ExtlLinkDefs // link targets that are visible outside this File
}

type BaseInfo struct {
	RelFilePath string
	// AbsFilePath FU.AbsFilePath // necessary ceremony (problem in sqlx)
	AbsFilePath FU.AbsFilePath `db:"absfilepath"` // necessary ceremony
	Creatime string // ISO-8601 / RFC 3339
}

type Analysis struct {
	MimeType    string
	Mtype       string
	RootTag     string
	RootAtts    string // e.g. <html lang="en">
	XmlContype  string
	XmlDoctype  string
	DitaContype string
}

var TableSpec_Content = TableSpec {
      "content",
      []string { "inbatch" }, // FK
      nil, // intFields
      nil, // intRanges
      []string { "relfilepath", "absfilepath",
				"creatime", "meta_raw", "text_raw",
				"mimetype", "mtype", "roottag", "rootatts",
				"xmlcontype", "xmldoctype", "ditacontype" },
      []string { "Rel.FP (from CLI)",
								 "Absolute filepath",
  							 "Creation date+time",
								 "Meta/header (raw)",
								 "Text/body (raw)",
								 "MIME type",
  							 "M-Type",
								 "Root tag",
								 "Root attrs",
								 "XML contype",
								 "XML doctype",
								 "DITA contype" },
  	}

// GetContentAll gets all content in the DB.
func (p *MmmcDB) GetContentAll() (pp []*Content) {
	pp = make([]*Content, 0, 16)
	rows, err := p.DB.Queryx("SELECT * FROM CONTENT")
	if err != nil {
		panic("GetContentAll")
	}
	for rows.Next() {
		p := new(Content)
		err := rows.StructScan(p)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	return pp
}

// InsertContent adds a content item (i.e. a file) to the DB.
func (p *MmmcDB) InsertContent(pC *Content) (idx int, e error) {
	var err error
	var rslt sql.Result

	tx := p.MustBegin()
	// []string { "relfilepath", "absfilepath",
	// 	"creatime", "meta_raw", "text_raw",
	// 	"mimetype", "mtype", "roottag", "rootatts",
	// 	"xmlcontype", "xmldoctype", "ditacontype" },
	s := "INSERT INTO CONTENT(" +
		"relfilepath, absfilepath, creatime, meta_raw, text_raw, " +
		"mimetype, mtype, roottag, rootatts, xmlcontype, xmldoctype, ditacontype" +
		") VALUES(" +
		":relfilepath, :absfilepath, :creatime, :meta_raw, :text_raw, " +
		":mimetype, :mtype, :roottag, :rootatts, :xmlcontype, :xmldoctype, :ditacontype)" // " RETURNING i_INB", p)
	rslt, err = tx.NamedExec(s, pC)
	tx.Commit()
	if err != nil {
		panic(err)
	}

	liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertFile: ID %d nR %d \n", liid, naff)
	return int(liid), nil
}
