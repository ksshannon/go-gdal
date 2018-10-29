// Copyright (c) 2017, Kyle Shannon All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var symbols = []struct {
	file    string
	class   string
	goName  string
	match   *regexp.Regexp
	funcs   []string
	funcMap map[string]struct{}
}{
	{
		"gdal.go",
		"GDALDataset",
		"Dataset",
		regexp.MustCompile(`T GDALDataset::[A-Za-z][A-Za-z0-9]*\(`),
		nil,
		make(map[string]struct{}),
	},
	{
		"band.go",
		"GDALRasterBand",
		"Band",
		regexp.MustCompile(`T GDALRasterBand::[A-Za-z][A-Za-z0-9]*\(`),
		nil,
		make(map[string]struct{}),
	},
	{
		"layer.go",
		"OGRLayer",
		"Layer",
		regexp.MustCompile(`T OGRLayer::[A-Za-z][A-Za-z0-9]*\(`),
		nil,
		make(map[string]struct{}),
	},
	{
		"geometry.go",
		"OGRGeometry",
		"Geometry",
		regexp.MustCompile(`T OGRGeometry::[A-Za-z][A-Za-z0-9]*\(`),
		nil,
		make(map[string]struct{}),
	},
	{
		"osr.go",
		"OGRSpatialReference",
		"SpatialReference",
		regexp.MustCompile(`T OGRSpatialReference::[A-Za-z][A-Za-z0-9]*\(`),
		nil,
		make(map[string]struct{}),
	},
}

func variants(f string) []string {
	var v []string
	v = append(v, f)
	if strings.HasPrefix(f, "Get") || strings.HasPrefix(f, "get") {
		v = append(v, f[len("Get"):])
	}
	if f[0] >= 'A' && f[0] <= 'Z' {
		v = append(v, strings.ToLower(f[:1])+f[1:])
	} else {
		v = append(v, strings.ToUpper(f[:1])+f[1:])
	}
	if strings.HasPrefix(f, "export") || strings.HasPrefix(f, "import") {
		s := f[len("export"):]
		v = append(v, s)
		if i := strings.Index(s, "Wkt"); i >= 0 {
			v = append(v, strings.Replace(s, "Wkt", "WKT", -1))
			v = append(v, strings.Replace(f, "Wkt", "WKT", -1))
		} else if i := strings.Index(s, "Url"); i >= 0 {
			v = append(v, strings.Replace(s, "Url", "URL", -1))
			v = append(v, strings.Replace(f, "Url", "WKT", -1))
		} else if i := strings.Index(s, "Json"); i >= 0 {
			v = append(v, strings.Replace(s, "Json", "JSON", -1))
			v = append(v, strings.Replace(f, "Json", "JSON", -1))
		}
	}
	return v
}

func main() {
	libFlag := flag.String("l", "/usr/local/lib/libgdal.so", "path to gdal so")
	flag.Parse()
	var buf bytes.Buffer
	cmd := exec.Command("nm", "-D", "-C", "-g", *libFlag)
	cmd.Stdout = &buf
	cmd.Run()

	for _, sym := range symbols {
		bb := bytes.NewBuffer(buf.Bytes())
		scn := bufio.NewScanner(bb)
		for scn.Scan() {
			if matches := sym.match.FindStringIndex(scn.Text()); len(matches) > 0 {
				s := scn.Text()[matches[0]+len("T "+sym.class+"::") : strings.Index(scn.Text(), "(")]
				// If the function is a constructor or destructor, drop it
				if s == sym.class || s[0] == '~' {
					continue
				}
				sym.funcMap[s] = struct{}{}
				sym.funcs = append(sym.funcs, s)
			}
		}
		fmt.Printf("found %d functions for %s to check\n", len(sym.funcs), sym.class)
		file := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ksshannon", "go-gdal", sym.file)
		fin, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer fin.Close()

		b, _ := ioutil.ReadAll(fin)
		src := string(b)

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, sym.file, src, 0)
		if err != nil {
			panic(err)
		}

		apiMap := make(map[string]string)
		// Inspect the AST and print all identifiers and literals.
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				rcv := "nil"
				if x.Recv != nil {
					for _, v := range x.Recv.List {
						switch xv := v.Type.(type) {
						case *ast.StarExpr:
							if si, ok := xv.X.(*ast.Ident); ok {
								rcv = si.Name
							}
						case *ast.Ident:
							rcv = xv.Name
						}
					}
				}
				apiMap[x.Name.String()] = rcv
			}
			return true
		})
		sort.Strings(sym.funcs)
		for _, k := range sym.funcs {
			fmt.Printf("checking for %s::%s...", sym.class, k)
			rcv := "nil"
			found := false
			var ok bool
			var fx string
			for _, fx = range variants(k) {
				if rcv, ok = apiMap[fx]; ok {
					found = true
					break
				} else {
					rcv = "nil"
				}
			}
			if found {
				fmt.Printf(" ((%s).%s) %c", rcv, fx, '\u2713')
			} else {
				fmt.Printf(" %c", '\u2717')
			}
			fmt.Println()
		}
	}
}
