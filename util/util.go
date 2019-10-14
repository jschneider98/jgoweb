package util

import (
	"regexp"
	"fmt"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
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


// my_var_id => MyVarId || my var id = MyVarId, etc.
func ToCamelCase(val string) string {

	if val == "" {
		return ""
	}

	val = strings.ToLower(val)
	val = strings.ReplaceAll(val, "_", " ")
	val = strings.Title(val)
	val = strings.ReplaceAll(val, " ", "")

	return val
}

// MyVarId => myVarId, my_var_id => myVarId, etc.
func ToLowerCamelCase(val string) string {

	if val == "" {
		return ""
	}

	val = ToCamelCase(val)

	rune, size := utf8.DecodeRuneInString(val)
	return string(unicode.ToLower(rune)) + val[size:]
}

// MyVarId => MVI
func ToAcronym(val string) string {

	if val == "" {
		return ""
	}

	re := regexp.MustCompile("[A-Z]+")
	letters := re.FindAllString(val, -1)

	return strings.Join(letters, "")
}

// MyVarId => mvi
func ToLowerAcronym(val string) string {

	if val == "" {
		return ""
	}

	val = ToAcronym(val)

	return strings.ToLower(val)
}
