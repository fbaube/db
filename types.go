package dbutils

// DbDescr is generic.
type DbDescr Datum

// DbColSpec specifies a datum (i.e. struct field and/or DB column)
// and its generic/portable/DB-independent representation (based on
// the enumeration Datum.TxtIntKeyEtc). Some values for common DB
// columns are defined in the D_* series in the file semblox/types.go
type DbColSpec DbDescr

// DbColIRL describes a column as-is in the DB (as obtained via
// reflection), and has a slot to include the value (as a string).
type DbColIRL DbDescr

// DbTblSpec specifies a DB table (but not its columns!).
type DbTblSpec DbDescr

var D_RelFP = DbColSpec{D_TXT, "relfp", "Rel. path", "Rel.FP (from CLI)"}
var D_AbsFP = DbColSpec{D_TXT, "absfp", "Abs. path", "Absolute filepath"}
var D_TmCre = DbColSpec{D_TXT, "t_cre", "Cre. time", "Creation date+time"}
var D_TmImp = DbColSpec{D_TXT, "t_imp", "Imp. time", "DB import date+time"}
var D_TmEdt = DbColSpec{D_TXT, "t_edt", "Edit time", "Last edit date+time"}
