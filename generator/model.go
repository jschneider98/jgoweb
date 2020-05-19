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

	_, err := ~StructAcronym~.Ctx.Update("~FullTableName~").
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
func (mg *ModelGenerator) IsSoftDelete() bool {

	for key := range mg.Fields {
		if mg.Fields[key].DbFieldName == "deleted_at" {
			return true
		}
	}

	return false
}

//
func (mg *ModelGenerator) GetDeleteCode() string {
	var code string
	ph := make(map[string]string)

	if mg.IsSoftDelete() {
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

	_, err := ~StructAcronym~.Ctx.DeleteFrom("~FullTableName~").
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

	_, err := ~StructAcronym~.Ctx.Update("~FullTableName~").
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

	_, err := ~StructAcronym~.Ctx.Update("~FullTableName~").
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

// **************

//
func (mg *ModelGenerator) GenerateTest() string {
	var code string

	code = mg.GetTestImportCode()
	code += mg.GetTestFetchByIdCode()
	code += mg.GetTestSetterGetterCode()
	code += mg.GetTestInsertCode()
	code += mg.GetTestUpdateCode()
	code += mg.GetTestDeleteCode()

	if mg.IsSoftDelete() {
		code += mg.GetTestUndeleteCode()
	}

	code += mg.GetTestNewWithDataCode()
	code += mg.GetTestProcessSubmitCode()

	return code
}

//
func (mg *ModelGenerator) GetTestImportCode() string {
	return `// +build integration

package models

import (
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb"
	"net/http"
	"testing"
	"strings"
)
`
}

//
func (mg *ModelGenerator) GetTestFetchByIdCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName

	code += util.NamedSprintf(`
//
func TestFetch~ModelName~ById(t *testing.T) {
	jgoweb.InitMockCtx()
	InitMock~ModelName~()

	// force not found
	id := "00000000-0000-0000-0000-000000000000"
	~StructAcronym~, err := Fetch~ModelName~ById(jgoweb.MockCtx, id)

	if err != nil {
		t.Errorf("\nERROR: Failed to fetch ~ModelName~ by id: %%v\n", err)
		return
	}

	if ~StructAcronym~ != nil {
		t.Errorf("\nERROR: Should have failed to find ~ModelName~: %%v\n", id)
		return
	}

	~StructAcronym~, err = Fetch~ModelName~ById(jgoweb.MockCtx, Mock~ModelName~.GetId())

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
		return
	}

	if ~StructAcronym~ == nil {
		t.Errorf("\nERROR: Should have found ~ModelName~ with Id: %%v\n", Mock~ModelName~.GetId())
		return
	}

	if ~StructAcronym~.GetId() != Mock~ModelName~.GetId() {
		t.Errorf("\nERROR: Fetch mismatch. Expected: %%v Got: %%v\n", Mock~ModelName~.GetId(), ~StructAcronym~.GetId())
		return
	}
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetTestSetterGetterCode() string {
	code := ""

	for _, field := range mg.Fields {
		code += mg.GetTestSetterGetterByField(field.FieldName)
	}

	return code
}

//
func (mg *ModelGenerator) GetTestSetterGetterByField(fieldName string) string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~FieldName~"] = fieldName

	code += util.NamedSprintf(`
//
func Test~ModelName~~FieldName~(t *testing.T) {
	InitMock~ModelName~()
	origVal := Mock~ModelName~.Get~FieldName~()
	testVal := "test"

	Mock~ModelName~.Set~FieldName~("")

	if Mock~ModelName~.~FieldName~.Valid {
		t.Errorf("ERROR: ~FieldName~ should be invalid.\n")
	}

	if Mock~ModelName~.Get~FieldName~() != "" {
		t.Errorf("ERROR: Set ~FieldName~ failed. Should have a blank value. Got: %%s", Mock~ModelName~.Get~FieldName~())
	}

	Mock~ModelName~.Set~FieldName~(testVal)

	if !Mock~ModelName~.~FieldName~.Valid {
		t.Errorf("ERROR: ~FieldName~ should be valid.\n")
	}

	if Mock~ModelName~.Get~FieldName~() != testVal {
		t.Errorf("ERROR: Set ~FieldName~ failed. Expected: %%s, Got: %%s", testVal, Mock~ModelName~.Get~FieldName~())
	}

	Mock~ModelName~.Set~FieldName~(origVal)
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) IsTimestamp(fieldName string) bool {
	return fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "DeletedAt"
}

//
func (mg *ModelGenerator) GetTestInsertCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~SetSaveVals~"] = ""
	ph["~Setters~"] = ""
	ph["~Assert~"] = fmt.Sprintf("%s == nil", mg.StructAcronym)

	for _, field := range mg.Fields {

		if !mg.IsTimestamp(field.FieldName) && field.FieldName != "Id" {
			ph["~SetSaveVals~"] += fmt.Sprintf("%s := \"%s Insert\"\n\t", field.FieldName, field.FieldName)
			ph["~Setters~"] += fmt.Sprintf("%s.Set%s(%s)\n\t", mg.StructAcronym, field.FieldName, field.FieldName)
			ph["~Assert~"] += fmt.Sprintf(" || %s.Get%s() != %s", mg.StructAcronym, field.FieldName, field.FieldName)
		}
	}

	code += util.NamedSprintf(`
//
func Test~ModelName~Insert(t *testing.T) {
	InitMock~ModelName~()
	~SetSaveVals~
	~StructAcronym~, err := New~ModelName~(jgoweb.MockCtx)

	if err != nil {
		t.Errorf("\nERROR: New~ModelName~() failed. %%v\n", err)
	}

	~Setters~
	err = ~StructAcronym~.Save()

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	if !~StructAcronym~.Id.Valid {
		t.Errorf("\nERROR: ~ModelName~.Id should be set.\n")
	}

	// verify write
	~StructAcronym~, err = Fetch~ModelName~ById(jgoweb.MockCtx, ~StructAcronym~.GetId())

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	if ~Assert~ {
		t.Errorf("\nERROR: ~ModelName~ does not match save values. Insert failed.\n")
	}
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetTestUpdateCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~SetSaveVals~"] = ""
	ph["~Setters~"] = ""
	ph["~Assert~"] = fmt.Sprintf("%s == nil", mg.StructAcronym)

	for _, field := range mg.Fields {

		if !mg.IsTimestamp(field.FieldName) && field.FieldName != "Id" {
			ph["~SetSaveVals~"] += fmt.Sprintf("%s := \"%s Update\"\n\t", field.FieldName, field.FieldName)
			ph["~Setters~"] += fmt.Sprintf("Mock%s.Set%s(%s)\n\t", mg.ModelName, field.FieldName, field.FieldName)
			ph["~Assert~"] += fmt.Sprintf(" || %s.Get%s() != %s", mg.StructAcronym, field.FieldName, field.FieldName)
		}
	}

	code += util.NamedSprintf(`
//
func Test~ModelName~Update(t *testing.T) {
	InitMock~ModelName~()
	~SetSaveVals~
	~Setters~
	err := Mock~ModelName~.Save()

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	// verify write
	~StructAcronym~, err := Fetch~ModelName~ById(jgoweb.MockCtx, Mock~ModelName~.GetId())

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	if ~Assert~ {
		t.Errorf("\nERROR: ~ModelName~ does not match save values. Update failed.\n")
	}
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetTestDeleteCode() string {

	if mg.IsSoftDelete() {
		return mg.GetTestSoftDeleteCode()
	}

	return mg.GetTestHardDeleteCode()
}

//
func (mg *ModelGenerator) GetTestSoftDeleteCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName

	code += util.NamedSprintf(`
//
func Test~ModelName~Delete(t *testing.T) {
	InitMock~ModelName~()
	err := Mock~ModelName~.Delete()

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	// verify write
	~StructAcronym~, err := Fetch~ModelName~ById(jgoweb.MockCtx, Mock~ModelName~.GetId())

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	if !~StructAcronym~.DeletedAt.Valid {
		t.Errorf("\nERROR: ~ModelName~ does not match save values. Delete failed.\n")
	}
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetTestHardDeleteCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName

	code += util.NamedSprintf(`
//
func Test~ModelName~Delete(t *testing.T) {
	InitMock~ModelName~()
	err := Mock~ModelName~.Delete()

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	// verify write
	~StructAcronym~, err := Fetch~ModelName~ById(jgoweb.MockCtx, Mock~ModelName~.GetId())

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	if ~StructAcronym~ != nil {
		t.Errorf("\nERROR: Delete failed. Fetch should return nil.\n")
		return
	}

	Mock~ModelName~ = nil
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetTestUndeleteCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName

	code += util.NamedSprintf(`
//
func Test~ModelName~Undelete(t *testing.T) {
	InitMock~ModelName~()
	err := Mock~ModelName~.Delete()

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	err = Mock~ModelName~.Undelete()

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	// verify write
	~StructAcronym~, err := Fetch~ModelName~ById(jgoweb.MockCtx, Mock~ModelName~.GetId())

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	if ~StructAcronym~ == nil || ~StructAcronym~.DeletedAt.Valid {
		t.Errorf("\nERROR: ~ModelName~ does not match save values. Undelete failed.\n")
	}
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetTestNewWithDataCode() string {
	var code string
	ph := make(map[string]string)

	ph["~ModelName~"] = mg.ModelName

	code += util.NamedSprintf(`
//
func TestNew~ModelName~WithData(t *testing.T) {
	httpReq, err := http.NewRequest("GET", "http://example.com", nil)

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}

	req := &web.Request{}
	req.Request = httpReq

	_, err = New~ModelName~WithData(jgoweb.MockCtx, req)

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
	}
}
`, ph)

	return code
}

//
func (mg *ModelGenerator) GetTestProcessSubmitCode() string {
	var code string
	ph := make(map[string]string)

	ph["~StructAcronym~"] = mg.StructAcronym
	ph["~ModelName~"] = mg.ModelName
	ph["~PostVals~"] = "z=post"

	for _, field := range mg.Fields {

		if !mg.IsTimestamp(field.FieldName) && field.FieldName != "Id" {
			ph["~PostVals~"] += fmt.Sprintf("&%s=%s", field.FieldName, field.FieldName)
		}
	}

	code += util.NamedSprintf(`
//
func Test~ModelName~ProcessSubmit(t *testing.T) {
	~StructAcronym~, err := New~ModelName~(jgoweb.MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
		return
	}

	httpReq, err := http.NewRequest("POST", "http://example.com", strings.NewReader("~PostVals~"))

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req := &web.Request{}
	req.Request = httpReq

	msg, saved, err := ~StructAcronym~.ProcessSubmit(req)

	if err != nil {
		t.Errorf("\nERROR: %%v\n", err)
		return
	}

	if !saved {
		t.Errorf("\nERROR: %%v", msg)
	}
}
`, ph)

	return code
}
