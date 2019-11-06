package generator

import (
	"fmt"
	"strings"
	"github.com/jschneider98/jgoweb/util"
)

type ReportViewGenerator struct {
	BaseStructName string `json:"model_name"`
	InstanceName string `json:"instance_name"`
	StructAcronym string
	Fields []string
}

//
func NewReportViewGenerator(baseStructName string, fields []string) (*ReportViewGenerator) {
	return &ReportViewGenerator{BaseStructName: baseStructName, Fields: fields}
}

//
func (rvg *ReportViewGenerator) MakeInstanceName() {
	rvg.InstanceName = util.ToLowerCamelCase(rvg.BaseStructName + "Report")
}

//
func (rvg *ReportViewGenerator) MakeStructInstanceName() {
	rvg.StructAcronym = util.ToLowerAcronym(rvg.BaseStructName + "Report")
}

//
func  (rvg *ReportViewGenerator) IsHiddenField(fieldName string) bool {
	return fieldName == "Id" || fieldName == "AccountId" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "DeletedAt"
}


//
func (rvg *ReportViewGenerator) GetField(str string) string {
	parts := strings.Split(str, " as ")

	// alias
	if len(parts) > 1 {
		return parts[1]
	}

	parts = strings.Split(str, ".")

	if len(parts) > 1 {
		return parts[1]
	}

	return str
}

//
func (rvg *ReportViewGenerator) GenerateView() string {
	var code string

	code += rvg.GetViewTitleCode()
	code += rvg.GetViewBodyCode()
	code += rvg.GetViewScriptCode()

	return code
}

//
func (rvg *ReportViewGenerator) GetViewTitleCode() string {
	var code string

	code += fmt.Sprintf("[[define \"title\"]]%s[[end]]\n", util.ToWords(rvg.BaseStructName))

	return code
}

//
func (rvg *ReportViewGenerator) GetViewBodyCode() string {
	var code string
	var inputs string
	var header string
	var cells string

	for key := range rvg.Fields {
		fieldName := rvg.GetField(rvg.Fields[key])

		if rvg.IsHiddenField(fieldName) {
			inputs += fmt.Sprintf(`
			<input type="hidden" name="%s" v-model="%s"/>`, fieldName, fieldName)
		} else {
			inputs += fmt.Sprintf(`
			<div class="col-sm-3 my-1">
				<div class="form-group">
					<label for="%s">%s</label>
					<input type="text" class="form-control" id="%s" name="%s" aria-describedby="%sHelp" placeholder="Enter %s" v-model="filter.%s">
				</div>
			</div>
`, fieldName, util.ToWords(util.ToCamelCase(fieldName) ), fieldName, fieldName, fieldName, util.ToWords(util.ToCamelCase(fieldName) ), fieldName)

			header += fmt.Sprintf(`
				<th scope="col">%s</th>`, util.ToWords( util.ToCamelCase(fieldName)) )

			cells += fmt.Sprintf(`
				<td>{{ item.%s.String }}</td>`, fieldName)
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

<div v-cloak v-show="results.length > 0">
	<table class="table table-striped">
		<thead>
			<tr>
%s
			</tr>
		</thead>
		<tbody>
			<tr v-for="item in results">
%s
			</tr>
		</tbody>
	</table>
</div>

[[end]]
`, rvg.InstanceName, util.ToSnakeCase(rvg.BaseStructName), inputs, header, cells)

	return code
}

//
func (rvg *ReportViewGenerator) GetViewScriptCode() string {
	var code string
	var fields []string

	for key := range rvg.Fields {
		fieldName := rvg.GetField(rvg.Fields[key])
		fields = append(fields, fmt.Sprintf("\t\t\t\t%s: \"\"", fieldName) )
	}

	fieldCode := strings.Join(fields, ",\n")

	code += fmt.Sprintf(`
[[define "scripts"]]
<script>
	var app = new Vue({
		el: '#app',
		data: {
			offset: 0,
			offsetCount: 1,
			count: 0,
			filter: {
%s
			},
			query: "",
			results: []
		},
		methods: {
			submitForm: function() {
				this.updateQuery();
				this.getCount();
				this.getData();
			},
			updateQuery: function() {
				var tmp = [];

				for (var key in this.filter) {
					tmp.push(` + "`${key}=${this.filter[key]}`" + `);
				}

				this.query = encodeURI("?" + tmp.join("&"));
			},
			getData: function() {
				var url = "/ajax_%s" + this.query + encodeURI(` + "`&offset=${this.offset}`" + `);

				axios({ method: "GET", "url": url }).then(result => {
					this.results = (result.data != null) ? result.data : [];
				}, error => {
					console.error(error);
				});
			},
			getCount: function() {
				var url = "/ajax_%s" + this.query + encodeURI("&count=true");

				axios({ method: "GET", "url": url }).then(result => {
					this.count = (result.data != null) ? result.data : 0;
				}, error => {
					console.error(error);
				});
			}
		}
	})
</script>
[[end]]
`, fieldCode, util.ToSnakeCase(rvg.BaseStructName), util.ToSnakeCase(rvg.BaseStructName))

	return code
}
