package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	FP "path/filepath"

	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	XM "github.com/fbaube/xmlmodels"
	"github.com/jmoiron/sqlx"
)

// NewContentityRecord works for directories and symlinks too.
// It used to SetError(..), but no longer does.
func NewContentityRecord(pPP *FU.PathProps) *ContentityRecord {
	var e error
	pCR := new(ContentityRecord)
	pCR.PathProps = *pPP

	if !pPP.Exists() {
		pCR.SetError(errors.New("Does not exist"))
		return pCR
	}
	if pPP.IsOkayDir() || pPP.IsOkaySymlink() {
		// COMMENTING THIS OUT IS A FIX
		// pCR.SetError(errors.New("Is directory or symlink"))
		return pCR
	}
	if !pPP.IsOkayFile() {
		pCR.SetError(errors.New("Is not valid file"))
		return pCR
	}
	// OK, it's a valid file.
	pCR.Raw, e = pPP.FetchContent()
	if e != nil {
		L.L.Error("DB.newCnty: cannot fetch content: " + e.Error())
		pCR.SetError(fmt.Errorf("DB.newCnty: cannot fetch content: %w", e))
		return pCR
	}
	var pAR *XM.AnalysisRecord
	pAR, e = FU.AnalyseFile(pCR.Raw, FP.Ext(string(pPP.AbsFP())))
	if e != nil {
		L.L.Error("DB.newCnty: analyze file failed: " + e.Error())
		pCR.SetError(fmt.Errorf("fu.newCR: analyze file failed: %w", e))
		return pCR
	}
	if pAR == nil {
		panic("NIL pAR")
	}
	pCR.AnalysisRecord = *pAR
	// SPLIT FILE!
	if !pAR.ContentityStructure.HasNone() {
		L.L.Okay("Key element triplet: Root<%s> Meta<%s> Text<%s>",
			pAR.ContentityStructure.Root.String(),
			pAR.ContentityStructure.Meta.String(),
			pAR.ContentityStructure.Text.String())
	} else if pAR.FileType() == "MKDN" {
		// pAR.KeyElms.SetToAllText()
		L.L.Warning("TODO set MKDN all text, and ranges")
	} else {
		L.L.Warning("Found no key elms (root,meta,text)")
	}
	// fmt.Printf("D=> NewCR: %s \n", pCR.String())
	return pCR
}

// GetContentityAll gets all content in the DB.
func (p *MmmcDB) GetContentityAll() (pp []*ContentityRecord) {
	pp = make([]*ContentityRecord, 0, 16)
	rows, err := p.DB.Queryx("SELECT * FROM CONTENT")
	if err != nil {
		panic("GetContentityAll")
	}
	for rows.Next() {
		p := new(ContentityRecord)
		err := rows.StructScan(p)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	return pp
}

// InsertContentityRecord adds a content item (i.e. a file) to the DB.
func (p *MmmcDB) InsertContentityRecord(pC *ContentityRecord, pT *sqlx.Tx) (idx int, e error) {
	var err error
	var rslt sql.Result
	println("REL:", pC.RelFP())
	println("ABS:", pC.AbsFP())
	var s string
	s = fmt.Sprintf(
		"INSERT INTO CONTENTITY("+
			"relfp, absfp, "+
			"t_cre, t_imp, t_edt, "+
			"metaraw, textraw, "+
			"mimetype, mtype, roottag, rootatts, "+
			"xmlcontype, xmldoctype, ditaflavor, ditacontype"+
			") VALUES("+
			"\"%s\", \"%s\", "+
			"\"%s\", \"%s\", \"%s\", "+
			"\"%s\", \"%s\", "+
			"\"%s\", \"%s\", \"%s\", \"%s\", "+
			"\"%s\", \"%s\", \"%s\", \"%s\")",
		pC.RelFP(), pC.AbsFP(),
		pC.Created, pC.Imported, pC.Edited,
		pC.GetSpan(pC.Meta), pC.GetSpan(pC.Text),
		pC.MimeType, pC.MType, pC.Root.Name, pC.Root.Atts,
		pC.XmlContype, pC.Doctype, pC.DitaFlavor, pC.DitaContype)

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

func (p *ContentityRecord) HasError() bool {
	return p.error != nil && p.error.Error() != ""
}

// GetError is necessary cos "Error()"" dusnt tell you whether "error"
// is "nil", which is the indication of no error. Therefore we need
// this function, which can actually return the telltale "nil".
func (p *ContentityRecord) GetError() error {
	return p.error
}

// Error satisfies interface "error", but the
// weird thing is that "error" can be nil.
func (p *ContentityRecord) Error() string {
	if p.error != nil {
		return p.error.Error()
	}
	return ""
}

func (p *ContentityRecord) SetError(e error) {
	p.error = e
}
