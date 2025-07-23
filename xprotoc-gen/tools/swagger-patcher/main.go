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
	paths, _ := sw["paths"].(map[string]interface{})
	const multipart = "multipart/form-data"

	for _, rawPathItem := range paths {
		ops, _ := rawPathItem.(map[string]interface{})
		for _, rawOp := range ops {
			op, _ := rawOp.(map[string]interface{})

			// 1) ensure consumes: multipart/form-data
			if cons, ok := op["consumes"].([]interface{}); ok {
				found := false
				for _, c := range cons {
					if c == multipart {
						found = true
						break
					}
				}
				if !found {
					op["consumes"] = append(cons, multipart)
				}
			} else {
				op["consumes"] = []interface{}{multipart}
			}

			// 2) rewrite the single body→binary parameter
			if params, ok := op["parameters"].([]interface{}); ok {
				for i, rawParam := range params {
					p, _ := rawParam.(map[string]interface{})
					if p["in"] == "body" {
						if schema, ok := p["schema"].(map[string]interface{}); ok {
							if schema["type"] == "string" &&
								(schema["format"] == "binary" || schema["format"] == "bytes") {
								// mutate into a file upload
								delete(p, "schema")
								p["in"] = "formData"
								p["type"] = "file"
								// keep any existing description
								// p["description"] = schema["description"]
								params[i] = p
							}
						}
					}
				}
				op["parameters"] = params
			}
		}
	}
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
