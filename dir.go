package assets

import (
	"os"
	"time"
)

type DirInfo struct {
	Path  string
	Mode  os.FileMode
	CTime time.Time
	MTime time.Time
}
