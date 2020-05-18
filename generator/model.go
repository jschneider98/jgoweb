package generator

import (
	"fmt"
	"github.com/jschneider98/jgomodel"
	"github.com/jschneider98/jgoweb"
	"github.com/jschneider98/jgoweb/db/psql"
	"github.com/jschneider98/jgoweb/util"
	"strings"
)

type ModelGenerator struct {
	ModelName     string `json:"model_name"`
	InstanceName  string `json:"instance_name"`
	TrimSuffix    string
	StructAcronym string
	Model         *jgomodel.Model
	Fields        []psql.Field
}

//
func NewModelGenerator(ctx jgoweb.ContextInterface, schema string, table string, trimSuffix string) (*ModelGenerator, error) {
	var err error

	mg := &ModelGenerator{}
	mg.Model, err = jgomodel.NewModel(ctx, schema, table)

	if err != nil {
		return nil, err
	}

	mg.TrimSuffix = trimSuffix

	mg.MakeModelName()
	mg.MakeInstanceName()
	mg.MakeStructInstanceName()
	mg.Fields = mg.Model.Fields

	return mg, nil
}

//
func (mg *ModelGenerator) MakeModelName() {

	if mg.ModelName != "" {
		return
	}

	mg.ModelName = ""

	if mg.Model.Schema != "public" {
		mg.ModelName = util.ToCamelCase(mg.Model.Schema)
	}

	// Conditionally remove a suffix
	table := strings.TrimSuffix(mg.Model.Table, mg.TrimSuffix)

	mg.ModelName += util.ToCamelCase(table)
}

//
func (mg *ModelGenerator) MakeInstanceName() {
	mg.MakeModelName()
	mg.InstanceName = util.ToLowerCamelCase(mg.ModelName)
}

//
func (mg *ModelGenerator) MakeStructInstanceName() {
	mg.MakeModelName()
	mg.StructAcronym = util.ToLowerAcronym(mg.ModelName)
}

//
func (mg *ModelGenerator) Generate() string {
	var code string

	code = mg.GetImportCode()
	code += mg.GetStructCode()

	// code += mg.GetHydratorStructCode()
	// code += mg.GetHydratorIsValidCode()

	code += mg.GetNewCode()
	code += mg.GetSetDefaultCode()
	code += mg.GetNewWithDataCode()
	code += mg.GetFetchByIdCode()
	code += mg.GetProcessSubmit()
	code += mg.GetHydrateCode()
	code += mg.GetIsValidCode()
	code += mg.GetSaveCode()
	code += mg.GetInsertCode()
	code += mg.GetUpdateCode()
	code += mg.GetDeleteCode()
	code += mg.GetUndeleteCode()
	code += mg.GetSetterGetterCode()

	return code
}

//
func (mg *ModelGenerator) GetImportCode() string {
	return `package models

import(
	"time"
	"database/sql"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb"
	"github.com/jschneider98/jgoweb/util"
)
`
}

//
func (mg *ModelGenerator) GetStructCode() string {
	var code string

	code += fmt.Sprintf("// %s\n", mg.ModelName)
	code += fmt.Sprintf("type %s struct {\n", mg.ModelName)

	for key := range mg.Fields {
		// code += fmt.Sprintf("\t%s %s %s\n", mg.Fields[key].FieldName, mg.Fields[key].DataType, mg.Fields[key].Annotation)
		code += fmt.Sprintf("\t%s %s %s\n", mg.Fields[key].FieldName, "sql.NullString", mg.Fields[key].Annotation)
	}

	code += "\tCtx jgoweb.ContextInterface `json:\"-\" validate:\"-\"`\n"
	code += "}\n\n"

	return code
}

//
func (mg *ModelGenerator) GetHydratorStructCode() string {
	var code string

	code += fmt.Sprintf("// %s\n", mg.ModelName+"Hydrator")
	code += fmt.Sprintf("type %s struct {\n", mg.ModelName+"Hydrator")

	for key := range mg.Fields {
		code += fmt.Sprintf("\t%s %s %s\n", mg.Fields[key].FieldName, "string", mg.Fields[key].Annotation)
	}

	code += "\tCtx jgoweb.ContextInterface `json:\"-\" validate:\"-\"`\n"

	code += "}\n"

	return code
}

//
func (mg *ModelGenerator) GetHydratorIsValidCode() string {
	var code string
	acronym := mg.StructAcronym + "h"

	code += fmt.Sprintf(`
// Validate the hydrator
func (%s *%s) IsValid() error {
	return %s.Ctx.GetValidator().Struct(%s)
}
`, acronym, mg.ModelName+"Hydrator", acronym, acronym)

	return code
}

//
func (mg *ModelGenerator) GetNewCode() string {
	var code string

	code += fmt.Sprintf(`
// Empty new model
func New%s(ctx jgoweb.ContextInterface) (*%s, error) {
	%s := &%s{Ctx: ctx}
	%s.SetDefaults()

	return %s, nil
}
`, mg.ModelName, mg.ModelName, mg.StructAcronym, mg.ModelName, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetSetDefaultCode() string {
	var code string
	var defaults string

	for key := range mg.Fields {

		if mg.Fields[key].DbDefault.Valid && !mg.Fields[key].DbDefaultIsFunc {
			defaults += fmt.Sprintf("\t%s.Set%s(\"%s\")\n", mg.StructAcronym, mg.Fields[key].FieldName, mg.Fields[key].Default)
		} else if mg.Fields[key].FieldName == "CreatedAt" || mg.Fields[key].FieldName == "UpdatedAt" {
			defaults += fmt.Sprintf("\t%s.Set%s(%s)\n", mg.StructAcronym, mg.Fields[key].FieldName, " time.Now().Format(time.RFC3339) ")
		}
	}

	code += fmt.Sprintf(`
// Set defaults
func (%s *%s) SetDefaults() {
%s
}
`, mg.StructAcronym, mg.ModelName, defaults)

	return code
}

//
func (mg *ModelGenerator) GetNewWithDataCode() string {
	var code string

	code += fmt.Sprintf(`
// New model with data
func New%sWithData(ctx jgoweb.ContextInterface, req *web.Request) (*%s, error) {
	%s, err := New%s(ctx)

	if err != nil {
		return nil, err
	}

	err = %s.Hydrate(req)

	if err != nil {
		return nil, err
	}

	return %s, nil
}
`, mg.ModelName, mg.ModelName, mg.StructAcronym, mg.ModelName, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetFetchByIdCode() string {
	var code string

	code += fmt.Sprintf(`
// Factory Method
func Fetch%sById(ctx jgoweb.ContextInterface, id string) (*%s, error) {
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
`, mg.ModelName, mg.ModelName, mg.StructAcronym, mg.ModelName, mg.Model.FullTableName, mg.StructAcronym, mg.StructAcronym, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetProcessSubmit() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~Words~"] = util.ToWords(mg.ModelName)

	code += util.NamedSprintf(`
//
func (~StructAcronym~ *~ModelName~) ProcessSubmit(req *web.Request) (string, bool, error) {
	err := ~StructAcronym~.Hydrate(req)

	if err != nil {
		return "", false, err
	}

	err = ~StructAcronym~.Ctx.GetValidator().Struct(~StructAcronym~)

	if err != nil {
		return util.GetNiceErrorMessage(err, "</br>"), false, nil
	}

	err = ~StructAcronym~.Save()

	if err != nil {
		return "", false, err
	}

	return "~Words~ saved.", true, nil
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetHydrateCode() string {
	var code string
	assignments := ""

	for key := range mg.Fields {

		assignments += fmt.Sprintf("\t%s.Set%s(req.PostFormValue(\"%s\"))\n", mg.StructAcronym, mg.Fields[key].FieldName, mg.Fields[key].FieldName)

		// fieldName := mg.Fields[key].FieldName

		// switch mg.Fields[key].DataType {
		// case "sql.NullString":
		// 	assignments += fmt.Sprintf("\tif %sHydrator.%s != \"\" {\n", mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t\t%s.%s.String = %sHydrator.%s\n", mg.InstanceName, fieldName, mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t\t%s.%s.Valid = true\n", mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t}\n\n")
		// case "sql.NullInt64":
		// 	assignments += fmt.Sprintf("\tif %sHydrator.%s != \"\" {\n", mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t%s.%s.Int64 = int64(%sHydrator.%s)\n", mg.InstanceName, fieldName, mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t%s.%s.Valid = true\n", mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t}\n\n")
		// case "slq.NullBool":
		// 	assignments += fmt.Sprintf("\tif %sHydrator.%s != \"\" {\n", mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t%s.%s.Bool = bool(%sHydrator.%s)\n", mg.InstanceName, fieldName, mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t%s.%s.Valid = true\n", mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t}\n\n")
		// case "slq.NullFloat64":
		// 	assignments += fmt.Sprintf("\tif %sHydrator.%s != \"\" {\n", mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t%s.%s.Float64 = float64(%sHydrator.%s)\n", mg.InstanceName, fieldName, mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t%s.%s.Valid = true\n\n", mg.InstanceName, fieldName)
		// 	assignments += fmt.Sprintf("\t}\n\n")
		// }
	}

	code += fmt.Sprintf(`
// Hydrate the model with data
func (%s *%s) Hydrate(req *web.Request) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}

%s
	return nil
}
`, mg.StructAcronym, mg.ModelName, assignments)

	return code
}

//
func (mg *ModelGenerator) GetIsValidCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName

	code += util.NamedSprintf(`
// Validate the model
func (~StructAcronym~ *~ModelName~) IsValid() error {
	return ~StructAcronym~.Ctx.GetValidator().Struct(~StructAcronym~)
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetSaveCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName

	code += util.NamedSprintf(`
// Insert/Update based on pkey value
func (~StructAcronym~ *~ModelName~) Save() error {
	err := ~StructAcronym~.IsValid()

	if err != nil {
		return err
	}

	if !~StructAcronym~.Id.Valid {
		return ~StructAcronym~.Insert()
	} else {
		return ~StructAcronym~.Update()
	}
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetInsertCode() string {
	var code string
	var objCols []string
	var colList string
	ph := make(map[string]string)

	//
	for key := range mg.Fields {
		if mg.Fields[key].DbFieldName != "id" && mg.Fields[key].DbFieldName != "created_at" && mg.Fields[key].DbFieldName != "updated_at" {
			// (p.AccountId, p.Units, ...)
			objCols = append(objCols, fmt.Sprintf("%s.%s", mg.StructAcronym, mg.Fields[key].FieldName))
		}
	}

	query := mg.Model.GetInsertQuery()

	colList = strings.Join(objCols, ",\t\t\t")
	colList = strings.ReplaceAll(colList, ",", ",\n")

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~colList~"] = colList

	code += util.NamedSprintf(`
// Insert a new record
func (~StructAcronym~ *~ModelName~) Insert() error {
	query := `+"`\n"+query+"\n`"+`

	stmt, err := ~StructAcronym~.Ctx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(~colList~).Scan(&~StructAcronym~.Id)

	if err != nil {
		return err
	}

	return nil
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetUpdateCode() string {
	var code string
	columnList := ""
	updatedAt := ""
	ph := make(map[string]string)

	for key := range mg.Fields {

		if mg.Fields[key].FieldName != "CreatedAt" {
			columnList += fmt.Sprintf("\t\tSet(\"%s\", %s.%s).\n", mg.Fields[key].DbFieldName, mg.StructAcronym, mg.Fields[key].FieldName)
		}

		if mg.Fields[key].FieldName == "UpdatedAt" {
			updatedAt = fmt.Sprintf("%s.SetUpdatedAt( time.Now().Format(time.RFC3339) )", mg.StructAcronym)
		}
	}

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~updatedAt~"] = updatedAt
	ph["~FullTableName~"] = mg.Model.FullTableName
	ph["~columnList~"] = columnList

	code += util.NamedSprintf(`
// Update a record
func (~StructAcronym~ *~ModelName~) Update() error {
	if !~StructAcronym~.Id.Valid {
		return nil
	}

	~updatedAt~

	_, err = ~StructAcronym~.Ctx.Update("~FullTableName~").
~columnList~
		Where("id = ?", ~StructAcronym~.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetDeleteCode() string {
	var code string
	softDelete := false
	ph := make(map[string]string)

	for key := range mg.Fields {
		if mg.Fields[key].DbFieldName == "deleted_at" {
			softDelete = true
		}
	}

	if softDelete {
		return mg.GetSoftDeleteCode()
	}

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~FullTableName~"] = mg.Model.FullTableName

	code += util.NamedSprintf(`
// Hard delete a record
func (~StructAcronym~ *~ModelName~) Delete() error {

	if !~StructAcronym~.Id.Valid {
		return nil
	}

	_, err = ~StructAcronym~.Ctx.DeleteFrom("~FullTableName~").
		Where("id = ?", ~StructAcronym~.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetSoftDeleteCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~FullTableName~"] = mg.Model.FullTableName

	code += util.NamedSprintf(`
// Soft delete a record
func (~StructAcronym~ *~ModelName~) Delete() error {

	if !~StructAcronym~.Id.Valid {
		return nil
	}

	~StructAcronym~.SetDeletedAt( (time.Now()).Format(time.RFC3339) )

	_, err = ~StructAcronym~.Ctx.Update("~FullTableName~").
		Set("deleted_at", ~StructAcronym~.DeletedAt).
		Where("id = ?", ~StructAcronym~.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetUndeleteCode() string {
	var code string
	softDelete := false
	ph := make(map[string]string)

	for key := range mg.Fields {
		if mg.Fields[key].DbFieldName == "deleted_at" {
			softDelete = true
		}
	}

	if !softDelete {
		return ""
	}

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~FullTableName~"] = mg.Model.FullTableName

	code += util.NamedSprintf(`
// Soft undelete a record
func (~StructAcronym~ *~ModelName~) Undelete() error {

	if !~StructAcronym~.Id.Valid {
		return nil
	}

	~StructAcronym~.SetDeletedAt("")

	_, err = ~StructAcronym~.Ctx.Update("~FullTableName~").
		Set("deleted_at", ~StructAcronym~.DeletedAt).
		Where("id = ?", ~StructAcronym~.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetSetterGetterCode() string {
	var code string

	for key := range mg.Fields {

		code += mg.GetStringGetterCode(mg.Fields[key])
		code += mg.GetStringSetterCode(mg.Fields[key])

		// switch mg.Fields[key].DataType {
		// case "sql.NullString":
		// 	code += mg.GetStringGetterCode(mg.Fields[key])
		// 	code += mg.GetStringSetterCode(mg.Fields[key])
		// case "sql.NullInt64":
		// 	code += mg.GetIntGetterCode(mg.Fields[key])
		// 	code += mg.GetIntSetterCode(mg.Fields[key])
		// case "sql.NullFloat64":
		// 	code += mg.GetFloatGetterCode(mg.Fields[key])
		// 	code += mg.GetFloatSetterCode(mg.Fields[key])
		// case "sql.NullBool":
		// 	code += mg.GetBoolGetterCode(mg.Fields[key])
		// 	code += mg.GetBoolSetterCode(mg.Fields[key])
		// }
	}

	return code
}

//
func (mg *ModelGenerator) GetStringGetterCode(field psql.Field) string {
	var code string
	var slice string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	if field.DbDataType == "date" {
		slice = "[0:10]"
	}

	code += fmt.Sprintf(`
//
func (%s *%s) Get%s() string {

	if %s.Valid {
		return %s.String%s
	}

	return ""
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName, slice)

	return code
}

//
func (mg *ModelGenerator) GetStringSetterCode(field psql.Field) string {
	var code string
	var slice string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	if field.DbDataType == "date" {
		slice = "[0:10]"
	}

	code += fmt.Sprintf(`
//
func (%s *%s) Set%s(val string) {

	if val == "" {
		%s.Valid = false
		%s.String = ""

		return
	}

	%s.Valid = true
	%s.String = val%s
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName, fullFieldName, fullFieldName, slice)

	return code
}

//
func (mg *ModelGenerator) GetIntGetterCode(field psql.Field) string {
	var code string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	code += fmt.Sprintf(`
//
func (%s *%s) Get%s() int64 {

	if %s.Valid {
		return %s.Int64
	}

	return 0
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName)

	return code
}

//
func (mg *ModelGenerator) GetIntSetterCode(field psql.Field) string {
	var code string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	code += fmt.Sprintf(`
//
func (%s *%s) Set%s(val int64) {

	if val == 0 {
		%s.Valid = false
		%s.Int64 = 0

		return
	}

	%s.Valid = true
	%s.Int64 = val
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName, fullFieldName, fullFieldName)

	return code
}

//
func (mg *ModelGenerator) GetFloatGetterCode(field psql.Field) string {
	var code string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	code += fmt.Sprintf(`
//
func (%s *%s) Get%s() float64 {

	if %s.Valid {
		return %s.Float64
	}

	return 0
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName)

	return code
}

//
func (mg *ModelGenerator) GetFloatSetterCode(field psql.Field) string {
	var code string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	code += fmt.Sprintf(`
//
func (%s *%s) Set%s(val float64) {

	if val == 0 {
		%s.Valid = false
		%s.Float64 = 0

		return
	}

	%s.Valid = true
	%s.Float64 = val
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName, fullFieldName, fullFieldName)

	return code
}

//
func (mg *ModelGenerator) GetBoolGetterCode(field psql.Field) string {
	var code string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	code += fmt.Sprintf(`
//
func (%s *%s) Get%s() bool {

	if %s.Valid {
		return %s.Bool
	}

	%s.Valid = true
	%s.Bool = false

	return false
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName, fullFieldName, fullFieldName)

	return code
}

//
func (mg *ModelGenerator) GetBoolSetterCode(field psql.Field) string {
	var code string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	code += fmt.Sprintf(`
//
func (%s *%s) Set%s(val bool) {
	%s.Valid = true
	%s.Bool = val
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName)

	return code
}

//
func (mg *ModelGenerator) IsHiddenField(fieldName string) bool {
	return fieldName == "Id" || fieldName == "AccountId" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "DeletedAt"
}
