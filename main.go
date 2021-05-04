package main

import (
	"encoding/json"
	"github.com/Slidem/ftreedepth"
	"github.com/Slidem/inplaceenvsubst"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const defaultDepth = 15

type SearchInput struct {
	Patterns []string `json:"patterns"`
	Files    []string `json:"files"`
	Depth    int      `json:"depth"`
}

func (si *SearchInput) FindAll() bool {

	return len(si.Patterns) == 0 && len(si.Files) == 0
}

func main() {

	cfg := &inplaceenvsubst.Config{
		FailOnMissingVariables: failOnMissingVariables(),
		RunInParallel:          replaceInParallel(),
		ErrorListener:          getErrorListener(),
	}
	inplaceenvsubst.ProcessFiles(getFilesToReplace(), cfg)
}

func getFilesToReplace() []string {
	si := getSearchInput()
	var toReplace []string
	ftreedepth.WalkTree(si.Depth, getWorkingDirectory(), walkTreeFunc(si, &toReplace))
	return toReplace
}

func walkTreeFunc(si SearchInput, toReplace *[]string) ftreedepth.CallbackFunc {
	return 	func(path string, info os.FileInfo, err error) {
		fn := filepath.Base(path)
		if si.FindAll() {
			*toReplace = append(*toReplace, path)
			return
		}
		for _, r := range si.Patterns {
			if matched, _ := regexp.MatchString(r, fn); matched {
				*toReplace = append(*toReplace, path)
				return
			}
		}
		for _, n := range si.Files {
			if fn == n {
				*toReplace = append(*toReplace, path)
			}
		}
	}
}

func replaceInParallel() bool {
	v := os.Getenv("INPUT_REPLACE_IN_PARALLEL")
	if v == "true" {
		return true
	} else {
		return false
	}
}

func getWorkingDirectory() string {

	wd := os.Getenv("INPUT_WORKING-DIRECTORY")
	d := os.Getenv("GITHUB_WORKSPACE")
	if wd != "" {
		d = filepath.Join(d, wd)
	}
	return d
}

func getSearchInput() SearchInput {
	i := SearchInput{}
	inputJson := os.Getenv("INPUT_SEARCH_INPUT")
	if inputJson == "" {
	    inputJson = "{}"
	}

	err := json.Unmarshal([]byte(inputJson), &i)
	if err != nil {
		log.Fatalf("Could not convert search input to json value. Input: %s\n", inputJson)
	}
	if i.Depth <= 0 {
		i.Depth = defaultDepth
	}
	return i
}

func failOnMissingVariables() bool {
	return os.Getenv("INPUT_FAIL_ON_MISSING_VARIABLES") == "true"
}

func getErrorListener() inplaceenvsubst.ErrorListener {
	if failOnMissingVariables() {
		return &ExitOnMissingEnv{}
	} else {
		return &inplaceenvsubst.ConsoleErrorListener{}
	}
}

type ExitOnMissingEnv struct {
}

func (l *ExitOnMissingEnv) ErrorFound(filepath string, err error) {
	log.Fatal(err)
}
