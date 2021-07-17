package dbutils

func checkerr(e error) {
	if e == nil {
		return
	}
	panic("Sqlite3 FAILURE: " + e.Error())
}

func mustExecStmt(s string) {
	stmt, e = theDB.DB.Prepare(s)
	checkerr(e)
	_, e := stmt.Exec() // rslt,e := ...
	checkerr(e)
	// liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("DD:mustExecStmt: ID %d nR %d \n", liid, naff)
}

// MustExistTable makes sure it exists but
// does NOT drop an already-existing table.
func MustExistTable(s string) {
	println("db.MustExistTable: WTH: ", s)
	mustExecStmt(s)
}

// ForceEmpty is a convenience function. It first makes a backup.
func (p *MmmcDB) ForceEmpty() {
	if theDB == nil {
		panic("db.forcempty.uninitd.L193")
	}
	p.MoveCurrentToBackup()
	p.ForceExistDBandTables()
}
