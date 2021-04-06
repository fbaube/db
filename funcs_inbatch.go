package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// GetInbatchesAll gets all input batches in the system.
func (p *MmmcDB) GetInbatchesAll() (pp []*Inbatch) {
	pp = make([]*Inbatch, 0, 16)
	rows, err := p.DB.Queryx("SELECT * FROM INBATCH")
	if err != nil {
		panic("GetInbatchesAll")
	}
	for rows.Next() {
		p := new(Inbatch)
		err := rows.StructScan(p)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	return pp
}

// InsertInbatch adds an input batch to the system.
func (p *MmmcDB) InsertInbatch(pIB *Inbatch) (idx int, e error) {
	var err error
	var rslt sql.Result
	if pIB.FilCt == 0 {
		pIB.FilCt = 1
	} // HACK

	tx := p.MustBegin()
	s := "INSERT INTO INBATCH(" +
		"descr, filct, t_cre, relfp, absfp" +
		") VALUES(" +
		":descr, :filct, :t_cre, :relfp, :absfp)" // " RETURNING i_INB", p)
	rslt, err = tx.NamedExec(s, pIB)
	tx.Commit()
	println("=== ### ===")
	if err != nil {
		panic(err)
	}
	/*
			Query(...) (*sql.Rows, error) - unchanged
			QueryRow(...) *sql.Row - unchanged
			Extensions:
			Queryx(...) (*sqlx.Rows, error) - Query, but return an sqlx.Rows
			QueryRowx(...) *sqlx.Row -- QueryRow, but return an sqlx.Row
			New semantics:
			Get(dest interface{}, ...) error // to fetch one scannable
			Select(dest interface{}, ...) error // to fetch multi scannables
			Scannable means: simple datum not struct OR struct w no exported fields OR
		implements sql.Scanner f
			"SELECT * FROM INBATCH"
	*/

	// func StructScan(rows rowsi, dest interface{}) error
	// StructScan all rows from an sql.Rows or an sqlx.Rows into the dest slice.
	// StructScan will scan in the entire rows result; to get fewer, use Queryx
	// and see sqlx.Rows.StructScan. If rows is sqlx.Rows, it will use its mapper,
	// otherwise it will use the default.
	// ============

	db := p.DB
	// var err error
	// func TestInbatch(rows sql.Rows, S *Inbatch)
	var egInb = Inbatch{}

	rows, err := db.Query("SELECT * FROM INBATCH")
	fmt.Printf("rows: %+v \n", rows)
	TestInbatch(rows, &egInb)
	sqlx.StructScan(rows, &egInb)
	fmt.Printf("StructScan got: %+v \n", egInb)
	println("=== ### ===")

	rowsx, err := db.Queryx("SELECT * FROM INBATCH")
	fmt.Printf("rowsx: %+v \n", rowsx)
	TestInbatch(rows, &egInb)
	println("=== ### ===")
	for rowsx.Next() {
		var inb1 Inbatch
		err = rowsx.StructScan(&inb1)
	}
	rrow := db.QueryRow("SELECT * FROM INBATCH")
	fmt.Printf("rrow: %+v \n", rrow)
	println("=== ### ===")
	var inb2 Inbatch
	err = rrow.Scan(&inb2) // could chain this

	rrowx := db.QueryRowx("SELECT * FROM INBATCH")
	fmt.Printf("rrowx: %+v \n", rrowx)
	println("=== ### ===")
	var inb3 Inbatch
	err = rrowx.Scan(&inb3) // could chain this

	inb4 := Inbatch{}
	inb4s := []Inbatch{}

	// this will pull the first place directly into p
	err = db.Get(&inb4, "SELECT * FROM INBATCH LIMIT 1")
	fmt.Printf("inb4: %+v \n", inb4)
	println("=== ### ===")

	// this will pull places with telcode > 50 into the slice pp
	err = db.Select(&inb4s, "SELECT * FROM INBATCH")
	fmt.Printf("inb4s: %+v \n", inb4s)
	println("=== ### ===")

	// ============

	println("NEED: RETURNING (inbatch ID)")
	liid, err := rslt.LastInsertId()
	if err != nil {
		panic(err)
	}
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertInbatch: ID=%d (nR=%d) \n", liid, naff)
	return int(liid), nil
}
