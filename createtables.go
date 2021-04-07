package db

import (
	SB "github.com/fbaube/semblox"
)

var AllTableConfigs = []TableConfig{
	TableConfig_Inbatch,
	TableConfig_Contentity,
	// TableConfig_Topicref,
}

// Topicref describes a reference from a Map to a Topic. Note that
// this does NOT (necessarily) refer to the DITA `topictref` element!
//
// The relationship is N-to-N btwn Maps and Topics, so `schemaTREF` might
// not be unique because a topic might be explicitly referenced more than
// once by a map. So for simplicity, let's create only one `schemaTREF` per
// topic per map file, and see if it creates problems elsewhere later on.
//
type Topicref struct {
	Idx         int `db:"idx_topicref"`
	Idx_MapCnty int
	Idx_TpcCnty int
}

// TableSpec_Topicref describes the table.
var TableSpec_Topicref = DbTblSpec{SB.D_TBL,
	"TRF", "topicref", "Reference from map to topic"}

var ColumnSpecs_Topicref = []DbColSpec{
	// NONE!
}

var TableConfig_Topicref = TableConfig{
	"topicref",
	// ONLY foreign keys
	[]string{"contentity", "contentity"},
	ColumnSpecs_Topicref,
}
