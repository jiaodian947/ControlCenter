package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"tools/upgrade/upgrader"

	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	srcpath = flag.String("s", "", "work  path")
	rename  = flag.Bool("r", false, "rename path")
)

func ConvertToUTF8(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	var O *transform.Reader
	if s[0] == 0xff && s[1] == 0xfe {
		O = transform.NewReader(I, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder())
	} else if s[0] == 0xfe && s[1] == 0xff { //unicode
		O = transform.NewReader(I, unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder())
	} else {
		O = transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	}

	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

var defs []string

func LoadAllDef(srcpath string) {
	dir, _ := os.Open(srcpath)
	files, _ := dir.Readdir(0)

	for _, f := range files {
		p := srcpath + "/" + f.Name()
		if !f.IsDir() {
			ext := strings.ToLower(path.Ext(p))
			if ext == ".h" || ext == ".cpp" || ext == ".hpp" || ext == ".c" || ext == ".filters" || ext == ".vcxproj" {
				defs = append(defs, p)
			}
			continue
		}
		LoadAllDef(p)
	}
}

var paths []string

func LoadAllPath(srcpath string) {
	dir, _ := os.Open(srcpath)
	files, _ := dir.Readdir(0)

	for _, f := range files {
		p := srcpath + "/" + f.Name()
		if !f.IsDir() || strings.HasPrefix(f.Name(), ".") {
			continue
		}
		LoadAllPath(p)
		paths = append(paths, p)
	}
}

func convtoutf8(f string) {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	file.Close()
	if len(data) == 0 {
		return
	}
	//utf8 boom
	if len(data) >= 3 && data[0] == 0xef && data[1] == 0xbb && data[2] == 0xbf {
		return
	}

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
}

var regfilterinc = regexp.MustCompile(`\<Filter\sInclude="(?P<file>[\w\./\\]+)"`)
var regcic = regexp.MustCompile(`\<ClCompile\sInclude="(?P<file>[\w\./\\]+)"`)
var regfilter = regexp.MustCompile(`\<Filter\>(?P<file>[\w\./\\]+)\</Filter\>`)
var regcii = regexp.MustCompile(`\<ClInclude\sInclude="(?P<file>[\w\./\\]+)"`)

func replace(reg *regexp.Regexp, file string, rep string) string {
	return reg.ReplaceAllStringFunc(file, func(s string) string {
		params := upgrader.GetParams(reg, s)
		if p, has := params["file"]; has {
			p = upgrader.ReplaceUpper(p)
			//fmt.Println(fmt.Sprintf(rep, upgrader.SnakeString(p)))
			return fmt.Sprintf(rep, upgrader.SnakeString(p))
		}
		return s
	})
}

func changeFilter(f string) {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()

	fd := string(data)

	fd = replace(regfilterinc, fd, `<Filter Include="%s"`)
	fd = replace(regcic, fd, `<ClCompile Include="%s"`)
	fd = replace(regfilter, fd, `<Filter>%s</Filter>`)
	fd = replace(regcii, fd, `<ClInclude Include="%s"`)

	//fmt.Println(fd)
	ioutil.WriteFile(f, []byte(fd), 0666)

}

func main() {
	flag.Parse()
	if *srcpath == "" {
		*srcpath = "./"
	}

	if *rename {
		LoadAllPath(*srcpath)
		//fmt.Println(paths)
		for _, v := range paths {
			base := path.Base(v)
			base = upgrader.ReplaceUpper(base)
			newname := upgrader.SnakeString(base)
			fmt.Println(v, "=>", path.Dir(v)+"/"+newname)
			if err := os.Rename(v, path.Dir(v)+"/"+newname); err != nil {
				fmt.Println("err:", err)
			}
			//time.Sleep(time.Second)
		}
		return
	}
	LoadAllDef(*srcpath)

	for _, v := range defs {
		if strings.ToLower(path.Ext(v)) == ".filters" || strings.ToLower(path.Ext(v)) == ".vcxproj" {
			changeFilter(v)
			continue
		}

		convtoutf8(v)
		log.Println("Convert encoding:", v, " ==> UTF8.")
		log.Print("Upgrade ")
		upgrader.UpgradeFile(v)
	}
}
