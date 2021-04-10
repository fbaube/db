package db

import (
	FU "github.com/fbaube/fileutils"
	SB "github.com/fbaube/semblox"
)

// Inbatch describes a single import batch at the CLI.
type Inbatch struct {
	Idx_Inbatch int
	FilCt       int
	RelFP       string
	AbsFP       FU.AbsFilePath
	T_Cre       string
	Descr       string
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
