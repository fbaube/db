package dbutils

import (
	"errors"
	"fmt"
	"os"

	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
)

// MoveCurrentToBackup makes a best effort but can fail if the
// backup destination is a directory or has a permissions problem.
// The current DB is renamed and so "disappears" from production.
func (p *MmmcDB) MoveCurrentToBackup() error {
	if !p.PathProps.Exists() {
		L.L.Info("No current DB to move to backup")
		return nil
	}
	var cns = NowAsYMDHM()
	var fromFP string = p.PathProps.AbsFP.S()
	var toFP string = FU.AbsFilePath(p.PathProps.AbsFP).BaseName() + "-" + cns + ".db"
	// func os.Rename(oldpath, newpath string) error
	e := os.Rename(fromFP, toFP)
	if e != nil {
		reterr := "can't move current DB to <" + toFP + ">: "
		L.L.Error(reterr + e.Error())
		return errors.New(fmt.Sprintf(reterr+"%w", e))
	}
	L.L.Progress("Old DB moved to: " + toFP)
	return nil
}

// DupeCurrentToBackup makes a best effort but can fail if the
// backup destination is a directory or has a permissions problem.
// The current DB is not affected.
func (p *MmmcDB) DupeCurrentToBackup() error {
	if !p.PathProps.Exists() {
		L.L.Info("No current DB to duplicate to backup")
		return nil
	}
	var cns = NowAsYMDHM()
	var fromFP string = p.PathProps.AbsFP.S()
	var toFP string = FU.AbsFilePath(p.PathProps.AbsFP).BaseName() + "-" + cns + ".db"

	e := FU.CopyFromTo(fromFP, toFP)
	if e != nil {
		reterr := fmt.Sprintf("Can't copy current DB to <%s>: %w: ", toFP, e)
		L.L.Error(reterr)
		return errors.New(reterr)
	}
	L.L.Info("Old DB copied to backup at: " + toFP)
	return nil
}
