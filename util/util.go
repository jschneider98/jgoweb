package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"gopkg.in/go-playground/validator.v9"
	template "html/template"
	"math"
	"regexp"
	"runtime"
	"strings"
	textTemplate "text/template"
	"time"
	"unicode"
	"unicode/utf8"
)

// UniqueIntArray
// map used for O(1) performance
// struct{} used because it doesn't take up any extra space
type UniqueIntArray struct {
	Set  map[int]struct{}
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
	// var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	// val = matchFirstCap.ReplaceAllString(val, "${1}_${2}")

	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	val = matchAllCap.ReplaceAllString(val, "${1}_${2}")

	return strings.ToLower(val)
}

// MyVarId => My Var Id
func ToWords(val string) string {
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	return matchAllCap.ReplaceAllString(val, "${1} ${2}")
}

//
func GetAdvancedHtmlAlerts(msgType string, canClose bool, messages ...string) template.HTML {
	var msgs string

	if messages == nil {
		return template.HTML("")
	}

	msgs = fmt.Sprintf("<div class=\"alert alert-%s\" role=\"alert\">\n", msgType)

	if canClose {
		msgs += `
		<button type="button" class="close" data-dismiss="alert" aria-label="Close">
			<span aria-hidden="true">&times;</span>
		</button>
		`
	}

	for key := range messages {
		msgs += fmt.Sprintf("\t%s</br>\n", messages[key])
	}

	msgs += "</div>"

	return template.HTML(msgs)
}

//
func GetHtmlAlerts(msgType string, messages ...string) template.HTML {
	return GetAdvancedHtmlAlerts(msgType, true, messages...)
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
		case "errorMsg":
			msg = append(msg, fmt.Sprintf("%v", err.Value()))
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

// var params map[string]string
// params = make(map[string]string)
// params["@test"@] = "one"
// params["@test2@"] = "two"
// params["@test3@"] = "three"
// str, newParams, err := util.PrepareString("This is a @test@ @test2@ @test3@", params, "@", ?")
func PrepareString(str string, holders map[string]string, match string, replace string) (string, []interface{}, error) {
	var params []interface{}
	pattern := match + "[0-9A-Za-z_]+" + match

	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(str, -1)

	str = re.ReplaceAllLiteralString(str, replace)

	for key := range matches {
		val, ok := holders[matches[key]]

		if ok {
			params = append(params, val)
		}
	}

	if len(matches) != len(params) {
		err := errors.New("Parameter to placeholder mismatch in util.PrepareString")

		return str, params, err
	}

	return str, params, nil
}

//
func PrepareQuery(str string, holders map[string]string) (string, []interface{}, error) {
	return PrepareString(str, holders, "@", "?")
}

//
func PrepareQueryQuoted(str string, holders map[string]string) (string, error) {
	str, params, err := PrepareString(str, holders, "@", "'%s'")

	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(str, params...)

	return result, nil
}

//
func NamedSprintf(str string, holders map[string]string) string {
	var params []interface{}

	str, params, _ = PrepareString(str, holders, "~", "%s")

	return fmt.Sprintf(str, params...)
}

//
func TemplateToString(tmplStr string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)

	tmpl, err := textTemplate.New("test").Delims("{{{", "}}}").Parse(tmplStr)

	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, data)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

//
func HtmlTemplateToString(tmplStr string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.New("test").Delims("{{{", "}}}").Parse(tmplStr)

	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, data)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

//
func GetIsoCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// StrPad returns the input string padded on the left, right or both sides using padType to the specified padding length padLength.
// https://gist.github.com/asessa/3aaec43d93044fc42b7c6d5f728cb039 (Andrea Sessa)
//
// Example:
// input := "Codes";
// StrPad(input, 10, " ", "RIGHT")        // produces "Codes     "
// StrPad(input, 10, "-=", "LEFT")        // produces "=-=-=Codes"
// StrPad(input, 10, "_", "BOTH")         // produces "__Codes___"
// StrPad(input, 6, "___", "RIGHT")       // produces "Codes_"
// StrPad(input, 3, "*", "RIGHT")         // produces "Codes"
func StrPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

//
func JsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

//
func FindInStringArray(slice []string, val string) (int, bool) {

	for i, item := range slice {
		if item == val {
			return i, true
		}
	}

	return -1, false
}
