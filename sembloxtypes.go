package dbutils

import "fmt"

type TxtIntKeyEtc string

const (
	D_TXT TxtIntKeyEtc = "A-z" // SQLite "TEXT"
	D_INT              = "0-9" // SQLite "INT"
	D_KEY              = "Key" // KEY (primary or foreign, SQLite "INTEGER")
	D_TBL              = "Tbl" // This is describing a table !
	D_LST              = "lst" // table type, one per enumetarion ??
	D_TBD              = "??"  // possible future expansion
)

// Datum is a datum descriptor. It describes a single datum or field/column
// or field/column value, and is most useful as an element of an enumeration.
//
type Datum struct {
	// TxtIntKeyEtc is only D_TXT or D_INT (or D_KEY). See this func's
	// comment above re. why this is sufficient, at least for SQLite.
	TxtIntKeyEtc
	// Code is a short unique string token - no spaces or punctuation.
	// We use string codes rather than iota-based integer values for
	// robustness, because values based on iota could change.
	// (1) When a Datum describes a DB column (or table or row's column
	// value or table!), Code is the actual name of the DB field/table,
	// and should be all lower case.
	// (2) In other uses, Code should be all upper case, and all Codes
	// in a particular enumeration "should" be of the same length.
	Code string
	// Name is a short-form name, for common use incl. column headers.
	Name string
	// Descr is a long-form description.
	Descr string
}

func (d Datum) String() string {
	return fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\"",
		d.TxtIntKeyEtc, d.Code, d.Name, d.Descr)
}

// FLDTP describes a simple field that has semantics.
// Field validation is highly desirable but TBD.
type FLDTP Datum

// BLKTP is TBS.
type BLKTP Datum

// InUI is meant to be embedded in any struct
// that uses a Datum to store an entry in a UI.
type InUI struct {
	Visible, Enabled, Selected bool
	Position                   int
}

// FLDTPs are semantic simple fields.
//
var FLDTPs = []FLDTP{
	// INTEGERS (5)
	{D_INT, "INTEG", "Integer", "Generic integer, size unspecified"},
	{D_INT, "BOOL_", "Boolean", "Boolean (0|1)"},
	{D_INT, "PRKEY", "Pri. key", "Primary table key (unique, non-NULL"},
	{D_INT, "FRKEY", "For. key", "Foreign table key"},
	{D_TBD, "FLOAT", "Float", "Generic non-integer, size unspecified"},
	// TEXTS (8)
	{D_TXT, "STRNG", "String", "Generic string, not text"},
	{D_TXT, "TOKEN", "Token", "Generic token (no spaces, punc.)"},
	{D_TXT, "FTEXT", "Free text", "Generic free-flowing text"},
	{D_TXT, "MTEXT", "Markdown", "Markdown (or plain) text"},
	{D_TXT, "JTEXT", "JSON", "JSON content"},
	{D_TXT, "XTEXT", "XML text", "XML text such as (Lw)DITA"},
	{D_TXT, "HTEXT", "HTML5 text", "HTML5 (or previous) text"},
	{D_TXT, "MCFMT", "Microformat", "Microformat record"},
	// TEXT-BASED MISC. (5)
	{D_TXT, "FONUM", "Phone nr.", "Telephone number"},
	{D_TXT, "EMAIL", "Email", "Email address"},
	{D_TXT, "URLIN", "URL/URI/URN", "Generic path ID (URL, URI, URN)"},
	{D_TXT, "DATIM", "Date / Time", "Date and/or time (ISO-8601/RFC-3339)"},
	{D_TXT, "SEMVR", "Sem. ver. nr.", "Semantic version number (x.y.z)"},
}

// BLKTPs are data structures related to semantic field types.
//
var BLKTPs = []BLKTP{
	// LISTS (8)
	{D_LST, "OLIST", "Ord'd list", "Generic ordered list"},
	{D_LST, "ULIST", "Unord. list", "Generic unordered list"},
	{D_LST, "RLIST", "Ranked list", "List ordered by ranking"},
	{D_LST, "SLIST", "Seq. list", "List ordered as a sequence"},
	{D_LST, "ELIST", "Enum. list", "List of enumerated elements"},
	{D_LST, "ENUME", "Enum. item", "Element of enumeration"},
	{D_LST, "XLIST", "Excl. list", "List (select one only)"}, // rbn
	{D_LST, "MLIST", "Mult. list", "List (select multiple)"}, // cbx
	// MISC. STRUX (3)
	{D_TBD, "TABLE", "Table", "Table / dataframe"},
	{D_TBD, "UTREE", "Unord. tree", "Tree of unordered children (e.g. XML data)"},
	{D_TBD, "OTREE", "Ord. tree", "Tree of ordered children (e.g. XML mixed content)"},
}
