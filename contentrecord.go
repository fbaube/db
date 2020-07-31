package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	FP "path/filepath"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	"github.com/jmoiron/sqlx"
)

type ContentitySections struct {
	Raw string // The entire input file
	// Text_raw + Meta_raw = Raw (maybe plus surrounding tags)
	Text_raw   string
	Meta_raw   string
	MetaFormat string
	MetaProps  SU.PropSet
}

// ContentRecord is basically the content plus its "dead properties" -
// properties that are set by the user, rather than dynamically determined.
type ContentRecord struct {
	error
	Idx         int // `db:"idx_content"`
	Idx_Inbatch int // NOTE: Maybe rename to FILESET. And, could be multiple!
	RelFilePath string
	FU.AbsFilePath
	Times
	ContentitySections
	FU.AnalysisRecord
	// For these next two fields, instead put the refs & defs
	//   into another table that FKEY's into this table.
	// ExtlLinkRefs // links that point outside this File
	// ExtlLinkDefs // link targets that are visible outside this File
	// Linker = an outgoing link
	// Linkee = the target of an outgoing link
	// Linkable = a symbol that CAN be a Linkee
}

// NewCheckedContent works for directories and symlinks too.
func NewContentRecord(pPI *FU.PathProps) *ContentRecord {
	var e error
	pCC := new(ContentRecord)

	// pCC.PathInfo = *pPI
	pCC.AbsFilePath = FU.AbsFilePath(pPI.AbsFP())
	pCC.RelFilePath = pPI.RelFP()

	if pPI.IsOkayDir() || pPI.IsOkaySymlink() {
		return pCC
	}
	if !pPI.IsOkayFile() {
		pCC.SetError(errors.New("Is not valid file, directory, or symlink"))
		return pCC
	}
	// OK, it's a file.
	pCC.Raw, e = pPI.FetchContent()
	if e != nil {
		pCC.SetError(errors.New("Could not fetch content"))
		return pCC
	}
	// pCC.BasicAnalysis.FileIsOkay = true
	pBA, e := FU.AnalyseFile(pCC.Raw, FP.Ext(string(pPI.AbsFP())))
	if e != nil {
		// panic(e)
		pCC.SetError(fmt.Errorf("fu.CC: analyze file failed: %w", e))
		return pCC
	}
	pCC.AnalysisRecord = *pBA
	// println("NewCC OK!")
	return pCC
}

func NewContentRecordFromPath(path string) *ContentRecord {
	bp := FU.NewPathProps(path)
	return NewContentRecord(bp)
}

var TableSpec_Content = TableSpec{
	"content",
	[]string{"inbatch"}, // FK
	nil,                 // intFields
	nil,                 // intRanges
	[]string{
		"relfilepath", "absfilepath", // Paths
		"created", "imported", "edited", // Times
		"meta_raw", "text_raw",
		// Analysis
		"mimetype", "mtype", "roottag", "rootatts",
		"xmlcontype", "xmldoctype", "ditamarkuplg", "ditacontype"},
	[]string{"Rel.FP (from CLI)",
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
		"DITA contype"},
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
	println("REL:", pC.RelFilePath)
	println("ABS:", pC.AbsFilePath)
	var s string
	s = fmt.Sprintf(
		"INSERT INTO CONTENT("+
			"relfilepath, absfilepath, "+
			"created, imported, edited, "+
			"meta_raw, text_raw, "+
			"mimetype, mtype, roottag, rootatts, "+
			"xmlcontype, xmldoctype, ditamarkuplg, ditacontype"+
			") VALUES("+
			"\"%s\", \"%s\", "+
			"\"%s\", \"%s\", \"%s\", "+
			"\"%s\", \"%s\", "+
			"\"%s\", \"%s\", \"%s\", \"%s\", "+
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

// === Implement interface Errable

func (p *ContentRecord) HasError() bool {
	return p.error != nil && p.error.Error() != ""
}

// GetError is necessary cos "Error()"" dusnt tell you whether "error"
// is "nil", which is the indication of no error. Therefore we need
// this function, which can actually return the telltale "nil".
func (p *ContentRecord) GetError() error {
	return p.error
}

// Error satisfies interface "error", but the
// weird thing is that "error" can be nil.
func (p *ContentRecord) Error() string {
	if p.error != nil {
		return p.error.Error()
	}
	return ""
}

func (p *ContentRecord) SetError(e error) {
	p.error = e
}
