package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

func usage() {
	fmt.Println("usage: ")
	fmt.Println("    lint <schema-file> <yaml-files>...")
}

func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	// schemaFile := os.Args[1]
	schemaData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error reading schema file: %s, err: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)

	errs := []string{} // collect all errors in a round, to avoid run lint tool multiple times
	for _, yamlFile := range os.Args[2:] {
		yamlData, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			errs = append(errs, fmt.Sprintf("file: %s; readfile err: %v\n", yamlFile, err))
			continue
		}

		objects := make(map[string]interface{}) // NOTE: it doesn't work if yaml is a list
		err = yaml.Unmarshal(yamlData, &objects)
		if err != nil {
			errs = append(errs, fmt.Sprintf("file: %s; unmarshal err: %v\n", yamlFile, err))
			continue
		}

		documentLoader := gojsonschema.NewGoLoader(objects)
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			errs = append(errs, fmt.Sprintf("file: %s, schema validation err: %v\n", yamlFile, err))
			continue
		}

		if result.Valid() {
			continue
		}

		errs = append(errs, fmt.Sprintf("\nfile: %s is invalid: ", yamlFile))
		for _, desc := range result.Errors() {
			errs = append(errs, fmt.Sprintf("- %s", desc))
		}
	}

	if len(errs) > 0 {
		for _, v := range errs {
			fmt.Println(v)
		}
		os.Exit(1)
	}
}
