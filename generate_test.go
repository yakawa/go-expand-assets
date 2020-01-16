package assets

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func setup() error {
	if err := os.Mkdir("tmp", 0777); err != nil {
		return err
	}
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}
	atime := time.Date(2001, time.January, 1, 0, 0, 0, 0, jst)
	mtime := time.Date(2001, time.January, 1, 0, 0, 0, 0, jst)

	os.Chtimes("tmp", atime, mtime)
	return nil
}

func teardown() error {
	if err := os.RemoveAll("tmp"); err != nil {
		return err
	}
	return nil
}

func TestGenerator_addDirectory(t *testing.T) {
	g := Generator{}
	files, err := ioutil.ReadDir("tmp")
	if err != nil {
		t.Fatal("Read Dir Info Error")
	}

	for _, info := range files {
		if info.Name() == "tmp" {
			err := g.addDirectory("tmp", info)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		os.Exit(1)
	}
	m.Run()
	if err := teardown(); err != nil {
		os.Exit(2)
	}
}
