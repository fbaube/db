package db

import (
	"fmt"

	L "github.com/fbaube/mlog"
	"github.com/jmoiron/sqlx"
)

// TryColumns tries some sqlx stuff.
func (p *MmmcDB) TryColumns(tableName string) {
	var e error
	var rows *sqlx.Rows
	var cols []interface{}

	rows, e = p.DB.Queryx("SELECT * FROM " + tableName + " LIMIT 1")
	if e != nil {
		L.L.Error("TryColumns-1 failed: %v", e)
		return
	}
	n := 0
	for rows.Next() {
		n++
		// cols is an []interface{} of all of the column results
		cols, e = rows.SliceScan()
		if e != nil {
			panic(e)
		} else {
			fmt.Printf("    COLUMNS as SLICE: %+v \n", cols)
		}
	}
	// fmt.Printf("    db.chk-cols: c-slice-n: %d \n", n)

	rows, e = p.DB.Queryx("SELECT * FROM " + tableName + " LIMIT 1")
	if e != nil {
		L.L.Error("CheckColumns-2 failed: %v", e)
		return
	}
	n = 0
	for rows.Next() {
		n++
		results := make(map[string]interface{})
		e = rows.MapScan(results)
		if e != nil {
			panic(e)
		} else {
			fmt.Printf("    COLUMNS as MAP: %+v \n", results)
		}
	}
	// fmt.Printf("    db.chk-cols: str-map-n: %d \n", n)
}
