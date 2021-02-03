package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

func usage() {
	fmt.Println("usage: ")
	fmt.Println("    unique <yaml-files>...")
}

func doneOrDieWithMessage(err error, msg string) {
	if err == nil {
		return
	}

	fmt.Println("unexpected error:", err)
	os.Exit(1)
}

// inplace sort files by ModTime
func sortFilesByMTime(files []string) {
	// it's hard to do this in shell

	sort.Slice(files, func(i, j int) bool {
		statI, err := os.Stat(files[i])
		doneOrDieWithMessage(err, fmt.Sprintf("file: %s; stat err: %v", files[i], err))
		statJ, err := os.Stat(files[j])
		doneOrDieWithMessage(err, fmt.Sprintf("file: %s; stat err: %v", files[j], err))

		return statI.ModTime().Before(statJ.ModTime())
	})
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	files := os.Args[1:]
	sortFilesByMTime(files)

	names := make(map[string]string) // map[case-name]filename
	errMsgs := []string{}
	for _, filename := range files {
		// TODO: path of name depends on schema, maybe pass key path via os.args
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("file: %s; readfile err: %v", filename, err))
			continue
		}

		meta := make(map[string]interface{})
		err = yaml.Unmarshal(data, &meta)
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("file: %s; unmarshal err: %v", filename, err))
			continue
		}

		if _, ok := meta["name"]; !ok {
			errMsgs = append(errMsgs, fmt.Sprintf("file: %s; missing expected key: [name]", filename))
			continue
		}

		name, ok := meta["name"].(string)
		if !ok {
			errMsgs = append(errMsgs, fmt.Sprintf("file: %s; invalid type of [name], expect string", filename))
			continue
		}

		if _, ok := names[name]; ok {
			errMsgs = append(errMsgs, fmt.Sprintf("file: %s; name in %s conflicts with name in %s", filename, filename, names[name]))
			continue
		} else {
			names[name] = filename
		}
	}

	if len(errMsgs) > 0 {
		for _, v := range errMsgs {
			fmt.Println(v)
		}
		os.Exit(1)
	}
}
