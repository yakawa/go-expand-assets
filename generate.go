package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Path  string
	Mode  os.FileMode
	MTime time.Time
	Data  string
	Sha1  []byte
}

type DirInfo struct {
	Path  string
	Mode  os.FileMode
	MTime time.Time
}

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

	d := DirInfo{
		Path:  rpath,
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

	b64buf := g.encodeB64(buf)

	f := FileInfo{
		Path:  rpath,
		MTime: info.ModTime(),
		Mode:  info.Mode(),
		Sha1:  bs,
		Data:  b64buf,
	}

	g.FileList[rpath] = f

	return nil
}

func (g *Generator) encodeB64(b bytes.Buffer) string {
	b64str := base64.StdEncoding.EncodeToString(b.Bytes())
	var buf bytes.Buffer
	for k, c := range strings.Split(b64str, "") {
		buf.WriteString(c)
		if k%76 == 75 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
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
	writer := &bytes.Buffer{}

	body := `package main

import (
	"time"

	"github.com/yakawa1128/go-expand-assets"
)

var (
%s
)

var AssetsExpander = assets.NewExpander(
%s,
%s
)
`
	fileDefine := ""
	fileArg := "[]assets.ExpandFileInfo{"
	dirArg := "[]assets.ExpandDirInfo{"
	for _, v := range g.FileList {
		id := fmt.Sprintf("%x", v.Sha1)
		t := fmt.Sprintf("\t_Fi_%s_ = `%s`\n", id, v.Data)
		fileDefine += t
		fileArg += fmt.Sprintf("\n\tassets.ExpandFileInfo{\n")
		fileArg += fmt.Sprintf("\t\tPath: \"%s\",\n", v.Path)
		fileArg += fmt.Sprintf("\t\tData: %s,\n", "_Fi_"+id+"_")
		fileArg += fmt.Sprintf("\t\tMode: %#o,\n", v.Mode.Perm())
		fileArg += fmt.Sprintf("\t\tSha1: \"%x\",\n", v.Sha1)
		fileArg += fmt.Sprintf("\t\tMTime: time.Unix(%d, 0),\n", v.MTime.Unix())
		fileArg += fmt.Sprintf("\t},")
	}
	fileArg += "\n}"
	for _, v := range g.DirList {
		dirArg += fmt.Sprintf("\n\tassets.ExpandDirInfo{\n")
		dirArg += fmt.Sprintf("\t\tPath: \"%s\",\n", v.Path)
		dirArg += fmt.Sprintf("\t\tMode: %#o,\n", v.Mode.Perm())
		dirArg += fmt.Sprintf("\t\tMTime: time.Unix(%d, 0),\n", v.MTime.Unix())
		dirArg += fmt.Sprintf("\t},")
	}
	dirArg += "\n}"
	fmt.Printf(body, fileDefine, fileArg, dirArg)
	fmt.Fprintf(writer, body)
	return nil
}
