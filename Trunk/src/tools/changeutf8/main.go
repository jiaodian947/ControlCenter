package main

import (
	"bufio"
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"

	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	defs    []string
	srcpath = flag.String("s", "", "work  path")
)

func ConvertToUTF8(s []byte) ([]byte, error) {
	//utf8 boom
	if len(s) >= 3 && s[0] == 0xef && s[1] == 0xbb && s[2] == 0xbf {
		return s, nil
	}
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func LoadAllDef(srcpath string) {
	dir, _ := os.Open(srcpath)
	files, _ := dir.Readdir(0)

	for _, f := range files {
		p := srcpath + "/" + f.Name()
		if !f.IsDir() {
			ext := strings.ToLower(path.Ext(p))
			if ext == ".xml" || ext == ".ini" || ext == ".txt" {
				defs = append(defs, p)
			}
			continue
		}
		LoadAllDef(p)
	}

}

func process(f string) {
	log.Println("process", f)
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	file.Close()
	utf8data, err := ConvertToUTF8(data)

	if err != nil {
		panic(err)
	}

	outfile, err := os.Create(f)
	if err != nil {
		panic(err)
	}

	bom := []byte{0xef, 0xbb, 0xbf}
	writer := bufio.NewWriter(outfile)
	//write boom
	writer.Write(bom)
	writer.Write(utf8data)
	writer.Flush()
	outfile.Close()
	log.Println("process", f, "done.")
}

func main() {
	flag.Parse()
	if *srcpath == "" {
		*srcpath = "."
	}

	LoadAllDef(*srcpath)

	for _, v := range defs {
		process(v)
	}
}
