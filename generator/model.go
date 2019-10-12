package generator

import (
	"github.com/jschneider98/jgoweb"
	"github.com/jschneider98/db/psql"
)

type ModelGenerator struct {
	Schema string `json:"schema"`
	Table string `json:"table"`
	Ctx jgoweb.ContextInterface `json:"-"`
}

//
func NewModelGenerator(ctx jgoweb.ContextInterface, schema string, table string) *ModelGenerator {
	return &ModelGenerator{
		Schema: schema,
		Table: table,
		Ctx: ctx
	}
}

func (mg *ModelGenerator) SetFromSession() error {

}
