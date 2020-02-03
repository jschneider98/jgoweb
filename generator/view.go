package generator

import (
	"fmt"
	"strings"
	"github.com/jschneider98/jgoweb/util"
)

//
func (mg *ModelGenerator) GenerateView() string {
	var code string

	code += mg.GetViewTitleCode()
	code += mg.GetViewSubNavCode()
	code += mg.GetViewBodyCode()
	code += mg.GetViewScriptCode()

	return code
}

//
func (mg *ModelGenerator) GetViewTitleCode() string {
	var code string

	code += fmt.Sprintf("[[define \"title\"]]%s[[end]]\n", util.ToWords(mg.ModelName))

	return code
}

//
func (mg *ModelGenerator) GetViewSubNavCode() string {

	tmplParams := struct {
		ModelName string
	}{}

	tmplParams.ModelName = util.ToSnakeCase(mg.ModelName)

	str := `
[[define "sub-nav"]]
<div class="p-2 text-muted" style="border-top: 1px solid #ddd">
	<div class="col-auto">
		<a class="nav-link" href="/{{{.ModelName}}}_list">List</a>
	</div>
</div>
[[end]]
`

	code, err := util.TemplateToString(str, tmplParams)

	if err != nil {
		code = fmt.Sprintf("\n**** ERROR ****\n%s\n********\n", err)
	}

	return code
}

//
func (mg *ModelGenerator) GetViewBodyCode() string {
	var code string
	var inputs string
	var required string
	var labelRequired string
	var fieldName string


	for _, field := range mg.Fields {
		fieldName = field.FieldName
		required = ""
		labelRequired = ""

		if mg.IsHiddenField(fieldName) {
			inputs += fmt.Sprintf(`
			<input type="hidden" name="%s" v-model="%s"/>`, fieldName, fieldName)
		} else {

			if field.NotNull {
				required = "required"
				labelRequired = "*"
			}

			inputs += fmt.Sprintf(`
			<div class="col-sm-3 my-1">
				<div class="form-group">
					<label for="%s">%s%s</label>
					<input type="text" class="form-control" id="%s" name="%s" aria-describedby="%sHelp" placeholder="Enter %s" v-model="%s" %s>
				</div>
			</div>
`, fieldName, util.ToWords(fieldName), labelRequired, fieldName, fieldName, fieldName, util.ToWords(fieldName), fieldName, required)
		}
	}

	code += fmt.Sprintf(`
[[define "body"]]
<div>
	<form id="%sForm" action="/%v" method="POST">
		<div class="form-row">
%s
		</div>
		<button type="submit" class="btn btn-primary">Submit</button>
	</form>
</div>
[[end]]
`, mg.InstanceName, util.ToSnakeCase(mg.ModelName), inputs)

	return code
}

//
func (mg *ModelGenerator) GetViewScriptCode() string {
	var code string
	var fields []string

	for key := range mg.Fields {
		fieldName := mg.Fields[key].FieldName
		fields = append(fields, fmt.Sprintf("\t\t\t\t%s: [[.%s.Get%s]]", fieldName, mg.ModelName, fieldName) )
	}

	fieldCode := strings.Join(fields, ",\n")

	code += fmt.Sprintf(`
[[define "scripts"]]
<script>
	var app = new Vue({
		el: '#app',
		data: {
%s
		}
	});
</script>
[[end]]
`, fieldCode)

	return code
}
