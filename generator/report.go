package generator

import (
	"fmt"
	"strings"
	"github.com/jschneider98/jgoweb/util"
)

type ReportGenerator struct {
	BaseStructName string `json:"model_name"`
	InstanceName string `json:"instance_name"`
	StructAcronym string
	Fields []string
}

//
func NewReportGenerator(baseStructName string, fields []string) (*ReportGenerator) {
	return &ReportGenerator{BaseStructName: baseStructName, Fields: fields}
}

//
func (rg *ReportGenerator) MakeInstanceName() {
	rg.InstanceName = util.ToLowerCamelCase(rg.BaseStructName + "Report")
}

//
func (rg *ReportGenerator) MakeStructInstanceName() {
	rg.StructAcronym = util.ToLowerAcronym(rg.BaseStructName + "Report")
}

//
func (rg *ReportGenerator) Generate() string {
	var code string

	code += rg.GetIncludeCode()
	code += rg.GetStructCode()
	code += rg.GetNewCode()
	code += rg.GetRunCode()
	code += rg.GetCountCode()
	code += rg.GetQueryCode()
	code += rg.GetParamsCode()

	return code
}

//
func (rg *ReportGenerator) GetIncludeCode() string {
	var code string

	code += fmt.Sprintf(`package reports

import(
	"fmt"
	"database/sql"
	"net/url"
	"github.com/jschneider98/jgoweb"
	"github.com/jschneider98/jgoweb/util"
	"github.com/jschneider98/medex/models"
)

`)

	return code
}

//
func (rg *ReportGenerator) GetField(str string) string {
	parts := strings.Split(str, " as ")

	// alias
	if len(parts) > 1 {
		return util.ToCamelCase(parts[1])
	}

	parts = strings.Split(str, ".")

	if len(parts) > 1 {
		return util.ToCamelCase(parts[1])
	}

	return str
}

func (rg *ReportGenerator) GetRooteField(str string) string {
	parts := strings.Split(str, " as ")
	return parts[0]
}

//
func (rg *ReportGenerator) GetStructCode() string {
	var code string

	code += fmt.Sprintf(
"type %sReport struct {\n" +
"	Ctx jgoweb.ContextInterface `json:\"-\" validate:\"-\"`\n" +
"	User *models.User\n" +
"}\n\n", rg.BaseStructName)


	code += fmt.Sprintf("type %sResult struct {\n", rg.BaseStructName)

	for key := range rg.Fields {
		field := rg.GetField(rg.Fields[key])

		code += fmt.Sprintf("\t%s sql.NullString `json:\"%s\"`\n", field, field)
	}

	code += "}\n\n"

	return code
}

//
func (rg *ReportGenerator) GetNewCode() string {
	params := make(map[string]string)

	params["~Name~"] = rg.BaseStructName + "Report"

	code := util.NamedSprintf(`
// Empty new report
func New~Name~(ctx jgoweb.ContextInterface, user *models.User) (*~Name~) {
	return &~Name~{Ctx: ctx, User: user}
}

`, params)

	return code
}

//
func (rg *ReportGenerator) GetRunCode() string {
	var code string

	code += fmt.Sprintf(`
//
func (r *%sReport) Run(params url.Values) ([]%sResult, error) {
	var results []%sResult
	var query string
	var sqlParams []interface{}
	var err error

	// make sure we generate the correct query for this method
	params.Set("count", "")

	query, sqlParams, err = r.GetQuery(params)

	if err != nil {
		return nil, err
	}

	stmt := r.Ctx.SelectBySql(query, sqlParams...)

	_, err = stmt.Load(&results)

	if err !=nil {
		return nil, err
	}

	return results, nil
}

`, rg.BaseStructName, rg.BaseStructName, rg.BaseStructName)

		return code
	}



//
func (rg *ReportGenerator) GetCountCode() string {
	var code string

	code += fmt.Sprintf(`
//
func (r *%sReport) GetCount(params url.Values) (int, error) {
	var count int
	var query string
	var sqlParams []interface{}
	var err error

	// make sure we generate the correct query for this method
	params.Set("count", "true")

	query, sqlParams, err = r.GetQuery(params)

	if err != nil {
		return 0, err
	}

	stmt := r.Ctx.SelectBySql(query, sqlParams...)

	_, err = stmt.Load(&count)

	if err !=nil {
		return 0, err
	}

	return count, nil
}

`, rg.BaseStructName)

		return code
}


//
func (rg *ReportGenerator) GetQueryCode() string {
	var code string
	var placeholderStr string

	for key := range rg.Fields {
		field := rg.GetField(rg.Fields[key])

		placeholderStr += fmt.Sprintf("\t~%s~\n", field)
	}

	code += fmt.Sprintf("\n" +
"//\n" +
"func (r *%sReport) GetQuery(params url.Values) (string, []interface{}, error) {\n" +
"	var query string\n\n" +
"	strParams, qParams := r.GetQueryParams(params)\n\n" +
"	query = util.NamedSprintf(`\n" +
`SELECT
~fields~
< Query Body Here>
WHERE account_id = @AccountId@
%s
~order_by~
~limit~
~offset~
`+ "`, strParams)" +
`
	return util.PrepareQuery(query, qParams)
}

`, rg.BaseStructName, placeholderStr)

		return code
}


//
func (rg *ReportGenerator) GetParamsCode() string {
	var code string
	var fieldStr string
	var strParamsStr string

	for key := range rg.Fields {
		field := rg.GetField(rg.Fields[key])
		rootField := rg.GetRooteField(rg.Fields[key])

		fieldStr += fmt.Sprintf("\t%s,\n", rg.Fields[key])

		strParamsStr += fmt.Sprintf(`
	strParams["~%s~"] = ""

	if params.Get("%s") != "" {
		strParams["~%s~"] = "AND %s ilike @%s@"
		qParams["@%s@"] = "%s" + params.Get("%s") + "%s"
	}

`, field, field, field, rootField, field, field, "%", field, "%")
	}

	code += fmt.Sprintf(`
//
func (r *%sReport) GetQueryParams(params url.Values) (map[string]string, map[string]string) {
	strParams := make(map[string]string)
	qParams := make(map[string]string)

	qParams["@AccountId@"] = r.User.AccountId.String

`, rg.BaseStructName)

	code += `	strParams["~limit~"] = ""
	strParams["~offset~"] = ""
	strParams["~order_by~"] = ""
	strParams["~fields~"] = "\tcount(*) as count\n"

	if params.Get("count") == "" {

		if params.Get("limit") != "" {
			strParams["~limit~"] = fmt.Sprintf("LIMIT %s", params.Get("limit"))
		}

		if params.Get("offset") != "" {
			strParams["~offset~"] = fmt.Sprintf("OFFSET %s", params.Get("offset"))
		}
`

	code += fmt.Sprintf(`
		strParams["~fields~"] = ` + "`" + `
%s
` + "`", fieldStr)

	code += `
		strParams["~order_by~"] = "<ORDER BY>"

	} else {
		strParams["~limit~"] = "LIMIT 1"
	}
`

	code += strParamsStr
	code += fmt.Sprintf("\treturn strParams, qParams\n")
	code += "}\n"

	return code
}
