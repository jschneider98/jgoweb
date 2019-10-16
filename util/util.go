package util

import (
	"regexp"
	"fmt"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
	"html/template"
	"github.com/gocraft/web"
	"gopkg.in/go-playground/validator.v9"
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

// MyVarId => my_far_id, etc
func ToSnakeCase(val string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")

	val = matchFirstCap.ReplaceAllString(val, "${1}_${2}")

	return strings.ToLower(val)
}

// MyVarId => My Var Id
func ToWords(val string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")

	return matchFirstCap.ReplaceAllString(val, "${1} ${2}")
}

//
func GetHtmlAlerts(msgType string, messages ...string) template.HTML {
	var msgs string

	if messages == nil {
		return template.HTML("")
	}

	msgs =  fmt.Sprintf("<div class=\"alert alert-%s\" role=\"alert\">\n", msgType)

	for key := range messages {
		msgs += fmt.Sprintf("\t%s</br>\n", messages[key])
	}

	msgs += "</div>"

	return template.HTML(msgs)
}

//
func GetNiceErrorMessage(errs error, seperator string) string {
	var msg []string
	var field string

	for _, err := range errs.(validator.ValidationErrors) {
		field = ToWords(err.Field())

		switch err.Tag() {
		case "required":
			msg = append(msg, fmt.Sprintf("%s is required.", field))
		case "email":
			msg = append(msg, fmt.Sprintf("%s must be a valid email address", field))
		case "max":
			msg = append(msg, fmt.Sprintf("%s is too long.", field))
		case "min":
			msg = append(msg, fmt.Sprintf("%s is too short.", field))
		case "rfc3339", "rfc3339WithoutZone":
			msg = append(msg, fmt.Sprintf("%s must be a valid date/time", field))
		case "date":
			msg = append(msg, fmt.Sprintf("%s must be a valid date", err.Field()))
		case "int":
			msg = append(msg, fmt.Sprintf("%s must be a valid whole number", field))
		case "float":
			msg = append(msg, fmt.Sprintf("%s must be a valid decimal number", field))
		case "notNull":
			msg = append(msg, fmt.Sprintf("%s must not be blank.", field))
		default:
			msg = append(msg, fmt.Sprintf("%s is invalid.", field))
		}

		
			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			// fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()
	}

	return strings.Join(msg, seperator)
}
