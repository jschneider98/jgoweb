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
	code += mg.GetViewBodyCode()
	code += mg.GetViewScriptCode()

	return code
}

//
func (mg *ModelGenerator) GetViewTitleCode() string {
	var code string

	code += fmt.Sprintf("[[define \"title\"]]%s[[end]]\n", mg.ModelName)

	return code
}

//
func (mg *ModelGenerator) GetViewBodyCode() string {
	var code string
	var inputs string


	for key := range mg.Fields {
		fieldName := mg.Fields[key].FieldName

		if mg.IsHiddenField(fieldName) {
			inputs += fmt.Sprintf(`
			<input type="hidden" name="%s" v-model="%s"/>`, fieldName, fieldName)
		} else {
			inputs += fmt.Sprintf(`
			<div class="col-sm-3 my-1">
				<div class="form-group">
					<label for="%s">%s</label>
					<input type="text" class="form-control" id="%s" name="%s" aria-describedby="%sHelp" placeholder="Enter %s" v-model="%s">
				</div>
			</div>
`, fieldName, util.ToWords(fieldName), fieldName, fieldName, fieldName, util.ToWords(fieldName), fieldName)
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