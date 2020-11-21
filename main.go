package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	dir := filepath.Dir(os.Args[0])
	fi, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range fi {
		if !f.IsDir() && strings.HasSuffix(f.Name(), "dat") {
			imgname, err := Dat2Image(path.Join(dir, f.Name()))
			if err != nil {
				fmt.Println("ERROR, dat file name:", f.Name(), err)
			} else {
				fmt.Println("SUCCESS, image:", imgname)
			}
		}
	}
}

const (
	jpg0 = 0xFF
	jpg1 = 0xD8
	gif0 = 0x47
	gif1 = 0x49
	png0 = 0x89
	png1 = 0x50
)

func Dat2Image(datpath string) (string, error) {
	b, err := ioutil.ReadFile(datpath)
	if err != nil {
		return "", err
	}
	if len(b) < 2 {
		return "", errors.New("image size error")
	}

	j0 := b[0] ^ jpg0
	j1 := b[1] ^ jpg1
	g0 := b[0] ^ gif0
	g1 := b[1] ^ gif1
	p0 := b[0] ^ png0
	p1 := b[1] ^ png1
	var v byte
	var ext string
	if j0 == j1 {
		v = j0
		ext = "jpg"
	} else if g0 == g1 {
		v = g0
		ext = "gif"
	} else if p0 == p1 {
		v = p0
		ext = "png"
	} else {
		return "", errors.New("unknown image format")
	}

	for i := range b {
		b[i] = b[i] ^ v
	}

	imgpath := datpath[0:len(datpath)-len(ext)] + ext
	err = ioutil.WriteFile(imgpath, b, os.ModePerm)
	return imgpath, err
}
