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
func (rg *ReportGenerator) GetStructCode() string {
	var code string

	code += fmt.Sprintf(
"type %sReport struct {\n" +
"	Ctx jgoweb.ContextInterface `json:\"-\" validate:\"-\"`\n" +
"	User *models.User\n" +
"}\n\n", rg.BaseStructName)


	code += fmt.Sprintf("type %sResult struct {\n", rg.BaseStructName)

	for key := range rg.Fields {
		field := rg.Fields[key]
		parts := strings.Split(field, ".")

		if len(parts) > 1 {
			field = parts[1]
		}

		code += fmt.Sprintf("\t%s sql.NullString `json:\"%s\"`\n", util.ToCamelCase(field), field)
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

	params := make(map[string]string)
	params["~Name~"] = rg.BaseStructName

	code += util.NamedSprintf(`
//
func (r *~Name~Report) Run(params url.Values) ([]~Name~Result, error) {
	var results []~Name~Result
	var query string
	var err error

	strParams := make(map[string]string)
	qParams := make(map[string]string)

	qParams["@AccountId@"] = r.User.AccountId.String

`, params)

	code += `	strParams["~limit~"] = ""
	strParams["~offset~"] = ""

	if params.Get("limit") != "" {
		strParams["~limit~"] = fmt.Sprintf("LIMIT %s", params.Get("limit"))
	}

	if params.Get("offset") != "" {
		strParams["~offset~"] = fmt.Sprintf("OFFSET %s", params.Get("offset"))
	}
`

	for key := range rg.Fields {
		field := rg.Fields[key]
		parts := strings.Split(field, ".")

		if len(parts) > 1 {
			field = parts[1]
		}

		code += fmt.Sprintf(`
	strParams["~%s~"] = ""

	if params.Get("%s") != "" {
		strParams["~%s~"] = "AND %s ilike @%s@"
		qParams["@%s@"] = params.Get("%s") + "%s"
	}

`, field, field, field, field, field, field, field, "%")
	
	}

	code += "\tquery = util.NamedSprintf(`<query here>`, strParams)"
	code += `

	query, sqlParams, err := util.PrepareQuery(query, qParams)

	stmt := r.Ctx.SelectBySql(query, sqlParams...)

	_, err = stmt.Load(&results)

	if err !=nil {
		return nil, err
	}

	return results, nil
}
`

	return code
}
