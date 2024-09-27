package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
				Usage:   "Shows files that are being scanned etc.",
			},
			&cli.StringFlag{
				Name:    "sourceDir",
				Aliases: []string{"src"},
				Usage:   "Directory to scan recursively for go files.",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:    "translationsFile",
				Aliases: []string{"file"},
				Usage:   "File to read from and write translations to.",
				Value:   "translations/en.json",
			},
		},
		Action: func(ctx *cli.Context) error {
			sourceDir := ctx.String("sourceDir")
			translationsFile := ctx.String("translationsFile")
			files := ctx.Args().Slice()

			if len(files) == 0 {
				err := filepath.Walk(sourceDir, func(path string, fi fs.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if fi.IsDir() {
						return nil
					}

					if !fi.Mode().IsRegular() {
						return nil
					}

					if filepath.Ext(path) != ".go" {
						return nil
					}

					files = append(files, path)

					return nil
				})
				if err != nil {
					return err
				}
			}

			if ctx.Bool("verbose") {
				fmt.Printf("files: %v\n", files)
			}

			translations := make(map[string]interface{})

			// get current translations
			f, err := os.Open(translationsFile)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			defer f.Close()

			if f != nil {
				dec := json.NewDecoder(f)
				if err := dec.Decode(&translations); err != nil {
					return err
				}
			}

			// update translations hash
			if err := updateTranslations(translations, files); err != nil {
				return err
			}

			// serialize in readable format for humas to change
			b, err := json.MarshalIndent(translations, "", "\t")
			if err != nil {
				return err
			}

			if ctx.Bool("verbose") {
				fmt.Printf("%s\n", string(b))
			}

			if len(translations) == 0 {
				fmt.Println("No translations found")
				return nil
			}

			// use temporary file to do atomic change
			nf, err := os.CreateTemp(filepath.Dir(translationsFile), filepath.Base(translationsFile)+"-*")
			if err != nil {
				return err
			}

			n, err := nf.Write(b)
			if err != nil {
				return err
			}
			if n != len(b) {
				return err
			}

			if err := nf.Chmod(0644); err != nil {
				return err
			}
			nf.Close()
			if err := os.Rename(nf.Name(), translationsFile); err != nil {
				return err
			}

			return nil
		},
	}
}

// -- 1: key (default)
// L:  Localize(in string, data ...any) string
// N:  LocalizePlural(in string, count int, data ...any) string
// -- 2: key, default
// X:  LocalizeKey(key, fallback string, data ...any) string
// XN: LocalizePluralKey(key, fallback string, count int, data ...any) string

func updateTranslations(m map[string]interface{}, srcs []string) error {
	for _, src := range srcs {
		fset := token.NewFileSet()
		af, err := parser.ParseFile(fset, src, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		ast.Walk(&visitor{m: m}, af)
	}

	return nil
}

type visitor struct {
	state    stateFn
	name     string
	key      string
	fallback string
	m        map[string]interface{}
}

// Visitor pattern using a state machine while walking the tree
func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	if v.state == nil {
		v.state = translateNew
		v.name = ""
		v.key = ""
		v.fallback = ""
	}

	v.state = v.state(v, node)

	return v
}

// State machine to pick out translation key and fallback from AST
type stateFn func(*visitor, ast.Node) stateFn

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

func translateCall(v *visitor, node ast.Node) stateFn {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return nil
	}

	v.name = ident.Name

	switch ident.Name {
	case "L", "Localise":
		return translateLocalize
	case "N", "LocalizePlural":
		return translateLocalize
	case "X", "LocalizeKey":
		return translateKey
	case "XN", "LocalizePluralKey":
		return translateKey
	}

	return nil
}

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

func translateFinish(v *visitor, node ast.Node) stateFn {
	// only adding new keys, ignoring changed or removed (ha!) ones
	// removing is dangerous as there could be dynamic keys that get removed

	_, found := v.m[v.key]
	if found {
		return nil
	}

	switch v.name {
	case "LocalizePlural", "LocalizePluralKey", "N", "XN":
		m := make(map[string]string)
		m["other"] = v.fallback
		v.m[v.key] = m
	default:
		v.m[v.key] = v.fallback
	}

	return nil
}
