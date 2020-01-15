package assets

import (
	"os"
	"time"
)

type FileInfo struct {
	Path  string
	Mode  os.FileMode
	CTime time.Time
	MTime time.Time
	Data  string
	Sha1  []byte
}
