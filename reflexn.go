package db

import (
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

// rows, _ := db.Query("select * from inbatch;")
func TestInbatch(rowsx *sqlx.Rows, S *Inbatch) {

	// Set up data structures for DB row
	columns, _ := rowsx.Columns()
	colTypes, _ := rowsx.ColumnTypes()
	colCount := len(columns)
	colValues := make([]interface{}, colCount)
	colValPtrs := make([]interface{}, colCount)
	// Map entries aren't sposta have stable addresses,
	// so this is kinda dodgy. Use ASAP and then discard.
	for i := 0; i < colCount; i++ {
		colValPtrs[i] = &colValues[i]
	}
	// Introspect the DB Row's columns
	for i := 0; i < colCount; i++ {
		fmt.Printf("%s(%s) ", columns[i], colTypes[i].DatabaseTypeName())
	}
	println(" ")
	/*
		for i := 0; i < colCount; i++ {
			c := columns[i]
			ct := colTypes[i]
			fmt.Printf("=== COLUMN %d :: %s :: %+v\n", i, c, ct)
			/*
				func (ci *ColumnType) DatabaseTypeName() string
				func (ci *ColumnType) DecimalSize() (precision, scale int64, ok bool)
				func (ci *ColumnType) Length() (length int64, ok bool)
				func (ci *ColumnType) Name() string
				func (ci *ColumnType) Nullable() (nullable, ok bool)
				func (ci *ColumnType) ScanType() reflect.Type
			* /
			decPrec, decScale, _ := ct.DecimalSize()
			leng, _ := ct.Length()
			nulbl, _ := ct.Nullable()
			fmt.Printf("DBTN<%s> DecSz<%d:%d> Len<%d> Name<%s> Nulbl<%t> ScanTp<%v> \n",
				ct.DatabaseTypeName(), decPrec, decScale, leng, ct.Name(), nulbl, ct.ScanType())
		}
	*/

	// Now work on the struct

	p := &Inbatch{}

	//long and bored code
	t := reflect.TypeOf(*p)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			fmt.Println(t.Field(i).Name)
		}
	} else {
		fmt.Println("not a stuct")
	}

	//shorthanded call
	fmt.Println(reflect.TypeOf(*p).Field(0).Name) //can panic if no field exists

	/*

			final_result := map[int]map[string]string{}
			result_id := 0
			for rows.Next() {

				rows.Scan(vaPtrs...)

				tmp_struct := map[string]string{}

				for i, col := range columns {
					var v interface{}
					val := values[i]
					b, ok := val.([]byte)
					if ok {
						v = string(b)
					} else {
						v = val
					}
					tmp_struct[col] = fmt.Sprintf("%s", v)
				}

				final_result[result_id] = tmp_struct
				result_id++
			}

			fmt.Println(final_result)

			// (2)

			for rows.Next() {

				// Scan the result into the column pointers...
				if err := rows.Scan(columnPointers...); err != nil {
					return err
				}

				// Create our map, and retrieve the value for each column from the ptrs
				// slice, storing it in the map with the name of the column as the key.
				m := make(map[string]interface{})
				for i, colName := range cols {
					val := columnPointers[i].(*interface{})
					m[colName] = *val
				}

				// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
				fmt.Print(m)
			}

			// (3) SQLX

			rows, err := db.Queryx("SELECT * FROM place")
			for rows.Next() {
				results := make(map[string]interface{})
				err = rows.MapScan(results)
			}

			// ...
		}

		// Like sql.Rows.Scan, but scans a single Row into a map[string]interface{}.
		// Use this to get results for SQL that might not be under your control
		// (for instance, if you're building an interface for an SQL server that
		// executes SQL from input).  Please do not use this as a primary interface!
		// This will modify the map sent to it in place, so do not reuse the same one
		// on different queries or you may end up with something odd!
		//
		// The resultant map values will be string representations of the various
		// SQL datatypes for existing values and a nil for null values.
		func MapScan(r ColScanner, dest map[string]interface{}) error {
			// ignore r.started, since we needn't use reflect for anything.
			columns, err := r.Columns()
			if err != nil {
				return err
			}

			values := make([]interface{}, len(columns))
			for i, _ := range values {
				values[i] = &sql.NullString{}
			}

			r.Scan(values...)

			for i, column := range columns {
				ns := *(values[i].(*sql.NullString))
				if ns.Valid {
					dest[column] = ns.String
				} else {
					dest[column] = nil
				}
			}

			return nil

	*/

}
