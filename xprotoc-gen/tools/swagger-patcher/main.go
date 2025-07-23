package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <swagger.(json|yaml)|dir> ...\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	for _, arg := range os.Args[1:] {
		matches, _ := filepath.Glob(arg)
		if len(matches) == 0 {
			matches = []string{arg}
		}
		for _, f := range matches {
			info, err := os.Stat(f)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			if info.IsDir() {
				filepath.WalkDir(f, func(p string, d fs.DirEntry, _ error) error {
					if d.IsDir() {
						return nil
					}
					ext := strings.ToLower(filepath.Ext(p))
					if ext == ".json" || ext == ".yaml" || ext == ".yml" {
						if err := fixFile(p); err != nil {
							fmt.Fprintln(os.Stderr, err)
						} else {
							fmt.Println("patched", p)
						}
					}
					return nil
				})
			} else if err := fixFile(f); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println("patched", f)
			}
		}
	}
}

func patch(sw map[string]interface{}) {
	defs, _ := sw["definitions"].(map[string]interface{})
	paths, _ := sw["paths"].(map[string]interface{})

	for _, p := range paths {
		methods, _ := p.(map[string]interface{})
		for _, raw := range methods {
			op, _ := raw.(map[string]interface{})

			var param map[string]interface{}
			if arr, ok := op["parameters"].([]interface{}); ok {
				for _, it := range arr {
					pmap, _ := it.(map[string]interface{})
					if pmap["in"] == "body" {
						param = pmap
						break
					}
				}
			}
			if param == nil {
				continue
			}

			sch, _ := param["schema"].(map[string]interface{})
			ref, _ := sch["$ref"].(string)
			if ref == "" || !strings.HasPrefix(ref, "#/definitions/") {
				continue
			}
			defName := strings.TrimPrefix(ref, "#/definitions/")
			defObj, _ := defs[defName].(map[string]interface{})
			if !hasBytesField(defObj) {
				continue
			}

			param["in"] = "formData"
			param["type"] = "file"
			delete(param, "schema")
			delete(param, "format")

			const ct = "multipart/form-data"
			if cons, ok := op["consumes"].([]interface{}); ok {
				op["consumes"] = append(cons, ct)
			} else {
				op["consumes"] = []interface{}{ct}
			}
		}
	}
}

func hasBytesField(def map[string]interface{}) bool {
	props, _ := def["properties"].(map[string]interface{})
	for _, v := range props {
		p, _ := v.(map[string]interface{})
		if p["type"] == "string" {
			if fmt, _ := p["format"].(string); fmt == "bytes" || fmt == "binary" {
				return true
			}
		}
	}
	return false
}

func fixFile(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ext := strings.ToLower(filepath.Ext(path))
	var jsonData []byte

	switch ext {
	case ".yaml", ".yml":
		jsonData, err = yaml.YAMLToJSON(src)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	case ".json":
		jsonData = src
	default:
		return nil
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	patch(doc)

	jsonOut, _ := json.MarshalIndent(doc, "", "  ")

	var out []byte
	if ext == ".json" {
		out = jsonOut
	} else {
		if out, err = yaml.JSONToYAML(jsonOut); err != nil {
			return err
		}
	}
	return os.WriteFile(path, out, 0o644)
}
