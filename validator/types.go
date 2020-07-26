package validator

import (
	"go/ast"
)

type LayerName = string

type FilePath = string

type Package = string

type File struct {
	FilePath FilePath
	Package  Package
	Imports  []*ast.ImportSpec
	Layers   []LayerName
}

type ValidationError error
