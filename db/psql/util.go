package psql

import (
	"fmt"
	"regexp"
	"strings"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
	"github.com/jschneider98/jgoweb"
	"github.com/jschneider98/jgoweb/util"
)

type Field struct {
	DbFieldName string `json:"db_field_name"`
	DbDataType string `json:"db_data_type"`
	DbDefault dbr.NullString `json:"db_default"`
	DbDefaultIsFunc bool `json:"db_default_is_func`
	FieldName string `json:"field_name"`
	DataType string `json:"data_type"`
	Default string `json:"default"`
	Annotation string `json:"annotation"`
	NotNull bool `json:"not_null"`
	SortNum int `json:"sort_num"`
}

func GetFields(ctx jgoweb.ContextInterface, schema string, table string) ([]Field, error) {
	var fields []Field

	stmt := ctx.Select(`
		a.attname as db_field_name,
		pg_catalog.format_type(a.atttypid, a.atttypmod) as db_data_type,
		(
			SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128)
			FROM pg_catalog.pg_attrdef d
			WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef
		) as db_default,
		a.attnotnull as not_null,
		a.attnum as sort_num
	`).
	From(dbr.I("pg_catalog.pg_class").As("c")).
	Join(dbr.I("pg_catalog.pg_namespace").As("n"), "n.oid = c.relnamespace").
	Join(dbr.I("pg_catalog.pg_attribute").As("a"), "c.oid = a.attrelid").
	Where("n.nspname = ?", schema).
	Where("c.relname = ?", table).
	Where("a.attnum > ?", 0).
	Where("NOT a.attisdropped").
	OrderBy("a.attnum")

	_, err := stmt.Load(&fields)

	if err != nil {
		return nil, err
	}

	for key := range fields {
		fields[key].SetStructVals()
	}

	return fields, nil
}

//
func (f *Field) SetStructVals() {
	f.SetFieldName()
	f.SetDefault()
	f.SetDataType()
	f.SetAnnotation()
}

//
func (f *Field) SetFieldName() {
	f.FieldName = util.ToCamelCase(f.DbFieldName)
}

//
func (f *Field) SetDefault() {

	if !f.DbDefault.Valid {
		return
	}

	if strings.Contains(f.DbDefault.String, "(") {
		f.DbDefaultIsFunc = true
	} else {
		f.Default = f.DbDefault.String
	}
}

//
func (f *Field) SetDataType() {

	// switch f.DbDataType {
	// case "smallint":
	// case "smallserial":
	// case "serial":
	// case "bigint":
	// case "bigserial":
	// case "integer":
	// 	f.DataType = "sql.NullInt64"
	// 	return
	// case "boolean":
	// 	f.DataType = "sql.NullBool"
	// 	return
	// case "double precision":
	// 	f.DataType = "sql.NullFloat64"
	// 	return
	// }

	f.DataType = "sql.NullString"
}

//
func (f *Field) SetAnnotation() {
	f.Annotation = fmt.Sprintf("`json:\"%s\" %s`", f.FieldName, f.GetValidation())
}

//
func (f *Field) GetValidation() string {
	val := "validate:"

	// if not null and no default = required (Special case, bool = notnull)
	// if not null with default = no insert/update (i.e., use default)
	// if nullable = omitempty
	if f.NotNull == true && !f.DbDefault.Valid {

		if f.DbDataType == "boolean" {
			val += `"notNull`
		} else {
			val += `"required`
		}
	} else {
		val += `"omitempty`
	}

	// if f.NotNull == true && !f.DbDefault.Valid {
	// 	val += `"required`
	// } else {
	// 	val += `"omitempty`
	// }	

	switch f.DbDataType {
	case "smallint", "smallserial", "serial", "integer", "bigint", "bigserial":
		val += ",int"
	case "double precision", "real":
		val += ",numeric"
	case "uuid":
		val += ",uuid"
	case "timestamp with time zone":
		val += ",rfc3339"
	case "timestamp without time zone":
		val += ",rfc3339WithoutZone"
	case "date":
		val += ",date"
	}

	if strings.HasPrefix(f.DbDataType, "character varying(") {
		re := regexp.MustCompile("[0-9]+")
		val += ",min=1,max=" + re.FindString(f.DbDataType)
	}

	if strings.HasPrefix(f.DbDataType, "character(") {
		re := regexp.MustCompile("[0-9]+")
		val += ",min=1,max=" + re.FindString(f.DbDataType)
	}

	val += `"`

	return val
}


//
func DebugSqlStatement(stmt *dbr.SelectStmt) error {
	var err error

	buf := dbr.NewBuffer()
	err = stmt.Build(dialect.PostgreSQL, buf)

	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", buf.String())

	for key := range buf.Value() {
		fmt.Printf("%v == %v\n", key, buf.Value()[key])
	}

	return nil
}
