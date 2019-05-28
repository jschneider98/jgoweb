package util

import (
	"fmt"
	"runtime"
	"strings"
	"github.com/gocraft/web"
)

// UniqueIntArray
// map used for O(1) performance
// struct{} used because it doesn't take up any extra space
type UniqueIntArray struct {
	Set map[int]struct{}
	Data []int
}

func NewUniqueIntArray() *UniqueIntArray {
	var newInt []int
	newMap := make(map[int]struct{})
	newArray := UniqueIntArray{newMap, newInt}

	return &newArray
}

func (a *UniqueIntArray) Append(value int) {

	if _, ok := a.Set[value]; ok {
		// element found
		return
	}

	// update set
	a.Set[value] = struct{}{}
	// update array
	a.Data = append(a.Data, value)
}

// Based on WhereAmI() by Jim Lawless
// https://github.com/jimlawless/whereami
func WhereAmI(depthList ...int) string {
	var depth int
	
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s~%s~%d", formatFilename(file), runtime.FuncForPC(function).Name(), line)
}

// return the source filename (without the extension) after the last slash
func formatFilename(original string) string {
	i := strings.LastIndex(original, "/")
	var filename string

	if i == -1 {
		filename = original
	} else {
		filename = original[i+1:]
	}

	return strings.Replace(filename, ".go", "", -1)
}

//
func GetBaseUrl(req *web.Request) string {
	scheme := req.URL.Scheme
	host := req.Host

	if scheme == "" {
		scheme = "http"
	}

	return scheme + "://" + host
}
