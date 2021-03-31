package db

import (
	"database/sql"
	"fmt"
	"log"

	FU "github.com/fbaube/fileutils"
)

// Inbatch describes a single import batch at the CLI.
type Inbatch struct {
	Idx         int `db:"idx_inbatch"`
	FileCt      int
	RelFilePath string
	AbsFilePath FU.AbsFilePath `db:"absfilepath"` // necessary ceremony
	Creatime    string         // RFC 3339
	Descr       string
}

/*
type TableSpec struct {
	tableName string
	forenKeys []string
	intFields []string
	intRanges []int // save to DBrec[0] // -1, 0, 1
	strFields []string
	strDescrs []string // save to DBrec[0]
} */

var TableSpec_Inbatch = TableSpec{
	"inbatch",
	// no foreign keys
	nil,
	// One int field
	[]string{"filect"},
	[]int{1}, // >=1
	// Four string fields
	[]string{"relfilepath", "absfilepath", "creatime", "descr"},
	[]string{"Rel.FP (from CLI)", "Absolute filepath",
		"Creation date+time", "Batch description"},
}

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
	if pIB.FileCt == 0 {
		pIB.FileCt = 1
	} // HACK

	tx := p.MustBegin()
	s := "INSERT INTO INBATCH(" +
		"descr, filect, creatime, relfilepath, absfilepath" +
		") VALUES(" +
		":descr, :filect, :creatime, :relfilepath, :absfilepath)" // " RETURNING i_INB", p)
	rslt, err = tx.NamedExec(s, pIB)
	tx.Commit()
	if err != nil {
		panic(err)
	}

	liid, err := rslt.LastInsertId()
	if err != nil {
		panic(err)
	}
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertInbatch: ID=%d (nR=%d) \n", liid, naff)
	return int(liid), nil
}
