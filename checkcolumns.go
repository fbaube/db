package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (pMDB *MmmcDB) CheckColumns(tableName string) {
	var e error
	var rows *sqlx.Rows
	var cols []interface{}

	rows, e = pMDB.theSqlxDB.Queryx("SELECT * FROM " + tableName + " LIMIT 1")
	if e != nil {
		println("==> CheckColumns-1 failed")
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
	fmt.Printf("    db.chk-cols: c-slice-n: %d \n", n)

	rows, e = pMDB.theSqlxDB.Queryx("SELECT * FROM " + tableName + " LIMIT 1")
	if e != nil {
		println("==> CheckColumns-2 failed")
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
	fmt.Printf("    db.chk-cols: str-map-n: %d \n", n)
}
