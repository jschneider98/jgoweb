package jgoweb

import (
	"fmt"
	"strings"
	"github.com/gocraft/dbr"
)

type SearchParams struct {
	Query string
	Limit uint64
}

//
func NewSearchParams() *SearchParams {
	return &SearchParams{}
}

//
func (sp *SearchParams) BuildNameCondition(alias string, firstName string, lastName string) (dbr.Builder, error) {
	var builder dbr.Builder

	if sp.Query == "" {
		return nil, nil
	}

	queryParts := strings.Split(sp.Query, " ")

	if len(queryParts) == 1 {
		builder = dbr.And(
			dbr.Or(
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", alias, firstName), sp.Query + "%" ),
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", alias, lastName), sp.Query + "%" ),
			),
		)

		return builder, nil
	}

	builder = dbr.Or(
		dbr.And(
			dbr.Expr( fmt.Sprintf("%s.%s ilike ?", alias, firstName), queryParts[0] + "%" ),
			dbr.Expr( fmt.Sprintf("%s.%s ilike ?", alias, lastName), queryParts[1] + "%" ),
		),
		dbr.And(
			dbr.Expr( fmt.Sprintf("%s.%s ilike ?", alias, firstName), queryParts[1] + "%" ),
			dbr.Expr( fmt.Sprintf("%s.%s ilike ?", alias, lastName), queryParts[0] + "%" ),
		),
	)

	return builder, nil
}
