package db

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"fmt"
	"log"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
)

type ContentRecord struct {
	Idx int // `db:"idx_content"`
	Idx_Inbatch int
	FU.Paths
	Times
	Raw        string
	Meta_raw   string
	MetaFormat string
	Text_raw   string
	MetaProps  SU.PropSet
	FU.AnalysisRecord
	// For these next two fields, instead put the refs & defs
	//   into another table that FKEY's into this table.
	// ExtlLinkRefs // links that point outside this File
	// ExtlLinkDefs // link targets that are visible outside this File
}

var TableSpec_Content = TableSpec {
      "content",
      []string { "inbatch" }, // FK
      nil, // intFields
      nil, // intRanges
      []string {
				"relfilepath", "absfilepath", // Paths
				"created", "imported", "edited", // Times
				"meta_raw", "text_raw",
				// Analysis
				"mimetype", "mtype", "roottag", "rootatts",
				"xmlcontype", "xmldoctype", "ditamarkuplg", "ditacontype" },
      []string { "Rel.FP (from CLI)",
								 "Absolute filepath",
  							 "Creation date+time",
								 "DB import date+time",
								 "Last edit date+time",
								 "Meta/header (raw)",
								 "Text/body (raw)",
								 "MIME type",
  							 "M-Type",
								 "Root tag",
								 "Root attrs",
								 "XML contype",
								 "XML doctype",
								 "DITA markuplg",
								 "DITA contype" },
  	}

// GetContentAll gets all content in the DB.
func (p *MmmcDB) GetContentAll() (pp []*ContentRecord) {
	pp = make([]*ContentRecord, 0, 16)
	rows, err := p.DB.Queryx("SELECT * FROM CONTENT")
	if err != nil {
		panic("GetContentAll")
	}
	for rows.Next() {
		p := new(ContentRecord)
		err := rows.StructScan(p)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	return pp
}

// InsertContentRecord adds a content item (i.e. a file) to the DB.
func (p *MmmcDB) InsertContentRecord(pC *ContentRecord, pT *sqlx.Tx) (idx int, e error) {
	var err error
	var rslt sql.Result

	// []string { "relfilepath", "absfilepath",
	// 	"creatime", "meta_raw", "text_raw",
	// 	"mimetype", "mtype", "roottag", "rootatts",
	// 	"xmlcontype", "xmldoctype", "ditacontype" },
	println("REL:", pC.RelFilePath)
	println("ABS:", pC.AbsFilePath)
	var s string
	s = fmt.Sprintf(
		"INSERT INTO CONTENT(" +
		"relfilepath, absfilepath, " +
		"created, imported, edited, " +
		"meta_raw, text_raw, " +
		"mimetype, mtype, roottag, rootatts, " +
		"xmlcontype, xmldoctype, ditamarkuplg, ditacontype" +
		") VALUES(" +
			"\"%s\", \"%s\", " +
			"\"%s\", \"%s\", \"%s\", " +
			"\"%s\", \"%s\", " +
			"\"%s\", \"%s\", \"%s\", \"%s\", " +
			"\"%s\", \"%s\", \"%s\", \"%s\")",
		pC.RelFilePath, pC.AbsFilePath,
		pC.Created, pC.Imported, pC.Edited,
		pC.Meta_raw, pC.Text_raw,
		pC.MimeType, pC.MType, pC.RootTag, pC.RootAtts,
		pC.XmlContype, pC.XmlDoctype, pC.DitaMarkupLg, pC.DitaContype)

	println("EXEC:", s)

	rslt, err = pT.NamedExec(s, pC)
	if err != nil {
		println("========")
		println("DB: NamedExec: ERROR:", err.Error())
		println("========")
		panic("INSERT CONTENT failed")
	}
	liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertFile: ID %d nR %d \n", liid, naff)
	return int(liid), nil
}
