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
		os.Exit(1)
	}

	// schemaFile := os.Args[1]
	schemaData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error reading schema file: %s, err: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)

	for _, yamlFile := range os.Args[2:] {
		yamlData, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			fmt.Printf("error reading yaml file: %s, err: %v\n", yamlFile, err)
			os.Exit(1)
		}

		objects := make(map[string]interface{}) // NOTE: it doesn't work if yaml is a list
		err = yaml.Unmarshal(yamlData, &objects)
		if err != nil {
			fmt.Printf("error unmarshaling yaml file: %s, err: %v\n", yamlFile, err)
			os.Exit(1)
		}

		documentLoader := gojsonschema.NewGoLoader(objects)
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			fmt.Printf("error validating yaml file: %s, err: %v\n", yamlFile, err)
			os.Exit(1)
		}

		if result.Valid() {
			return
		}

		fmt.Printf("file: %s is invalid: \n", yamlFile)
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}

		os.Exit(1)
	}
}
