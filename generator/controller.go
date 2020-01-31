package generator

import (
	"fmt"
	// "strings"
	"github.com/jschneider98/jgoweb/util"
)

//
func (mg *ModelGenerator) GenerateController() string {
	var code string

	code += mg.GetContollerImportCode()
	code += mg.GetControllerMainCode()
	code += mg.GetControllerReqCode()


	return code
}

//
func (mg *ModelGenerator) GetContollerImportCode() string {
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
func (mg *ModelGenerator) GetControllerMainCode() string {

	tmplParams := struct {
		FuncName string
		ModelName string
	}{}

	tmplParams.FuncName = util.ToLowerCamelCase( util.ToSnakeCase(mg.ModelName) )
	tmplParams.ModelName = mg.ModelName

	str := `
// 
func (ctx *WebContext) {{{.FuncName}}}(rw web.ResponseWriter, req *web.Request) {
	var err error

	params := struct {
		{{{.ModelName}}} *models.{{{.ModelName}}}
		Messages template.HTML
	}{}

	// 
	params.{{{.ModelName}}}, err = ctx.get{{{.ModelName}}}FromRequest(req)

	if err != nil {
		ctx.JobError(util.WhereAmI(), err)
	}

	if req.Method == "POST" {
		msg, isValid, err := params.{{{.ModelName}}}.ProcessSubmit(req)

		if err != nil {
			ctx.JobError(util.WhereAmI(), err)
		}

		if isValid {
			params.Messages = util.GetHtmlAlerts("success", msg)
		} else {
			params.Messages = util.GetHtmlAlerts("danger", msg)
		}
	}

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
func (mg *ModelGenerator) GetControllerReqCode() string {

	tmplParams := struct {
		ModelName string
		Id string
	}{}

	tmplParams.ModelName = mg.ModelName
	tmplParams.Id = util.ToLowerCamelCase( util.ToSnakeCase(mg.ModelName + `Id`) )

	str := `
//
func (ctx *WebContext) get{{{.ModelName}}}FromRequest(req *web.Request) (*models.{{{.ModelName}}}, error) {
	var mdl *models.{{{.ModelName}}}
	var err error

	params := req.URL.Query()
	id := params.Get("{{{.Id}}}")

	if id != "" {
		mdl, err = models.Fetch{{{.ModelName}}}ById(ctx, id)

		if err != nil {
			return nil, err
		}

		if mdl != nil {
			return mdl, nil
		}
	}

	mdl, err = models.New{{{.ModelName}}}(ctx)

	if err != nil {
		return nil, err
	}

	if req.Method == "POST" {
		err = mdl.Hydrate(req)

		if err != nil {
			return nil, err
		}

		return mdl, nil
	}

	mdl.AccountId = ctx.Provider.AccountId

	return mdl, nil
}
`

	code, err := util.TemplateToString(str, tmplParams)

	if err != nil {
		code = fmt.Sprintf("\n**** ERROR ****\n%s\n********\n", err)
	}

	return code
}
