package db

// TableSpec assumes that for a DB table, it is enough to use strings and
// int(eger)s. A primary key is assumed, and foreign keys are allowed. Any
// field can be nil, except the first. Obviously the two fields "int*""
// should be the same length, and also the two fields "str*".
//
// Date-time's are not an issue for SQLite, since either a string or an int
// can be used. We will favor using strings ("TEXT") as ISO-8601 / RFC 3339.
// See https://www.sqlite.org/datatype3.html#date_and_time_datatype:
//  - TEXT: "YYYY-MM-DD HH:MM:SS.SSS" (or with "T" in the blank position).
//  - REAL as Julian day numbers: the day count since 24 November 4714 BC.
//  - INTEGER as Unix time: the seconds count since 1970-01-01 00:00:00 UTC.
//
type TableSpec struct {
	tableName string
	forenKeys []string
	intFields []string
	intRanges []int    // save to DBrec[0] // -1, 0, 1
	strFields []string
	strDescrs []string // save to DBrec[0]
}
