package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	FP "path/filepath"
	"time"

	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	XM "github.com/fbaube/xmlmodels"
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
	pAR, e = FU.AnalyseFile(pCR.Raw, FP.Ext(string(pPP.AbsFP)))
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
func (p *MmmcDB) InsertContentityRecord(pC *ContentityRecord) (int, error) {
	var rslt sql.Result
	var stmt string
	var e error
	println("REL:", pC.RelFP)
	println("ABS:", pC.AbsFP)

	/*
		if pC.TextRaw() == "" {
			pC.textraw = "thetextraw"
		}
		if pC.MetaRaw() == "" {
			pC.metaraw = "themetaraw"
		}
	*/

	pC.T_Cre = time.Now().UTC().Format(time.RFC3339)
	tx := p.MustBegin()
	stmt = "INSERT INTO CONTENTITY(" +
		"idx_inbatch, descr, relfp, absfp, " +
		"t_cre, t_imp, t_edt, " +
		// "metaraw, textraw, " +
		"mimetype, mtype, " +
		// roottag, rootatts, " +
		"xmlcontype, xmldoctype, ditaflavor, ditacontype" +
		") VALUES(" +

		// ":idx_inbatch, :pathprops.relfp, :pathprops.absfp, " +
		":idx_inbatch, :descr, :relfp, :absfp, " +

		// ":times.t_cre, :times.t_imp, :times.t_edt, " +
		":t_cre, :t_imp, :t_edt, " +

		// ":metaraw, :textraw, " +
		// ":mimetype, :mtype, " +
		":mimetype, :mtype, " +
		// ":root.name, :root.atts, " +
		// ":analysisrecord.contentitystructure.root.name, " +
		// ":analysisrecord.contentitystructure.root.atts, " +

		":xmlcontype, :doctype, :ditaflavor, :ditacontype);"
		// ":doctype, :ditaflavor, :ditacontype);"

	rslt, e = tx.NamedExec(stmt, pC)
	tx.Commit()
	println("=== ### ===")
	if e != nil {
		L.L.Error("DB.Add_Contentity: %s", e.Error())
	}
	if e != nil {
		println("========")
		println("DB: NamedExec: ERROR:", e.Error())
		println("========")
		fnam := "./insert-Contentity-failed.sql"
		e = ioutil.WriteFile(fnam, []byte(stmt), 0644)
		if e != nil {
			L.L.Error("Could not write file: " + fnam)
		} else {
			L.L.Dbg("Wrote \"INSERT INTO contentity ... \" to: " + fnam)
		}
		panic("INSERT CONTENTITY failed")
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
