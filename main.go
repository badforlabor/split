package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var f = flag.String("f", "", "-f=111.log")
var l = flag.Int("l", 0, "-l=160000")
var s = flag.Int("b", 0, "-b=16m")

func main() {
	flag.Parse()

	if len(*f) == 0 {
		flag.PrintDefaults()
		panic("invalid filename")
		return
	}

	if *l == 0 && *s == 0 {
		flag.PrintDefaults()
		panic("invalid args")
		return
	}

	if *l > 0 {
		SplitFileByLine(*f, *l)
	} else if *s > 0 {
		SplitFileByByteCount(*f, *s)
	}
}

func removeFiles(patten string) {
	var dir = filepath.Dir(patten)
	var files, _ = ioutil.ReadDir(dir)
	for _, v := range files {

		if v.IsDir() {
			continue
		}

		if strings.HasPrefix(v.Name(), patten) {
			os.Remove(v.Name())
		}
	}
}

func SplitFileByLine(file string, lines int) {
	var fileName = file
	var ext = filepath.Ext(file)
	if len(ext) == 0 {
		ext = ".log"
	}

	var splitName = func(num int) string {
		return file + ".part." + strconv.Itoa(num) + ext
	}

	// 删掉旧文件
	removeFiles(file + ".part.")


	fi, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	var count = 0
	br := bufio.NewReader(fi)
	var bw = bytes.NewBuffer([]byte{})
	var num = 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			if len(a) > 0 {
				bw.Write(a)
				bw.WriteString("\r\n")
			}

			ioutil.WriteFile(splitName(num), bw.Bytes(), 666)

			break
		}
		bw.Write(a)
		bw.WriteString("\r\n")
		count++

		if count >= lines {

			ioutil.WriteFile(splitName(num), bw.Bytes(), 666)
			bw.Reset()
			num++
			count = 0
		}
	}
}

func SplitFileByByteCount(file string, rawMaxSize int) {
	var fileName = file
	var ext = filepath.Ext(file)
	if len(ext) == 0 {
		ext = ".bin"
	}

	var splitName = func(num int) string {
		return file + ".part." + strconv.Itoa(num) + ext
	}

	// 删掉旧文件
	removeFiles(file + ".part.")


	fi, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	var count = 0
	br := bufio.NewReader(fi)
	var bw = bytes.NewBuffer([]byte{})
	var num = 0
	var buf = make([]byte, 4096)
	for {
		size, c := br.Read(buf)
		if c == io.EOF {
			if size > 0 {
				bw.Write(buf[0:size])
			}

			ioutil.WriteFile(splitName(num), bw.Bytes(), 666)

			break
		}
		bw.Write(buf[0:size])
		count += size

		if count >= rawMaxSize {

			ioutil.WriteFile(splitName(num), bw.Bytes(), 666)
			bw.Reset()
			num++
			count = 0
		}
	}
}
