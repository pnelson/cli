package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type usage struct {
	PkgName string
	VarName string
	Topics  []topic
	out     string
	root    string
}

type topic struct {
	Name string
	Data []byte
}

const indexPage = "README.md"

func main() {
	var (
		out     = flag.String("out", "usage.go", "filename")
		pkgName = flag.String("pkg", "main", "package name")
		varName = flag.String("var", "usage", "variable name")
		dirName = flag.String("dir", "docs", "path to usage docs")
	)
	flag.Parse()
	u := &usage{
		PkgName: *pkgName,
		VarName: *varName,
		Topics:  make([]topic, 0),
		out:     *out,
		root:    *dirName,
	}
	err := write(u)
	if err != nil {
		log.Fatal(err)
	}
}

func write(u *usage) error {
	var buf bytes.Buffer
	t := template.Must(template.New("program").Parse(source))
	root, err := filepath.Abs(u.root)
	if err != nil {
		return err
	}
	u.root = root
	err = filepath.Walk(root, u.walk)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, u)
	if err != nil {
		return err
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(u.out, data, 0644)
}

func (u *usage) walk(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	name := fi.Name()
	if name[0] == '.' {
		if fi.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}
	filename := path
	if fi.IsDir() {
		if name == indexPage {
			return fmt.Errorf("directory named '%s'", indexPage)
		}
		return nil
	}
	return u.add(filename)
}

func (u *usage) add(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename[len(u.root)+1:], indexPage)
	name = strings.TrimSuffix(name, ext)
	topic := topic{Name: name, Data: b}
	u.Topics = append(u.Topics, topic)
	return nil
}

const source = `
// generated by github.com/pnelson/cli/cmd/cli-usage-gen
package {{.PkgName}}

var {{.VarName}} = map[string][]byte{
{{range .Topics}}		"{{.Name}}": []byte({{.Data | printf "%#q"}}),
{{end}}
}
`
