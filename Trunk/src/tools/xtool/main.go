package main

import (
	"bufio"
	"html/template"
	"io"
	"os"
	"strings"

	"fmt"

	"bytes"
	"encoding/xml"
	"flag"
	"path/filepath"
)

type Property struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Desc string `xml:"desc,attr"`
}

type ColType struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Desc string `xml:"desc,attr"`
}

type Record struct {
	Name     string     `xml:"name,attr"`
	Cols     int        `xml:"cols,attr"`
	Desc     string     `xml:"desc,attr"`
	ColTypes []*ColType `xml:"column"`
}

type Include struct {
	Path string `xml:"name,attr"`
}

type Object struct {
	Path      string
	Name      string
	Propertys []*Property `xml:"properties>property"`
	Records   []*Record   `xml:"records>record"`
	Includes  []*Include  `xml:"includes>path"`
	Childs    []*Object
}

var (
	defs    []string
	srcpath = flag.String("p", "", "entity define file's path")
	output  = flag.String("o", "", "parser file's output path")
)

func LoadAllDef(path string) {
	dir, _ := os.Open(path)
	files, _ := dir.Readdir(0)

	for _, f := range files {
		p := path + "/" + f.Name()
		if !f.IsDir() {
			defs = append(defs, p)
			continue
		}
		LoadAllDef(p)
	}

}

func GetType(typ string) string {

	switch strings.ToLower(typ) {
	case "byte", "int8", "int16", "int32", "dword", "word":
		return "int"
	case "widestr":
		return "wstring"
	default:
		return strings.ToLower(typ)
	}
}

func strFirstToUpper(str string) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		upperStr += strings.Title(temp[y])
	}
	return upperStr
}

func NewReaderLabel(label string, input io.Reader) (io.Reader, error) {
	return input, nil
}

func readDef(f string) *Object {

	file, err := os.Open(f)
	if err != nil {
		return nil
	}

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = NewReaderLabel

	obj := &Object{}
	err = decoder.Decode(obj)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if obj != nil && len(obj.Records) > 0 {
		for _, v := range obj.Records {
			for k, c := range v.ColTypes {
				if c.Name == "" {
					c.Name = fmt.Sprintf("Col%d", k)
				}
			}
		}
	}
	return obj
}

func getBaseName(f string) string {
	return strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
}

func main() {
	flag.Parse()
	if *srcpath == "" || *output == "" {
		fmt.Println("usage:data -p path -o output")
		flag.PrintDefaults()
		return
	}

	defs = make([]string, 0, 128)
	LoadAllDef(*srcpath)

	tmp, err := template.New("class.tpl").Funcs(template.FuncMap{"getType": GetType, "setName": strFirstToUpper}).ParseFiles("class.tpl")
	if err != nil {
		panic(err)
	}

	for _, f := range defs {

		obj := readDef(f)
		if obj == nil || (len(obj.Propertys) == 0 && len(obj.Records) == 0 && len(obj.Includes) == 0) {
			fmt.Println("process:", f, "is empty. skip.")
			continue
		}

		var includes []*Object
		if len(obj.Includes) > 0 {
			includes = make([]*Object, 0, len(obj.Includes))
			for _, v := range obj.Includes {
				path := *srcpath + "/../" + v.Path
				o := readDef(path)
				if o == nil || (len(o.Propertys) == 0 && len(o.Records) == 0) {
					continue
				}
				o.Path = path
				includes = append(includes, o)
			}
		}
		obj.Childs = includes
		filebase := getBaseName(f)
		obj.Name = filebase
		obj.Path = f
		outfile, err := os.Create(*output + "/" + filebase + ".h")
		if err != nil {
			panic(err)
		}

		w := bytes.NewBuffer(nil)
		err = tmp.Execute(w, obj)
		if err != nil {
			outfile.Close()
			panic(err)
		}

		writer := bufio.NewWriter(outfile)
		writer.Write(w.Bytes())
		writer.Flush()
		outfile.Close()
	}

}
