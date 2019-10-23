package jgoweb

import (
	"fmt"
	"strings"
	"net/url"
	"github.com/gocraft/dbr"
)

type SearchParams struct {
	Query string
	Limit uint64
	Offset uint64
	TableAlias string
	FirstName string
	LastName string
	IdField string
	UrlParams url.Values
}

//
func NewSearchParams() *SearchParams {
	sp := &SearchParams{}

	sp.TableAlias = "main"
	sp.FirstName = "first_name"
	sp.LastName = "last_name"
	sp.IdField = "school_state_id"

	return sp
}

//
func (sp *SearchParams) BuildDefaultCondition() (dbr.Builder, error) {
	var builder dbr.Builder

	if sp.Query == "" {
		return nil, nil
	}

	queryParts := strings.Split(sp.Query, " ")

	if len(queryParts) == 1 {
		builder = dbr.And(
			dbr.Or(
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", sp.TableAlias, sp.FirstName), "%" + sp.Query + "%" ),
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", sp.TableAlias, sp.LastName), "%" + sp.Query + "%" ),
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", sp.TableAlias, sp.IdField), "%" + sp.Query + "%" ),
			),
		)

		return builder, nil
	}

	builder = dbr.And(
		dbr.Or(
			dbr.And(
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", sp.TableAlias, sp.FirstName), "%" + queryParts[0] + "%" ),
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", sp.TableAlias, sp.LastName), "%" + queryParts[1] + "%" ),
			),
			dbr.And(
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", sp.TableAlias, sp.FirstName), "%" + queryParts[1] + "%" ),
				dbr.Expr( fmt.Sprintf("%s.%s ilike ?", sp.TableAlias, sp.LastName), "%" + queryParts[0] + "%" ),
			),
			dbr.Expr( fmt.Sprintf("%s.%s ilike ?", sp.TableAlias, sp.IdField), "%" + sp.Query + "%" ),
		),
	)

	return builder, nil
}
