package picture

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sbinet/go-python"
)

func init() {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
}

// pyStr python字符串
var pyStr = python.PyString_FromString

// goStr go字符串
var goStr = python.PyString_AS_STRING

// FillWordsIntoPic 字符填充图片
func FillWordsIntoPic(oriImgPath, resultImgPath, filledWords, fontPath string, scale, step int) error {
	if filledWords == "" {
		return errors.New("filledWords can not be empty")
	}
	curLoc, _ := os.Getwd()
	curLoc += fmt.Sprintf("%c%s", os.PathSeparator, "scripts")
	insertBeforeSysPath("/usr/lib/python2.7/site-packages")
	log.Println(curLoc)
	module := importModule(curLoc, "fillwordsintopic")
	if module == nil {
		return errors.New("can not import package 'fillwordsintopic'")
	}
	fillFunc := module.GetAttrString("fillWordsIntoPic")
	if fillFunc == nil {
		return errors.New("could not getattr 'fillWordsIntoPic'")
	}
	defer fillFunc.DecRef()
	excResult := fillFunc.CallFunction(oriImgPath, resultImgPath, filledWords, fontPath, python.PyInt_FromLong(scale), python.PyInt_FromLong(step))
	if excResult == nil {
		return errors.New("could not call 'values.myfunc()'")
	}
	if goStr(excResult) != "Success" {
		return errors.New("unexpect appear when execute script, error is" + goStr(excResult))
	}
	excResult.DecRef()
	return nil
}

// insertBeforeSysPath will add given dir to python import path
func insertBeforeSysPath(p string) string {
	sysModule := python.PyImport_ImportModule("sys")
	path := sysModule.GetAttrString("path")
	python.PyList_Insert(path, 0, pyStr(p))
	return goStr(path.Repr())
}

// importModule will import python module from given directory
func importModule(dir, name string) *python.PyObject {
	sysModule := python.PyImport_ImportModule("sys") // import sys
	path := sysModule.GetAttrString("path")          // path = sys.path
	python.PyList_Insert(path, 0, pyStr(dir))        // path.insert(0, dir)
	return python.PyImport_ImportModule(name)        // return __import__(name)
}
