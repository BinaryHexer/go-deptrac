package validator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var Log = log.New(ioutil.Discard, "[go-deptrac] ", log.LstdFlags|log.Lshortfile)

type Validator struct {
	config *Config
	dirs   map[string][]File // dir -> files mapping
}

// NewValidator creates new Validator.
func NewValidator(config Config) *Validator {
	dirs := make(map[string][]File, 0)

	return &Validator{
		config: &config,
		dirs:   dirs,
	}
}

func (v *Validator) Validate(ignoreTests bool) (bool, []ValidationError, error) {
	var errors []ValidationError

	root := v.config.Paths[0]
	Log.Printf("starting in root: %s", root)
	// Walk through all the files and assign them layers
	err := filepath.Walk(root, v.walkFn(ignoreTests))
	if err != nil {
		return false, nil, err
	}

	// Walk through all the packages and validate their dependencies (imports)
	for _, files := range v.dirs {
		for _, file := range files {
			errs := v.validateImports(file)
			errors = append(errors, errs...)
		}
	}

	return len(errors) == 0, errors, nil
}

func (v *Validator) walkFn(ignoreTests bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if skip(path, info, ignoreTests) {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			panic(err)
		}

		Log.Printf("processing (%s) in package (%s)", path, f.Name.String())

		dir, err := filepath.Rel(v.config.baseDir, path)
		dir = filepath.Dir(dir)
		layers := v.layers(path)
		file := File{
			FilePath: path,
			Package:  f.Name.String(),
			Imports:  f.Imports,
			Layers:   layers,
		}

		if files, ok := v.dirs[dir]; ok {

			v.dirs[dir] = append(files, file)
		}
		v.dirs[dir] = []File{file}

		Log.Printf("file: %+v", file)

		return nil
	}
}

func skip(path string, fi os.FileInfo, ignoreTests bool) bool {
	if fi.IsDir() {
		return true
	}

	if !strings.HasSuffix(path, ".go") {
		return true
	}

	if ignoreTests && strings.HasSuffix(path, "_test.go") {
		return true
	}

	if strings.Contains(path, "/vendor/") {
		// todo - better check and flag
		return true
	}

	if strings.Contains(path, "/.") {
		return true
	}

	return false
}

func (v *Validator) layers(path string) []LayerName {
	var layers []LayerName

	for _, layer := range v.config.Layers {
		for _, collector := range layer.Collectors {
			if collector.regex.MatchString(path) {
				layers = append(layers, layer.Name)
			}
		}
	}

	return layers
}

func (v *Validator) validateImports(file File) []ValidationError {
	var errors []ValidationError

	for _, import_ := range file.Imports {
		errs := v.validateImport(file.FilePath, file.Layers, import_)
		errors = append(errors, errs...)
	}

	return errors
}

func (v *Validator) validateImport(path string, importerLayers []LayerName, import_ *ast.ImportSpec) []ValidationError {
	var errors []ValidationError

	root := v.config.Paths[0]
	root, err := filepath.Rel(v.config.baseDir, root)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(root + "/.*")

	importPath := import_.Path.Value
	importPath = strings.TrimSuffix(importPath, `"`)
	importPath = strings.TrimPrefix(importPath, `"`)
	importPath = re.FindString(importPath)
	files, ok := v.dirs[importPath]
	if !ok {
		Log.Printf("root: %s", root)
		Log.Printf("no files found for package: %s", importPath)
		Log.Printf("no files found for package: %s", import_.Path.Value)
	}

	for _, importerLayer := range importerLayers {
		for _, importedFile := range files {
			for _, importedLayer := range importedFile.Layers {
				if !v.importAllowed(importerLayer, importedLayer) {
					err := fmt.Errorf(
						"you cannot import %s Layer (%s) to %s Layer (%s)",
						importedLayer, importPath,
						importerLayer, path,
					)
					errors = append(errors, err)
				}
			}
		}
	}

	return errors
}

func (v *Validator) importAllowed(importerLayer LayerName, importedLayer LayerName) bool {
	allowedLayers := v.config.Ruleset[importerLayer]

	for _, allowedLayer := range allowedLayers {
		if importedLayer == allowedLayer {
			return true
		}
	}

	return false
}
