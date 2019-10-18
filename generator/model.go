package generator

import (
	"fmt"
	"strings"
	"github.com/jschneider98/jgoweb"
	"github.com/jschneider98/jgoweb/util"
	"github.com/jschneider98/jgoweb/db/psql"
	"github.com/jschneider98/jgomodel"
)

type ModelGenerator struct {
	ModelName string `json:"model_name"`
	InstanceName string `json:"instance_name"`
	StructAcronym string
	Model *jgomodel.Model
	Fields []psql.Field
}

//
func NewModelGenerator(ctx jgoweb.ContextInterface, schema string, table string) (*ModelGenerator, error) {
	var err error
	
	mg := &ModelGenerator{}
	mg.Model, err = jgomodel.NewModel(ctx, schema, table)

	if err != nil {
		return nil, err
	}

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

	// Unpluralize table name...
	table := strings.TrimSuffix(mg.Model.Table, "s")

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
	code += mg.GetHydrateCode()
	code += mg.GetIsValidCode()
	code += mg.GetSaveCode()
	code += mg.GetInsertCode()
	code += mg.GetUpdateCode()
	code += mg.GetDeleteCode()
	code += mg.GetSetterGetterCode()

	return code
}

// 
func (mg *ModelGenerator) GetImportCode() string {
return	`package models

import(
	"database/sql"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb"
	"github.com/jschneider98/jgomodel"
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
	code += "\t*jgomodel.Model `json:\"-\" validate:\"-\"`\n"
	code += "}\n\n"

	return code
}


//
func (mg *ModelGenerator) GetHydratorStructCode() string {
	var code string

	code += fmt.Sprintf("// %s\n", mg.ModelName + "Hydrator")
	code += fmt.Sprintf("type %s struct {\n", mg.ModelName + "Hydrator")

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
`, acronym, mg.ModelName + "Hydrator", acronym, acronym)

	return code
}

//
func (mg *ModelGenerator) GetNewCode() string {
	var code string

	code += fmt.Sprintf(`
// Empty new model
func New%s(ctx jgoweb.ContextInterface) (*%s, error) {
	%s, err := jgomodel.NewModel(ctx, "%s", "%s")

	if err != nil {
		return nil, err
	}

	%s := &%s{Model:%s}
	%s.Ctx = ctx
	%s.SetDefaults()

	return %s, nil
}
`, mg.ModelName, mg.ModelName, mg.InstanceName, mg.Model.Schema, mg.Model.Table, mg.StructAcronym, mg.ModelName, mg.InstanceName, mg.StructAcronym, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetSetDefaultCode() string {
	var code string
	var defaults string

	for key := range mg.Fields {

		if mg.Fields[key].DbDefault.Valid && !mg.Fields[key].DbDefaultIsFunc {
			defaults += fmt.Sprintf("\t%s.Set%s(%s)\n", mg.StructAcronym, mg.Fields[key].FieldName, mg.Fields[key].Default)
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
`, mg.ModelName, mg.ModelName, mg.InstanceName, mg.ModelName, mg.InstanceName, mg.InstanceName)

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
	%s[0].Model, err = jgomodel.NewModel(ctx, "%s", "%s")

	if err != nil {
		return nil, err
	}

	return &%s[0], nil
}
`, mg.ModelName, mg.ModelName, mg.InstanceName, mg.ModelName, mg.Model.FullTableName, mg.InstanceName, mg.InstanceName, mg.InstanceName, mg.InstanceName, mg.Model.Schema, mg.Model.Table, mg.InstanceName)

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

	code += fmt.Sprintf(`
// Validate the model
func (%s *%s) IsValid() error {
	return %s.Ctx.GetValidator().Struct(%s)
}
`, mg.StructAcronym, mg.ModelName, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetSaveCode() string {
	var code string

	code += fmt.Sprintf(`
// Insert/Update based on pkey value
func (%s *%s) Save() error {
	err := %s.IsValid()

	if err != nil {
		return err
	}

	if !%s.Id.Valid {
		return %s.Insert()
	} else {
		return %s.Update()
	}
}
`, mg.StructAcronym, mg.ModelName, mg.StructAcronym, mg.StructAcronym, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetInsertCode() string {
	var code string
	var objCols []string
	var colList string

	// 
	for key := range mg.Fields {
		if (mg.Fields[key].DbFieldName != "id" && mg.Fields[key].DbFieldName != "created_at" && mg.Fields[key].DbFieldName != "updated_at") {
			// (p.AccountId, p.Units, ...)
			objCols = append(objCols, fmt.Sprintf("%s.%s", mg.StructAcronym, mg.Fields[key].FieldName))
		}
	}

	colList = strings.Join(objCols, ",\t\t\t")
	colList = strings.ReplaceAll(colList, ",", ",\n")

	code += fmt.Sprintf(`
// Insert a new record
func (%s *%s) Insert() error {
	tx, err := %s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	query := %s.GetInsertQuery()

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(%s).Scan(&%s.Id)

	if err != nil {
		return err
	}

	return %s.Ctx.OptionalCommit(tx)
}
`, mg.StructAcronym, mg.ModelName, mg.StructAcronym, mg.StructAcronym, colList, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetUpdateCode() string {
	var code string
	columnList := ""

	for key := range mg.Fields {
		columnList += fmt.Sprintf("\t\tSet(\"%s\", %s.%s).\n", mg.Fields[key].DbFieldName, mg.StructAcronym, mg.Fields[key].FieldName)
	}

	code += fmt.Sprintf(`
// Update a record
func (%s *%s) Update() error {
	if !%s.Id.Valid {
		return nil
	}

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
`, mg.StructAcronym, mg.ModelName, mg.StructAcronym, mg.StructAcronym, mg.Model.FullTableName, columnList, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetDeleteCode() string {
	var code string
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

	if !%s.Id.Valid {
		return nil
	}

	tx, err := %s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom("%s").
		Where("id = ?", %s.Id).
		Exec()

	if err != nil {
		return err
	}

	return %s.Ctx.OptionalCommit(tx)
}
`, mg.StructAcronym, mg.ModelName, mg.StructAcronym, mg.StructAcronym, mg.Model.FullTableName, mg.StructAcronym, mg.StructAcronym)

	return code
}

//
func (mg *ModelGenerator) GetSoftDeleteCode() string {
	var code string

	code += fmt.Sprintf(`
// Soft delete a record
func (%s *%s) Delete() error {
	tx, err := %s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.Update("%s").
		Set("deleted_at = ?", "timezone('utc'::text, now())").
		Where("id = ?", %s.Id).
		Exec()

	if err != nil {
		return err
	}

	return p.Ctx.OptionalCommit(tx)
}
`, mg.StructAcronym, mg.ModelName, mg.StructAcronym, mg.Model.FullTableName, mg.StructAcronym)

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
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	code += fmt.Sprintf(`
//
func (%s *%s) Get%s() string {

	if %s.Valid {
		return %s.String
	}

	return ""
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName)

	return code
}

//
func (mg *ModelGenerator) GetStringSetterCode(field psql.Field) string {
	var code string
	fullFieldName := fmt.Sprintf("%s.%s", mg.StructAcronym, field.FieldName)

	code += fmt.Sprintf(`
//
func (%s *%s) Set%s(val string) {

	if val == "" {
		%s.Valid = false
		%s.String = ""

		return
	}

	%s.Valid = true
	%s.String = val
}
`, mg.StructAcronym, mg.ModelName, field.FieldName, fullFieldName, fullFieldName, fullFieldName, fullFieldName)

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
func  (mg *ModelGenerator) IsHiddenField(fieldName string) bool {
	return fieldName == "Id" || fieldName == "AccountId" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "DeletedAt"
}
