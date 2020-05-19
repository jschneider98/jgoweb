package jgoweb

import (
	"fmt"
	"strings"
)

type QueryIntBuilder struct {
	Query        string
	RawQuery     string
	InvalidQuery string
	ctx          ContextInterface
}

//
func NewQueryIntBuilder(ctx ContextInterface, rawQuery string) *QueryIntBuilder {
	builder := QueryIntBuilder{}

	builder.RawQuery = rawQuery
	builder.ctx = ctx

	return &builder
}

// @TODO: auto fix invalid query strings? Improve logic like crazy? Beter validation?
func (b *QueryIntBuilder) Build() (string, error) {
	var queryPart string
	var err error
	validQuery := ""

	parts := strings.Fields(b.Parse())
	addAnd := false

	for _, part := range parts {

		if b.IsValid() {
			validQuery = b.Query
		}

		opertor := b.GetOperator(part)

		if addAnd && (opertor == "" || opertor == "(" || opertor == "!") {
			b.Query += "&"
		}

		if len(opertor) > 0 {
			b.Query += opertor
			addAnd = false
			continue
		}

		queryPart, err = b.GetQuery(part)

		if err != nil {
			return "", err
		}

		if len(queryPart) > 0 {
			b.Query += queryPart
			addAnd = true
		}
	}

	if !b.IsValid() {
		b.InvalidQuery = b.Query
		// set to last known valid query
		b.Query = validQuery
	}

	return b.Query, nil
}

//
func (b *QueryIntBuilder) IsValid() bool {
	var temp string

	if b.Query == "" {
		return false
	}

	err := b.ctx.SelectBySql("SELECT (?::query_int)::text", b.Query).LoadOne(&temp)

	if err != nil {
		fmt.Printf("\n%v\n", err)
		return false
	}

	return true
}

//
func (b *QueryIntBuilder) Parse() string {
	query := strings.ToLower(b.RawQuery)
	query = strings.Replace(query, "(", " ( ", -1)
	query = strings.Replace(query, ")", " ) ", -1)

	return query
}

//
func (b *QueryIntBuilder) GetOperator(rawQuery string) string {

	switch rawQuery {
	case "and":
		return "&"
	case "or":
		return "|"
	case "not":
		return "!"
	case "(":
		return "("
	case ")":
		return ")"
	}

	return ""
}

//
func (b *QueryIntBuilder) GetQuery(rawQuery string) (string, error) {
	var result string

	sql := `
SELECT
	CASE WHEN query IS NOT NULL THEN
		format('(%s)', query)
	ELSE
		'(0)'
	END as query
FROM (
	SELECT
		array_to_string(array_agg(id), '|')  as query
	FROM tags
	WHERE name like ?
) as main
`

	stmt := b.ctx.SelectBySql(sql, "%"+rawQuery+"%")
	err := stmt.LoadOne(&result)

	if err != nil {
		return "", err
	}

	return result, nil
}
