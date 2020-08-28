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

// ContentitySections is embedded in db.ContentRecord
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
	FU.PathProps
	// RelFilePath string
	// FU.AbsFilePath
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

func (p *ContentRecord) String() string {
	return fmt.Sprintf("PP<%s> AR <%s>", p.PathProps.String(), p.AnalysisRecord.String())
}

// NewCheckedContent works for directories and symlinks too.
func NewContentRecord(pPP *FU.PathProps) *ContentRecord {
	var e error
	pCR := new(ContentRecord)
	pCR.PathProps = *pPP

	if (!pPP.Exists()) || pPP.IsOkayDir() || pPP.IsOkaySymlink() {
		return pCR
	}
	if !pPP.IsOkayFile() {
		pCR.SetError(errors.New("Is not valid file, directory, or symlink"))
		return pCR
	}
	// OK, it's a file.
	pCR.Raw, e = pPP.FetchContent()
	if e != nil {
		pCR.SetError(fmt.Errorf("db.newCR: Cannot fetch content: %w", e))
		return pCR
	}
	var pAR *FU.AnalysisRecord
	pAR, e = FU.AnalyseFile(pCR.Raw, FP.Ext(string(pPP.AbsFP())))
	if e != nil {
		pCR.SetError(fmt.Errorf("fu.newCR: analyze file failed: %w", e))
		return pCR
	}
	pCR.AnalysisRecord = *pAR
	fmt.Printf("D=> (B:NewCR) %s \n", pCR.String())
	return pCR
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
	println("REL:", pC.RelFP)
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
		pC.RelFP(), pC.AbsFilePath,
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
