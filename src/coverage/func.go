// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements the visitor that computes the (line, column)-(line-column) range for each function.
/*
modification history
----------------------
2015/11/15, by Xiaoye Jiang, modify, add GetCodeCov func for cov.baidu.com
*/

package coverage

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type CodeCov struct {
	LineVaild   int64
	LineCovered int64
	LineRate    float64
	FuncVaild   int64
	FuncCovered int64
	FuncRate    float64
}

func isSkipModules(fn string, skipModules []string) bool {
	for _, module := range skipModules {
		if strings.HasPrefix(fn, module) {
			return true
		}
	}
	return false
}

func GetCodeCov(coverFile string, skipModules []string) (CodeCov, error) {
	var codeCov CodeCov
	profile, err := os.Open(coverFile)
	if err != nil {
		return codeCov, err
	}
	profiles, err := ParseProfiles(profile)
	if err != nil {
		return codeCov, err
	}

	var total, covered int64
	var totalFunc, coveredFunc int64
	for _, pf := range profiles {
		fn := pf.FileName
		if isSkipModules(fn, skipModules) {
			continue
		}
		file, err := findFile(fn)
		if err != nil {
			return codeCov, err
		}
		funcs, err := findFuncs(file)
		if err != nil {
			return codeCov, err
		}

		for _, f := range funcs {
			c, t := f.coverage(pf)
			total += t
			covered += c
			totalFunc += 1
			if c != 0 {
				coveredFunc += 1
			}
		}
	}

	if total == 0 {
		total = 1
	}

	if totalFunc == 0 {
		totalFunc = 1
	}

	codeCov = CodeCov{
		LineVaild:   total,
		LineCovered: covered,
		FuncVaild:   totalFunc,
		FuncCovered: coveredFunc,
		LineRate:    float64(covered) / float64(total),
		FuncRate:    float64(coveredFunc) / float64(totalFunc),
	}
	return codeCov, nil
}

// findFuncs parses the file and returns a slice of FuncExtent descriptors.
func findFuncs(name string) ([]*FuncExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &FuncVisitor{
		fset:    fset,
		name:    name,
		astFile: parsedFile,
	}
	ast.Walk(visitor, visitor.astFile)
	return visitor.funcs, nil
}

// FuncExtent describes a function's extent in the source by file and position.
type FuncExtent struct {
	name      string
	startLine int
	startCol  int
	endLine   int
	endCol    int
}

// FuncVisitor implements the visitor that builds the function position list for a file.
type FuncVisitor struct {
	fset    *token.FileSet
	name    string // Name of file.
	astFile *ast.File
	funcs   []*FuncExtent
}

// Visit implements the ast.Visitor interface.
func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		start := v.fset.Position(n.Pos())
		end := v.fset.Position(n.End())
		fe := &FuncExtent{
			name:      n.Name.Name,
			startLine: start.Line,
			startCol:  start.Column,
			endLine:   end.Line,
			endCol:    end.Column,
		}
		v.funcs = append(v.funcs, fe)
	}
	return v
}

// coverage returns the fraction of the statements in the function that were covered, as a numerator and denominator.
func (f *FuncExtent) coverage(profile *Profile) (num, den int64) {
	// We could avoid making this n^2 overall by doing a single scan and annotating the functions,
	// but the sizes of the data structures is never very large and the scan is almost instantaneous.
	var covered, total int64
	// The blocks are sorted, so we can stop counting as soon as we reach the end of the relevant block.
	for _, b := range profile.Blocks {
		if b.StartLine > f.endLine || (b.StartLine == f.endLine && b.StartCol >= f.endCol) {
			// Past the end of the function.
			break
		}
		if b.EndLine < f.startLine || (b.EndLine == f.startLine && b.EndCol <= f.startCol) {
			// Before the beginning of the function
			continue
		}
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		total = 1 // Avoid zero denominator.
	}
	return covered, total
}
