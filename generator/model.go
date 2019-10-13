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
func NewModelGenerator(ctx jgoweb.ContextInterface, schema string, table string) (*ModelGenerator, error) {
	var err error
	
	mg := &ModelGenerator{
		Schema: schema,
		Table: table,
		Ctx: ctx,
	}

	mg.MakeModelName()

	mg.Fields, err = psql.GetFields(mg.Ctx, mg.Schema, mg.Table)

	if err != nil {
		return nil, err
	}

	return mg, nil
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
func (mg *ModelGenerator) GetAnnotation(field psql.Field) string {
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
		re := regexp.MustCompile("[0-9]+")
		val += ",stringlength(1|" + re.FindString(field.DataType) + ")"
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
func (mg *ModelGenerator) GetInstanceName() string {
	return strings.ToLower(mg.ModelName)
}

//
func (mg *ModelGenerator) GetFullTableName() string {
	return mg.Schema + "." + mg.Table
}

//
func (mg *ModelGenerator) GetUnsetPkeyVal() string {
	var unsetPkeyVal string

	if mg.Fields[0].DataType == "integer" {
		unsetPkeyVal = "0"
	} else {
		unsetPkeyVal = `""`
	}

	return unsetPkeyVal
}

//
func (mg *ModelGenerator) GetStructInstanceName() string{
	var structInstance string

	re := regexp.MustCompile("[A-Z]+")
	letters := re.FindAllString(mg.ModelName, -1)

	for key := range letters {
		structInstance += letters[key]
	}

	if structInstance == "" {
		structInstance = "m"
	}

	return strings.ToLower(structInstance)
}

//
func (mg *ModelGenerator) Generate() string {
	var code string


	code = mg.GetImportCode()
	code += mg.GetStructCode()

	code += mg.GetHydratorStructCode()
	code += mg.GetHydratorIsValidCode()

	code += mg.GetNewCode()
	code += mg.GetFetchByIdCode()
	code += mg.GetIsValidCode()
	code += mg.GetSaveCode()
	code += mg.GetInsertCode()
	code += mg.GetUpdateCode()
	code += mg.GetDeleteCode()

	return code
}

// 
func (mg *ModelGenerator) GetImportCode() string {
return	`
package models

import(
	"github.com/gocraft/dbr"
	"github.com/asaskevich/govalidator"
	"github.com/jschneider98/jgoweb"
)

`
}

//
func (mg *ModelGenerator) GetStructCode() string {
	var code string

	code += fmt.Sprintf("// %s\n", mg.ModelName)
	code += fmt.Sprintf("type %s struct {\n", mg.ModelName)

	for key := range mg.Fields {
		code += fmt.Sprintf("\t%s %s %s\n", mg.ToCamelCase(mg.Fields[key].FieldName), mg.ConvertDataType(mg.Fields[key]), mg.GetAnnotation(mg.Fields[key]))
	}

	code += "\tCtx ContextInterface `json:\"-\" valid:\"-\"`\n"
	code += "}\n"

	return code
}


//
func (mg *ModelGenerator) GetHydratorStructCode() string {
	var code string

	code += fmt.Sprintf("// %s\n", mg.ModelName + "Hydrator")
	code += fmt.Sprintf("type %s struct {\n", mg.ModelName + "Hydrator")

	for key := range mg.Fields {
		code += fmt.Sprintf("\t%s %s %s\n", mg.ToCamelCase(mg.Fields[key].FieldName), "string", mg.GetAnnotation(mg.Fields[key]))
	}

	code += "}\n"

	return code
}

//
func (mg *ModelGenerator) GetHydratorIsValidCode() string {
	var code string
	structInstance := mg.GetStructInstanceName() + "h"

	code += fmt.Sprintf(`
// Validate the hydrator
func (%s *%s) isValid() (bool, error) {
	return govalidator.ValidateStruct(%s)
}
`, structInstance, mg.ModelName + "Hydrator", structInstance)

	return code
}

//
func (mg *ModelGenerator) GetNewCode() string {
	var code string
	code += fmt.Sprintf(`
// Empty new model
func New%s(ctx ContextInterface) *%s {
	return &%s{Ctx: ctx}
}
`, mg.ModelName, mg.ModelName, mg.ModelName)

	return code
}

//
func (mg *ModelGenerator) GetFetchByIdCode() string {
	var code string
	fullTableName := mg.GetFullTableName()
	instanceName := mg.GetInstanceName()

	code += fmt.Sprintf(`
// Factory Method
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

	return code
}

//
func (mg *ModelGenerator) GetIsValidCode() string {
	var code string
	structInstance := mg.GetStructInstanceName()

	code += fmt.Sprintf(`
// Validate the model
func (%s *%s) isValid() (bool, error) {
	return govalidator.ValidateStruct(%s)
}
`, structInstance, mg.ModelName, structInstance)

	return code
}

//
func (mg *ModelGenerator) GetSaveCode() string {
	var code string
	structInstance := mg.GetStructInstanceName()
	unsetPkeyVal := mg.GetUnsetPkeyVal()

	code += fmt.Sprintf(`
// Insert/Update based on pkey value
func (%s *%s) Save() error {
	isValid, err := %s.isValid()

	if !isValid {
		return err
	}

	if %s.Id == %s {
		return %s.Insert()
	} else {
		return %s.Update()
	}
}
`, structInstance, mg.ModelName, structInstance, structInstance, unsetPkeyVal, structInstance, structInstance)


	return code
}

//
func (mg *ModelGenerator) GetInsertCode() string {
	var code string
	structInstance := mg.GetStructInstanceName()
	fullTableName := mg.GetFullTableName()
	columnList := ""

	for key := range mg.Fields {
		if (!mg.Fields[key].Default.Valid) {
			columnList += fmt.Sprintf(",\n\t\t\t\"%s\"", mg.Fields[key].FieldName)
		}
	}

	code += fmt.Sprintf(`
// Insert a new record
func (%s *%s) Insert() error {
	tx, err := %s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.InsertInto("%s").
		Columns(%s).
		Record(%s).
		Exec()

	if err != nil {
		return err
	}

	err = %s.Ctx.OptionalCommit(tx)

	return err
}
`, structInstance, mg.ModelName, structInstance, fullTableName, columnList, structInstance, structInstance)

	return code
}

//
func (mg *ModelGenerator) GetUpdateCode() string {
	var code string
	structInstance := mg.GetStructInstanceName()
	fullTableName := mg.GetFullTableName()
	columnList := ""

	for key := range mg.Fields {
		columnList += fmt.Sprintf("\t\tSet(\"%s\", %s.%s).\n", mg.Fields[key].FieldName, structInstance, mg.ToCamelCase(mg.Fields[key].FieldName))
	}

	code += fmt.Sprintf(`
// Update a record
func (%s *%s) Update() error {
	tx, err := %s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.Update("%s").
%s
		Where("id = ?", %s.Id).
		Exec()

	if err != nil {
		return err
	}

	err = %s.Ctx.OptionalCommit(tx)

	return err
}
`, structInstance, mg.ModelName, structInstance, fullTableName, columnList, structInstance, structInstance)

	return code
}

//
func (mg *ModelGenerator) GetDeleteCode() string {
	var code string
	structInstance := mg.GetStructInstanceName()
	fullTableName := mg.GetFullTableName()
	softDelete := false

	for key := range mg.Fields {
		if mg.Fields[key].FieldName == "deleted_at" {
			softDelete = true
		}
	}

	if softDelete {
		return mg.GetSoftDeleteCode()
	}

	code += fmt.Sprintf(`
// Hard delete a record
func (%s *%s) Delete() error {
	tx, err := %s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.Delete("%s").
		Where("id = ?", %s.Id).
		Exec()

	if err != nil {
		return err
	}

	return p.Ctx.OptionalCommit(tx)
`, structInstance, mg.ModelName, structInstance, fullTableName, structInstance)

	return code
}

//
func (mg *ModelGenerator) GetSoftDeleteCode() string {
	var code string
	structInstance := mg.GetStructInstanceName()
	fullTableName := mg.GetFullTableName()


	code += fmt.Sprintf(`
// Soft delete a record
func (%s *%s) Delete() error {
	tx, err := %s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.Update("%s").
		Set("deleted_at = ?", "timezone('utc'::text, now())"").
		Where("id = ?", %s.Id).
		Exec()

	if err != nil {
		return err
	}

	return p.Ctx.OptionalCommit(tx)
}
`, structInstance, mg.ModelName, structInstance, fullTableName, structInstance)

	return code
}
