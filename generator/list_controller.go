package generator

import (
	"fmt"
	// "strings"
	"github.com/jschneider98/jgoweb/util"
)

//
func (mg *ModelGenerator) GenerateListController() string {
	var code string

	code += mg.GetListContollerImportCode()
	code += mg.GetListControllerMainCode()
	code += mg.GetListControllerDeleteCode()

	return code
}

//
func (mg *ModelGenerator) GetListContollerImportCode() string {
	code := `
package web

import (
	"html/template"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
	"github.com/jschneider98/medex/models"
)
`
	return code
}

//
func (mg *ModelGenerator) GetListControllerMainCode() string {

	tmplParams := struct {
		Prefix string
		ModelName string
	}{}

	tmplParams.Prefix = util.ToLowerCamelCase( util.ToSnakeCase(mg.ModelName) )
	tmplParams.ModelName = mg.ModelName

	str := `
// route
func (ctx *WebContext) {{{.Prefix}}}List(rw web.ResponseWriter, req *web.Request) {
	var err error

	params := struct {
		Results []models.{{{.ModelName}}}
		Messages template.HTML
	}{}

	params.Messages, err = ctx.{{{.Prefix}}}Delete(req)

	if err != nil {
		ctx.JobError(util.WhereAmI(), err)
	}

	results, err := models.FetchAll{{{.ModelName}}}ByAccountId(ctx, ctx.MedexUser.GetAccountId())

	if err != nil {
		ctx.JobError(util.WhereAmI(), err)
	}

	if results == nil {
		results = make([]models.{{{.ModelName}}}, 0)
	}

	params.Results = results

	err = ctx.Template.Execute(rw, params)
	
	if err != nil {
		ctx.JobError(util.WhereAmI(), err)
	}

	ctx.JobSuccess()
}
`
	code, err := util.TemplateToString(str, tmplParams)

	if err != nil {
		code = fmt.Sprintf("\n**** ERROR ****\n%s\n********\n", err)
	}

	return code
}

//
func (mg *ModelGenerator) GetListControllerDeleteCode() string {

	tmplParams := struct {
		Prefix string
		ModelName string
	}{}

	tmplParams.Prefix = util.ToLowerCamelCase( util.ToSnakeCase(mg.ModelName) )
	tmplParams.ModelName = mg.ModelName

	str := `
//
func (ctx *WebContext) {{{.Prefix}}}Delete(req *web.Request) (template.HTML, error) {
	params := req.URL.Query()
	action := params.Get("action")
	{{{.Prefix}}}Id := params.Get("{{{.Prefix}}}Id")

	if action != "delete" || {{{.Prefix}}}Id == "" {
		return template.HTML(""), nil
	}

	mdl, err := models.Fetch{{{.ModelName}}}ById(ctx, {{{.Prefix}}}Id)

	if err != nil {
		return template.HTML(""), err
	}

	err = mdl.Delete()

	if err != nil {
		return template.HTML(""), err
	}

	return util.GetHtmlAlerts("success", "Data deleted."), nil
}
`
	code, err := util.TemplateToString(str, tmplParams)

	if err != nil {
		code = fmt.Sprintf("\n**** ERROR ****\n%s\n********\n", err)
	}

	return code
}
