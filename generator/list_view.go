package generator

import (
	"fmt"
	// "strings"
	"github.com/jschneider98/jgoweb/util"
)

//
func (mg *ModelGenerator) GenerateListView() string {
	var code string

	code += mg.GetListViewTitleCode()
	code += mg.GetListViewSubNavCode()
	code += mg.GetListViewBodyCode()
	code += mg.GetListViewScriptCode()

	return code
}

//
func (mg *ModelGenerator) GetListViewTitleCode() string {
	var code string

	code += fmt.Sprintf("[[define \"title\"]]%s List[[end]]\n", util.ToWords(mg.ModelName))

	return code
}

//
func (mg *ModelGenerator) GetListViewSubNavCode() string {

	code := `
[[define "sub-nav"]]
<div class="p-2 text-muted" style="border-top: 1px solid #ddd">
	<form>
		<div>
			<div class="col">
				<input class="form-control" id="filter" type="text" v-model="query" v-on:keyup="updateFilter" placeholder="Filter"></input>
			</div>
		</div>
	</form>
</div>
[[end]]
`

	return code
}

//
func (mg *ModelGenerator) GetListViewBodyCode() string {

	tmplParams := struct {
		Mg *ModelGenerator
	}{}

	tmplParams.Mg = mg

	str := `
[[define "body"]]
	<div class="p-2 text-muted" v-cloak>
		<h5 class="text-muted" style="font-size: calc(12px + 0.5vw)">{{ filteredResults.length }} out of {{ results.length }} total</h5>
	</div>

	<div v-cloak v-show="filteredResults.length > 0">
		<table class="table table-striped">
			<thead>
				<tr>?{{range $val := .Mg.Fields}}?
					<th scope="col">?{{$val.FieldName}}?</th>?{{end}}?
				</tr>
			</thead>
			<tbody>

				<tr v-for="item in filteredResults">?{{range $val := .Mg.Fields}}?
					<td>item.?{{ $val.FieldName}}?</td>?{{end}}?
				</tr>

			</tbody>
		</table>

		<h5 class="text-muted" style="font-size: calc(12px + 0.5vw)">{{ filteredResults.length }} out of {{ results.length }} total subscribers</h5>
	</div>
	<hr>
[[end]]
`
	code, err := util.TemplateToString(str, tmplParams)

	if err != nil {
		code = fmt.Sprintf("\n**** ERROR ****\n%s\n********\n", err)
	}

	return code
}

//
func (mg *ModelGenerator) GetListViewScriptCode() string {

	str := `
[[define "scripts"]]
<script>
	var app = new Vue({
		el: '#app',
		data: {
			query: "",
			results: [[ .Results ]],
			filteredResults: [[ .Results ]],
			loading: false
		},
		methods: {
			updateFilter: function() {
				var results = [];

				results = this.queryFilter();

				this.filteredResults = results;
			},
			queryFilter: function() {
				var results = [];

				for (var i = 0; i < this.results.length; i++) {
					if (this.queryMatch(this.results[i])) {
						results.push(this.results[i])
					}
				}

				return results;
			},
			queryMatch: function(row) {
				var parts = (this.query.toLowerCase()).split(" ");

				if (parts.length == 1) {
					if ( (row.FirstName.toLowerCase()).includes(parts[0]) ) {
						return true;
					}

					if ( (row.LastName.toLowerCase()).includes(parts[0]) ) {
						return true;
					}

					if ( (row.SchoolStateId.toLowerCase()).includes(parts[0]) ) {
						return true;
					}
				}

				if (parts.length > 1) {
					if ( (row.FirstName.toLowerCase()).includes(parts[0]) && (row.LastName.toLowerCase()).includes(parts[1])) {
						return true;
					}

					if ( (row.FirstName.toLowerCase()).includes(parts[1]) && (row.LastName.toLowerCase()).includes(parts[0])) {
						return true;
					}
				}

				return false;
			}
		}
	})
</script>
[[end]]
`
	code, err := util.TemplateToString(str, nil)

	if err != nil {
		code = fmt.Sprintf("\n**** ERROR ****\n%s\n********\n", err)
	}

	return code
}
