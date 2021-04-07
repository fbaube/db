package db

import (
	"database/sql"
	"time"

	L "github.com/fbaube/mlog"
	"github.com/jmoiron/sqlx"
)

// GetAll_Inbatch gets all input batches in the system.
func (p *MmmcDB) GetAll_Inbatch() (pp []*Inbatch) {
	var rowsx *sqlx.Rows
	var e error
	rowsx, e = p.DB.Queryx("SELECT * FROM INBATCH")
	if e != nil {
		L.L.Error("DB.GetAll_Inbatch: %w", e)
		return nil
	}
	pp = make([]*Inbatch, 0, 16)
	for rowsx.Next() {
		p := new(Inbatch)
		e = rowsx.StructScan(p)
		if e != nil {
			L.L.Error("DB.GetAll_Inbatch.StructScan: %w", e)
		}
		L.L.Dbg("Got Inbatch: %+v", *p)
		pp = append(pp, p)
	}
	return pp
}

// Add_Inbatch adds an input batch to the DB and returns its primary index.
func (p *MmmcDB) Add_Inbatch(pIB *Inbatch) (int, error) {
	var rslt sql.Result
	var e error
	if pIB.FilCt == 0 {
		pIB.FilCt = 1
	} // HACK

	pIB.T_Cre = time.Now().UTC().Format(time.RFC3339)
	tx := p.MustBegin()
	s := "INSERT INTO INBATCH(" +
		"descr, filct, t_cre, relfp, absfp" +
		") VALUES(" +
		":descr, :filct, :t_cre, :relfp, :absfp)" // " RETURNING i_INB", p)
	rslt, e = tx.NamedExec(s, pIB)
	tx.Commit()
	println("=== ### ===")
	if e != nil {
		L.L.Error("DB.Add_Inbatch: %w", e)
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

	var egInb = Inbatch{}
	var rowsx *sqlx.Rows
	rowsx, e = p.DB.Queryx("SELECT * FROM INBATCH")
	TestInbatch(rowsx, &egInb)

	// WORK HERE

	println("NEED: RETURNING (inbatch ID)")
	liid, err := rslt.LastInsertId()
	if err != nil {
		panic(err)
	}
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertInbatch: ID=%d (nR=%d) \n", liid, naff)
	return int(liid), nil
}
