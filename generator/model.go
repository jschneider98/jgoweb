package generator

import (
	"fmt"
	"regexp"
	"strings"
	"github.com/jschneider98/jgoweb"
	"github.com/jschneider98/jgoweb/db/psql"
)

type ModelGenerator struct {
	Schema string `json:"schema"`
	Table string `json:"table"`
	ModelName string `jsong:"model_name"`
	Fields []psql.Field
	Ctx jgoweb.ContextInterface `json:"-"`
}

//
func NewModelGenerator(ctx jgoweb.ContextInterface, schema string, table string) *ModelGenerator {
	
	mg := &ModelGenerator{
		Schema: schema,
		Table: table,
		Ctx: ctx,
	}

	mg.MakeModelName()


	return mg
}

//
func (mg *ModelGenerator) MakeModelName() {
	mg.ModelName = ""
	
	if mg.Schema != "public" {
		mg.ModelName = mg.ToCamelCase(mg.Schema)
	}

	mg.ModelName += mg.ToCamelCase(mg.Table)

}

//
func (mg *ModelGenerator) ConvertDataType(field psql.Field) string {

	switch field.DataType {
	case "integer":
		if field.NotNull == true  {
			return "int"
		} else {
			return "dbr.NullInt"
		}
	default:
		if field.NotNull == true  {
			return "string"
		} else {
			return "dbr.NullString"
		}
	}
}

//
func (mg *ModelGenerator) GetReflection(field psql.Field) string {
	return fmt.Sprintf("`json:\"%s\" %s`", field.FieldName, mg.GetValidation(field))
}

//
func (mg *ModelGenerator) GetValidation(field psql.Field) string {
	val := "valid:"

	if field.NotNull == true && !field.Default.Valid {
		val += `"required`
	} else {
		val += `"optional`
	}

	switch field.DataType {
	case "integer":
		val += ",int"
	case "uuid":
		val += ",uuid"
	case "timestamp with time zone":
		val += ",rfc3339"
	case "timestamp without time zone":
		val += ",rfc3339WithoutZone"
	}

	if strings.HasPrefix(field.DataType, "character varying(") {
		re:=regexp.MustCompile("[0-9]+")
		val += ",length(1|" + re.FindString(field.DataType) + ")"
	}

	val += `"`

	return val
}

//
func (mg *ModelGenerator) ToCamelCase(val string) string {
	val = strings.ToLower(val)
	val = strings.ReplaceAll(val, "_", " ")
	val = strings.Title(val)
	val = strings.ReplaceAll(val, " ", "")

	return val
}

//
func (mg *ModelGenerator) Generate() (string, error) {
	var code string
	var err error

	instanceName := strings.ToLower(mg.ModelName)
	fullTableName := mg.Schema + "." + mg.Table

	mg.Fields, err = psql.GetFields(mg.Ctx, mg.Schema, mg.Table)

	if err != nil {
		return "", err
	}

	code = `
package models

import(
	"github.com/gocraft/dbr"
	"github.com/asaskevich/govalidator"
	"github.com/jschneider98/jgoweb"
)

`
	code += fmt.Sprintf("// %s\n", mg.ModelName)
	code += fmt.Sprintf("type %s struct {\n", mg.ModelName)

	for key := range mg.Fields {
		code += fmt.Sprintf("\t%s %s %s\n", mg.ToCamelCase(mg.Fields[key].FieldName), mg.ConvertDataType(mg.Fields[key]), mg.GetReflection(mg.Fields[key]))
	}

	code += "\tCtx ContextInterface `json:\"-\" valid:\"-\"`\n"
	code += "}\n"

	code += fmt.Sprintf(`
//
func New%s(ctx ContextInterface) *%s {
	return &%s{Ctx: ctx}
}
`, mg.ModelName, mg.ModelName, mg.ModelName)

	code += fmt.Sprintf(`
// 
func Fetch%sById(ctx ContextInterface, id string) (*%s, error) {
	var %s []%s

	stmt := ctx.Select("*").
	From("%s").
	Where("id = ?", id).
	Limit(1)

	_, err := stmt.Load(&%s)

	if err != nil {
		return nil, err
	}

	if (len(%s) == 0) {
		return nil, nil
	}

	%s[0].Ctx = ctx

	return &%s[0], nil
}
`, mg.ModelName, mg.ModelName, instanceName, mg.ModelName, fullTableName, instanceName, instanceName, instanceName, instanceName)


	return code, nil
}
