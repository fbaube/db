package db

import (
	"fmt"
	"os"
	FU "github.com/fbaube/fileutils"
)

// MoveCurrentToBackup makes a best effort but can fail if the
// backup destination is a directory or has a permissions problem.
// The current DB is renamed and so "disappears" from production.
func (p *MmmcDB) MoveCurrentToBackup() error {
	if !p.PathInfo.Exists {
		println("    --> No current DB to move to backup")
		return nil
	}
	var cns = NowAsYMDHM()
	var fromFP string = p.PathInfo.AbsFP()
	var   toFP string = p.PathInfo.AbsFilePathParts.BaseName + "-" + cns + ".db"
	// func os.Rename(oldpath, newpath string) error
	e := os.Rename(fromFP, toFP)
	if e != nil {
		return fmt.Errorf("Can't move current DB to <%s>: %w: ", toFP, e)
	}
	println("    --> Old DB moved to:", toFP)
	return nil
}

// DupeCurrentToBackup makes a best effort but can fail if the
// backup destination is a directory or has a permissions problem.
// The current DB is not affected.
func (p *MmmcDB) DupeCurrentToBackup() error {
	if !p.PathInfo.Exists {
		println("    --> No current DB to duplicate to backup")
		return nil
	}
	var cns = NowAsYMDHM()
	var fromFP string = p.PathInfo.AbsFP()
	var   toFP string = p.PathInfo.BaseName + "-" + cns + ".db"

	e := FU.CopyFromTo(fromFP, toFP)
	if e != nil {
		return fmt.Errorf("Can't copy current DB to <%s>: %w: ", toFP, e)
	}
	println("    --> Old DB copied to backup at:", toFP)
	return nil
}
