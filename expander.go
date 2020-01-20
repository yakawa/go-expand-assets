package assets

import (
	"os"
	"time"
)

type Expander struct {
}

type ExpandFileInfo struct {
	Path  string
	Data  string
	Mode  os.FileMode
	MTime time.Time
	Sha1  string
}

type ExpandDirInfo struct {
	Path  string
	Mode  os.FileMode
	MTime time.Time
}

func NewExpander(fi []ExpandFileInfo, di []ExpandDirInfo) *Expander {
	e := &Expander{}

	return e
}
