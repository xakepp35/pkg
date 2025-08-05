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
		for method, rawOp := range ops {
			op, _ := rawOp.(map[string]interface{})
			method = strings.ToLower(method)

			params, _ := op["parameters"].([]interface{})
			isBinaryUpload := false

			// Определяем, есть ли binary поле в body
			for _, rawParam := range params {
				p, _ := rawParam.(map[string]interface{})
				if p["in"] == "body" {
					if schema, ok := p["schema"].(map[string]interface{}); ok {
						if schema["type"] == "string" &&
							(schema["format"] == "binary" || schema["format"] == "bytes") {
							isBinaryUpload = true
							break
						}
					}
				}
			}

			// 1) Binary POST → formData
			if method == "post" && isBinaryUpload {
				op["consumes"] = []interface{}{multipart}
				var newParams []interface{}
				for _, rawParam := range params {
					p, _ := rawParam.(map[string]interface{})
					if p["in"] == "body" {
						// binary → file
						if schema, ok := p["schema"].(map[string]interface{}); ok {
							if schema["type"] == "string" &&
								(schema["format"] == "binary" || schema["format"] == "bytes") {
								delete(p, "schema")
								p["in"] = "formData"
								p["type"] = "file"
							}
						}
					}
					if p["name"] == "file_name" {
						p["in"] = "formData"
						p["type"] = "string"
					}
					newParams = append(newParams, p)
				}
				op["parameters"] = newParams
				continue
			}

			// 2) GET с google.api.HttpBody → application/octet-stream
			if method == "get" {
				if produces, ok := op["produces"].([]interface{}); !ok || len(produces) == 0 {
					op["produces"] = []interface{}{"application/octet-stream"}
				}
				delete(op, "consumes")
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
