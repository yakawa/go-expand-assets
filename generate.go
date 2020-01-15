package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type Generator struct {
	PackageName string
	DirList     map[string]DirInfo
	FileList    map[string]FileInfo
}

func (g *Generator) addDirectory(path string, info os.FileInfo) error {
	wd, _ := os.Getwd()
	apath, _ := filepath.Abs(path)
	rpath, _ := filepath.Rel(wd, apath)

	if _, ok := g.DirList[rpath]; ok {
		return errors.New("Already Same path exists")
	}

	var s syscall.Stat_t
	syscall.Stat(path, &s)

	d := DirInfo{
		Path:  rpath,
		CTime: time.Unix(s.Ctimespec.Sec, 0),
		MTime: info.ModTime(),
		Mode:  info.Mode(),
	}

	g.DirList[rpath] = d
	return nil
}

func (g *Generator) addFile(path string, info os.FileInfo) error {
	wd, _ := os.Getwd()
	apath, _ := filepath.Abs(path)
	rpath, _ := filepath.Rel(wd, apath)

	if _, ok := g.FileList[rpath]; ok {
		return errors.New("Already Same path exists")
	}

	d := filepath.Dir(rpath)
	if _, ok := g.DirList[d]; !ok {
		i, _ := os.Stat(d)
		g.addDirectory(d, i)
	}
	var s syscall.Stat_t
	syscall.Stat(path, &s)

	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}

	h := sha1.New()
	h.Write(data)
	bs := h.Sum(nil)

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	zw.Write(data)
	zw.Close()

	b64buf := base64.StdEncoding.EncodeToString(buf.Bytes())

	f := FileInfo{
		Path:  rpath,
		CTime: time.Unix(s.Ctimespec.Sec, 0),
		MTime: info.ModTime(),
		Mode:  info.Mode(),
		Sha1:  bs,
		Data:  b64buf,
	}

	g.FileList[rpath] = f

	return nil
}

func (g *Generator) Add(d string) error {
	if g.DirList == nil {
		g.DirList = make(map[string]DirInfo)
	}
	if g.FileList == nil {
		g.FileList = make(map[string]FileInfo)
	}

	err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == true {
			g.addDirectory(path, info)
		} else {
			g.addFile(path, info)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) Write(w io.Writer) error {
	return nil
}

func (g *Generator) MakeFileList(s string) {

}
