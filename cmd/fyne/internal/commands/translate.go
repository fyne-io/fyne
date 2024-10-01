package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/urfave/cli/v2"
)

// Translate returns the cli command to scan for new translation strings.
func Translate() *cli.Command {
	return &cli.Command{
		Name:  "translate",
		Usage: "Scans for new translation strings.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Show files that are being scanned etc.",
			},
			&cli.BoolFlag{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update existing translations (use with care).",
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"n"},
				Usage:   "Scan without storing the results.",
			},
			&cli.StringFlag{
				Name:    "sourceDir",
				Aliases: []string{"d"},
				Usage:   "Directory to scan recursively for go files.",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:    "translationsFile",
				Aliases: []string{"f"},
				Usage:   "File to read from and write translations to.",
				Value:   "translations/en.json",
			},
		},
		Action: func(ctx *cli.Context) error {
			sourceDir := ctx.String("sourceDir")
			translationsFile := ctx.String("translationsFile")
			files := ctx.Args().Slice()

			// without any argument find all .go files in the source directory
			if len(files) == 0 {
				sources, err := findFilesExt(sourceDir, ".go")
				if err != nil {
					return err
				}
				files = sources
			}

			// update the translation file by scanning the given source files
			return updateTranslationsFile(ctx, translationsFile, files)
		},
	}
}

// Recursively walk the given directory and return all files with the matching extension
func findFilesExt(dir, ext string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ext || fi.IsDir() || !fi.Mode().IsRegular() {
			return nil
		}

		files = append(files, path)

		return nil
	})
	return files, err
}

// Create or add to translations file by scanning the given files for translation calls
func updateTranslationsFile(ctx *cli.Context, file string, files []string) error {
	translations := make(map[string]interface{})

	// try to get current translations first
	f, err := os.Open(file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	defer f.Close()

	// only try to parse translations that exist
	if f != nil {
		dec := json.NewDecoder(f)
		if err := dec.Decode(&translations); err != nil {
			return err
		}
	}

	if ctx.Bool("verbose") {
		fmt.Fprintf(os.Stderr, "scanning files: %v\n", files)
	}

	// update translations hash
	if err := updateTranslationsHash(ctx, translations, files); err != nil {
		return err
	}

	// serialize in readable format for humas to change
	b, err := json.MarshalIndent(translations, "", "\t")
	if err != nil {
		return err
	}

	// avoid writing an empty file
	if len(translations) == 0 {
		fmt.Fprintln(os.Stderr, "No translations found")
		return nil
	}

	// stop without making any changes
	if ctx.Bool("dry-run") {
		return nil
	}

	// support writing to stdout
	if file == "-" {
		fmt.Printf("%s\n", string(b))
		return nil
	}

	return writeTranslationsFile(b, file, f)
}

// Write data to given file, optionally using same permissions as the original file
func writeTranslationsFile(b []byte, file string, f *os.File) error {
	// default permissions
	perm := fs.FileMode(0644)

	// use same permissions as original file when possible
	if f != nil {
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		perm = fi.Mode().Perm()
	}

	// use temporary file to do atomic change
	nf, err := os.CreateTemp(filepath.Dir(file), filepath.Base(file)+"-*")
	if err != nil {
		return err
	}

	n, err := nf.Write(b)
	if err != nil {
		return err
	}

	if n < len(b) {
		return io.ErrShortWrite
	}

	if err := nf.Chmod(perm); err != nil {
		return err
	}

	if err := nf.Close(); err != nil {
		return err
	}

	// atomic switch to new file
	return os.Rename(nf.Name(), file)
}

// Update translations hash by scanning the given files
func updateTranslationsHash(ctx *cli.Context, m map[string]interface{}, srcs []string) error {
	for _, src := range srcs {
		// get AST by parsing source
		fset := token.NewFileSet()
		af, err := parser.ParseFile(fset, src, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		// walk AST to find known translation calls
		ast.Walk(&visitor{ctx: ctx, m: m}, af)
	}

	return nil
}

// Visitor pattern with state machine to find translations fallback (and key)
type visitor struct {
	ctx      *cli.Context
	state    stateFn
	name     string
	key      string
	fallback string
	m        map[string]interface{}
}

// Method to walk AST using interface for ast.Walk
func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	// start over any time there is no state
	if v.state == nil {
		v.state = translateNew
		v.name = ""
		v.key = ""
		v.fallback = ""
	}

	// run and get next state
	v.state = v.state(v, node)

	return v
}

// State machine to pick out translation key and fallback from AST
type stateFn func(*visitor, ast.Node) stateFn

// All translation calls need to start with the literal "lang"
func translateNew(v *visitor, node ast.Node) stateFn {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return nil
	}

	if ident.Name != "lang" {
		return nil
	}

	return translateCall
}

// A known translation method needs to be used
func translateCall(v *visitor, node ast.Node) stateFn {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return nil
	}

	v.name = ident.Name

	switch ident.Name {
	// simple cases: only the first argument is relevant
	case "L", "Localize":
		return translateLocalize
	case "N", "LocalizePlural":
		return translateLocalize
	// more complex cases: first two arguments matter
	case "X", "LocalizeKey":
		return translateKey
	case "XN", "LocalizePluralKey":
		return translateKey
	}

	return nil
}

// Parse first argument, use string as key and fallback, and finish
func translateLocalize(v *visitor, node ast.Node) stateFn {
	basiclit, ok := node.(*ast.BasicLit)
	if !ok {
		return nil
	}

	val, err := strconv.Unquote(basiclit.Value)
	if err != nil {
		return nil
	}

	v.key = val
	v.fallback = val

	return translateFinish(v, node)
}

// Parse first argument and use as key
func translateKey(v *visitor, node ast.Node) stateFn {
	basiclit, ok := node.(*ast.BasicLit)
	if !ok {
		return nil
	}

	val, err := strconv.Unquote(basiclit.Value)
	if err != nil {
		return nil
	}

	v.key = val

	return translateKeyFallback
}

// Parse second argument and use as fallback, and finish
func translateKeyFallback(v *visitor, node ast.Node) stateFn {
	basiclit, ok := node.(*ast.BasicLit)
	if !ok {
		return nil
	}

	val, err := strconv.Unquote(basiclit.Value)
	if err != nil {
		return nil
	}

	v.fallback = val

	return translateFinish(v, node)
}

// Finish scan for translation and add to translation hash with the right type (singular or plural)
func translateFinish(v *visitor, node ast.Node) stateFn {
	// only adding new keys, ignoring changed or removed (ha!) ones
	// removing is dangerous as there could be dynamic keys that get removed

	// Ignore existing translations to prevent unintentional overwriting
	_, found := v.m[v.key]
	if found {
		if !v.ctx.Bool("update") {
			if v.ctx.Bool("verbose") {
				fmt.Fprintf(os.Stderr, "ignoring: %s\n", v.key)
			}
			return nil
		}
		if v.ctx.Bool("verbose") {
			fmt.Fprintf(os.Stderr, "updating: %s\n", v.key)
		}
	} else {
		if v.ctx.Bool("verbose") {
			fmt.Fprintf(os.Stderr, "adding: %s\n", v.key)
		}
	}

	switch v.name {
	// Plural translations use a nested map
	case "LocalizePlural", "LocalizePluralKey", "N", "XN":
		m := make(map[string]string)
		m["other"] = v.fallback
		v.m[v.key] = m
	default:
		v.m[v.key] = v.fallback
	}

	return nil
}
