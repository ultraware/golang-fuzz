package main

import (
	"os"
	"strings"
	"text/template"
)

func createTemplate(tmplStr, fileName, pkgName, funcName, inputType string) func() {
	tmpl, err := template.New(``).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	inputCode, inputLen, imprts := getInputCode(inputType)

	var fuzzFile *os.File
	if strings.Contains(fileName, `*`) {
		fuzzFile, err = os.CreateTemp(`.`, fileName)
	} else {
		fuzzFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	}
	if err != nil {
		panic(err)
	}
	defer fuzzFile.Close()

	cleanup := func() {
		if !*keepFile {
			_ = os.Remove(fuzzFile.Name())
		}
	}

	err = tmpl.Execute(fuzzFile, tmplData{
		PkgName:   pkgName,
		Imports:   imprts,
		FuncName:  funcName,
		CorpusDir: *corpusDir,
		InputType: inputType,
		InputCode: inputCode,
		InputLen:  inputLen,
	})
	if err != nil {
		cleanup()
		panic(err)
	}

	return cleanup
}
