package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"strings"

	"path"

	"github.com/tealeg/xlsx"
)

const (
	MAX_SPACE_ROWS = 3
)

var (
	defs    []string
	srcpath = flag.String("i", "", "excel include path")
	outputs = flag.String("s", "", "server output path")
	outputc = flag.String("c", "", "client output path")
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

type Value struct {
	Index int
	Value string
}

type Object struct {
	Values []*Value
}

type Property struct {
	Type     string
	Name     string
	Sign     string
	Index    int
	IsServer bool
	IsClient bool
}

type ConfigRoot struct {
	Name      string
	Propertys []*Property
	Objects   []*Object
}

type GlobalValue struct {
	Name     string
	Value    string
	IsServer bool
	IsClient bool
}

type Global struct {
	KV []*GlobalValue
}

// 服务器属性数量
func (c *ConfigRoot) ServerPropertys() int {
	count := 0
	for _, v := range c.Propertys {
		if v != nil && v.IsServer {
			count++
		}
	}
	return count
}

// 客户端属性数量
func (c *ConfigRoot) ClientPropertys() int {
	count := 0
	for _, v := range c.Propertys {
		if v != nil && v.IsClient {
			count++
		}
	}
	return count
}

// 处理属性属于server还是client
func parseSign(p *Property, s string) {
	switch strings.ToLower(s) {
	case "server":
		p.IsServer = true
	case "client":
		p.IsClient = true
	}
}

// 处理全局参数属于server还是client
func parseGlobalSign(p *GlobalValue, s string) {
	switch strings.ToLower(s) {
	case "server":
		p.IsServer = true
	case "client":
		p.IsClient = true
	}
}

// 处理属性定义
func parseProperty(p *Property) bool {
	if p.Name == "" {
		return false
	}

	switch p.Type {
	case "int", "double", "bool", "string":
	default:
		return false
	}

	cs := strings.Split(p.Sign, "/")
	parseSign(p, cs[0])
	if len(cs) == 2 {
		parseSign(p, cs[1])
	}

	if p.IsClient || p.IsServer {
		return true
	}

	return false
}

// 解析一行
func parseRow(c *ConfigRoot, cells []*xlsx.Cell) {
	s, err := cells[0].String()
	if err != nil {
		return
	}

	if strings.HasPrefix(s, "//") || strings.HasPrefix(s, "#") {
		return
	}

	obj := &Object{}
	obj.Values = make([]*Value, 0, len(cells))
	for k, v := range cells {
		if k >= len(c.Propertys) || c.Propertys[k] == nil {
			continue
		}
		str, err := v.String()
		if err != nil || str == "" {
			continue
		}

		val := &Value{}
		val.Index = k
		val.Value = str

		obj.Values = append(obj.Values, val)
	}

	if len(obj.Values) > 0 {
		c.Objects = append(c.Objects, obj)
	}
}

// 解析全局参数
func parseGlobal(sheet *xlsx.Sheet) *Global {
	g := &Global{}
	g.KV = make([]*GlobalValue, 0, sheet.MaxRow)

	spaceRows := 0
	keyIndex := 1
	valueIndex := 2
	signIndex := 4
	var err error
	for i := 1; i < sheet.MaxRow; i++ {
		c := sheet.Rows[i].Cells
		if len(c) == 0 {
			spaceRows++
			if spaceRows > MAX_SPACE_ROWS {
				break
			}
			continue
		}

		if firstCell, err := c[0].String(); firstCell == "" || err != nil || len(c) < 5 {
			continue
		}

		spaceRows = 0

		gv := &GlobalValue{}
		gv.Name, err = c[keyIndex].String()
		if err != nil || gv.Name == "" {
			continue
		}
		gv.Name = strings.TrimSpace(gv.Name)

		gv.Value, err = c[valueIndex].String()
		if err != nil || gv.Value == "" {
			continue
		}
		gv.Value = strings.TrimSpace(gv.Value)

		var s string
		s, err = c[signIndex].String()
		if err != nil || s == "" {
			continue
		}
		s = strings.TrimSpace(s)

		ss := strings.Split(s, "/")
		parseGlobalSign(gv, ss[0])
		if len(ss) == 2 {
			parseGlobalSign(gv, ss[1])
		}

		if !gv.IsClient && !gv.IsServer {
			continue
		}

		g.KV = append(g.KV, gv)
	}

	return g
}

// 解析一页
func parseSheet(sheet *xlsx.Sheet) *ConfigRoot {
	ns := strings.Split(strings.ToLower(sheet.Name), "|")
	if len(ns) != 2 {
		return nil
	}
	if sheet.MaxRow <= 4 { // excel至少有四行，第1行注释，第二行类型，第三行名称，第四行server or client
		return nil
	}

	typeIndex := 1
	nameIndex := 2
	signIndex := 3

	config := &ConfigRoot{}
	maxCell := len(sheet.Rows[1].Cells)
	config.Propertys = make([]*Property, maxCell)
	config.Objects = make([]*Object, 0, sheet.MaxRow-4)
	config.Name = ns[1]

	// 处理定义
	var err error
	cols := 0
	for i := 0; i < maxCell; i++ {
		p := &Property{}
		p.Index = i
		p.Type, err = sheet.Rows[typeIndex].Cells[i].String()
		if err != nil {
			fmt.Println("get type error", i, err)
			continue
		}

		p.Type = strings.TrimSpace(p.Type)

		p.Name, err = sheet.Rows[nameIndex].Cells[i].String()
		if err != nil {
			fmt.Println("get name error", i, err)
			continue
		}
		p.Name = strings.TrimSpace(p.Name)

		p.Sign, err = sheet.Rows[signIndex].Cells[i].String()
		if err != nil {
			fmt.Println("get sign error", i, err)
			continue
		}
		p.Sign = strings.TrimSpace(p.Sign)

		if parseProperty(p) {
			config.Propertys[i] = p
			cols++
		}
	}

	if cols == 0 {
		return nil
	}

	// 处理数据部分
	spaceRows := 0
	for i := 4; i < sheet.MaxRow; i++ {
		c := sheet.Rows[i].Cells
		if len(c) == 0 {
			spaceRows++
			if spaceRows > MAX_SPACE_ROWS {
				break
			}
			continue
		}

		if firstCell, err := c[0].String(); firstCell == "" || err != nil {
			continue
		}
		spaceRows = 0
		parseRow(config, c)
	}

	return config
}

// GBK -> UTF8
func Encode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 导出服务器xml
func outputXmlServer(data interface{}, file string, name string, tpl string) {

	tmp, err := template.New(tpl).Funcs(template.FuncMap{}).ParseFiles(tpl)
	if err != nil {
		panic(err)
	}

	outfile, err := os.Create(file)
	if err != nil {
		panic(name + err.Error())
	}
	defer outfile.Close()

	err = tmp.Execute(outfile, data)
	if err != nil {
		panic(name + err.Error())
	}

	fmt.Println("write to:", file, "done")
}

// 导出客户端xml
func outputXmlClient(data interface{}, file string, name string, tpl string) {

	tmp, err := template.New(tpl).Funcs(template.FuncMap{}).ParseFiles(tpl)
	if err != nil {
		panic(name + err.Error())
	}

	outfile, err := os.Create(file)
	if err != nil {
		panic(name + err.Error())
	}
	defer outfile.Close()

	err = tmp.Execute(outfile, data)
	if err != nil {
		panic(name + err.Error())
	}

	fmt.Println("write to:", file, "done")
}

// 加载所有原始文件信息
func LoadAllDef(path string) {
	dir, _ := os.Open(path)
	files, _ := dir.Readdir(0)

	for _, f := range files {
		p := path + "/" + f.Name()
		if !f.IsDir() {
			defs = append(defs, p)
			continue
		}
		//LoadAllDef(p)
	}

}

// 获取文件的文件名(不包含扩展名)
func getBaseName(f string) string {
	return strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
}

// 导出excel
func exportXlsx(file string, spath, cpath string) {
	//basename := getBaseName(file)
	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		fmt.Println(file, err)
		os.Exit(1)
	}

	for _, sheet := range xlFile.Sheets {
		fmt.Println("process:", file, sheet.Name)
		// excel的页面字名格式为name1|name2 其中name2为导出的文件名,name1为策划用的
		ns := strings.Split(sheet.Name, "|")
		if len(ns) != 2 {
			continue
		}

		if strings.ToLower(ns[1]) == "global" { //全局参数
			g := parseGlobal(sheet)
			if g != nil {
				outputXmlServer(g, fmt.Sprintf("%s/%s.xml", spath, "global"), "Global", "global_server.tpl")
				outputXmlClient(g, fmt.Sprintf("%s/%s.xml", cpath, "global"), "Global", "global_client.tpl")
			}
			continue
		}

		c := parseSheet(sheet)
		if c != nil {
			if c.ServerPropertys() > 0 {
				outputXmlServer(c, fmt.Sprintf("%s/%s.xml", spath, c.Name), c.Name, "server.tpl")
			}
			if c.ClientPropertys() > 0 {
				outputXmlClient(c, fmt.Sprintf("%s/%s.xml", cpath, c.Name), c.Name, "client.tpl")
			}
		}
	}
}

func main() {
	flag.Parse()

	if *srcpath == "" || *outputs == "" || *outputc == "" {
		fmt.Println("usage:config_tool.exe -i path -s server path -c client path")
		flag.PrintDefaults()
		return
	}

	defs = make([]string, 0, 128)
	LoadAllDef(*srcpath)

	wg := &WaitGroupWrapper{}
	for _, v := range defs {
		if path.Ext(v) == ".xlsx" {
			f := v
			wg.Wrap(func() {
				exportXlsx(f, *outputs, *outputc)
			})
		}
	}

	wg.Wait()
}
