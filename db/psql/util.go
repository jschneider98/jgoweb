package psql

import (
	"github.com/gocraft/dbr"
	"github.com/jschneider98/jgoweb"
)

type Field struct {
	FieldName string `json:"field_name"`
	DataType string `json:"data_type"`
	Default dbr.NullString `json:"default"`
	NotNull bool `json:"not_null"`
	Attnum int `json:"attnum"`
}

func GetFields(ctx jgoweb.ContextInterface, schema string, table string) ([]Field, error) {
	var fields []Field

	stmt := ctx.Select(`
		a.attname as field_name,
		pg_catalog.format_type(a.atttypid, a.atttypmod) as data_type,
		(
			SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128)
			FROM pg_catalog.pg_attrdef d
			WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef
		) as default,
		a.attnotnull as not_null,
		a.attnum
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

	return fields, nil
}
