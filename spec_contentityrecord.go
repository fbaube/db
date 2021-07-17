package dbutils

import (
	"fmt"

	FU "github.com/fbaube/fileutils"
	XU "github.com/fbaube/xmlutils"
)

// ContentityRecord is basically the content plus its "dead properties" -
// properties that are set by the user, rather than dynamically determined.
type ContentityRecord struct {
	error
	Idx_Contentity int
	Idx_Inbatch    int // NOTE: Maybe rename to FILESET. And, could be multiple!
	Descr          string
	FU.PathProps
	Times
	XU.AnalysisRecord
	// For these next two fields, instead put the refs & defs
	//   into another table that FKEY's into this table.
	// ExtlLinkRefs // links that point outside this File
	// ExtlLinkDefs // link targets that are visible outside this File
	// Linker = an outgoing link
	// Linkee = the target of an outgoing link
	// Linkable = a symbol that CAN be a Linkee
}

func (p *ContentityRecord) String() string {
	return fmt.Sprintf("PP<%s> AR <%s>", p.PathProps.String(), p.AnalysisRecord.String())
}

var ColumnSpecs_Contentity = []DbColSpec{
	D_RelFP,
	D_AbsFP,
	D_TmCre,
	D_TmImp,
	D_TmEdt,
	DbColSpec{D_TXT, "descr", "Description", "Content item description"},
	// DbColSpec{D_TXT, "metaraw", "Meta (raw)", "Metadata/header (raw)"},
	// DbColSpec{D_TXT, "textraw", "Text (raw)", "Text/body (raw)"},
	DbColSpec{D_TXT, "mimetype", "MIME type", "MIME type"},
	DbColSpec{D_TXT, "mtype", "MType", "MType"},
	// DbColSpec{D_TXT, "roottag", "Root tag", "XML root tag"},
	// DbColSpec{D_TXT, "rootatts", "Root att's", "XML root tag attributes"},
	DbColSpec{D_TXT, "xmlcontype", "XML contype", "XML content type"},
	DbColSpec{D_TXT, "xmldoctype", "XML Doctype", "XML Doctype"},
	DbColSpec{D_TXT, "ditaflavor", "(Lw)DITA flavor", "(Lw)DITA flavor"},
	DbColSpec{D_TXT, "ditacontype", "(Lw)DITA contype", "(Lw)DITA content type"},
}

var TableConfig_Contentity = TableConfig{
	"contentity",
	// One foreign key
	[]string{"inbatch"},
	ColumnSpecs_Contentity,
}
