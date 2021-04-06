package db

import (
	FU "github.com/fbaube/fileutils"
	SB "github.com/fbaube/semblox"
)

// Inbatch describes a single import batch at the CLI.
type Inbatch struct {
	Idx   int `db:"idx_inbatch"`
	FilCt int
	RelFP string
	AbsFP FU.AbsFilePath // `db:"absfilepath"` // necessary ceremony
	T_Cre string         // RFC 3339
	Descr string
}

// TableSpec_Inbatch describes the table.
var TableSpec_Inbatch = DbTblSpec{SB.D_TBL, "INB", "inbatch", "Batch import of files"}

var ColumnSpecs_Inbatch = []DbColSpec{
	D_RelFP,
	D_AbsFP,
	D_TmCre,
	DbColSpec{SB.D_TXT, "descr", "Batch descr.", "Inbatch description"},
	DbColSpec{SB.D_INT, "filct", "Nr. of files", "Number of files"},
}

var TableConfig_Inbatch = TableConfig{
	"inbatch",
	// no foreign keys
	nil,
	ColumnSpecs_Inbatch,
}
