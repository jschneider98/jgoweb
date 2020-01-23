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

	tmplParams := struct {
		ModelName string
	}{}

	tmplParams.ModelName = util.ToSnakeCase(mg.ModelName)

	str := `
[[define "sub-nav"]]
<div class="p-2 text-muted" style="border-top: 1px solid #ddd">
	<form>
		<div class="form-row">
			<div class="col-auto">
				<a class="nav-link" href="/{{{ .ModelName }}}">Create</a>
			</div>
		
			<div class="col">
				<input class="form-control" id="filter" type="text" v-model="query" v-on:keyup="updateFilter" placeholder="Filter List"></input>
			</div>
		</div>
	</form>
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
func (mg *ModelGenerator) GetListViewBodyCode() string {

	tmplParams := struct {
		Mg *ModelGenerator
	}{}

	tmplParams.Mg = mg

	str := `
[[define "body"]]
	<div class="p-2 text-muted" v-cloak>
		<h5 class="text-muted" style="font-size: calc(12px + 0.5vw)">{{ filteredResults.length }} out of {{ results.length }}</h5>
	</div>

	<div v-cloak v-show="filteredResults.length > 0">
		<table class="table table-striped">
			<thead>
				<tr>
					<th class="no-print"></th>{{{range $val := .Mg.Fields}}}
					<th scope="col">{{{$val.FieldName}}}</th>{{{end}}}
				</tr>
			</thead>
			<tbody>

				<tr v-for="item in formatData">
					<td class="no-print">
						<div style="min-width: 50px;">
							<span class="text-primary">
								<a :href="getEditRoute(item.Id.String)"><i class="fa fa-edit fa-lg"></i></a>
								<a href="#"><i v-on:click="deleteConfirm(item.Id.String)" class="fa fa-trash fa-lg"></i></a>
							</span>
						</div>
					</td>{{{range $val := .Mg.Fields}}}
					<td>{{ item.{{{ $val.FieldName.String }}} }}</td>{{{end}}}
				</tr>

			</tbody>
		</table>

		<h5 class="text-muted" style="font-size: calc(12px + 0.5vw)">{{ filteredResults.length }} out of {{ results.length }}</h5>
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

	tmplParams := struct {
		EditRoute string
		ListRoute string
		Id string
	}{}

	tmplParams.EditRoute = util.ToSnakeCase(mg.ModelName)
	tmplParams.ListRoute = tmplParams.EditRoute + `_list`
	tmplParams.Id = util.ToLowerCamelCase( util.ToSnakeCase(mg.ModelName + `Id`) )

	str := `
[[define "scripts"]]
<script>
	window.addEventListener('load', (event) => {
		app.updateFilter();
	});

	var app = new Vue({
		el: '#app',
		data: {
			isCurrent: "true",
			query: "",
			results: [[ .Results ]],
			filteredResults: [[ .Results ]],
			loading: false
		},
		computed: {
			formatData: function() {

				for (let i = 0; i < this.filteredResults.length; i++) {
					this.filteredResults[i].CreatedAt.String = medex.formatDate(this.filteredResults[i].CreatedAt.String, false);

					this.filteredResults[i].UpdatedAt.String = medex.formatDate(this.filteredResults[i].UpdatedAt.String, false);
				}

				return this.filteredResults;
			}
		},
		methods: {
			getEditRoute: function({{{.Id}}}) {
				return "/{{{.EditRoute}}}?{{{.Id}}}="  + encodeURI({{{.Id}}});
			},
			getDeleteLink: function({{{.Id}}}) {
				var uri = encodeURI(` + "`?{{{.Id}}}=${ {{{.Id}}} }&action=delete`" + `)
				return "/{{{.ListRoute}}}" + uri;
			},
			updateFilter: function() {
				return this.queryFilter();

			},
			queryFilter: function() {
				var results = [];

				for (var i = 0; i < this.results.length; i++) {
					if (this.queryMatch(this.results[i])) {
						results.push(this.results[i])
					}
				}

				this.filteredResults = results;
			},
			queryMatch: function(row) {
				return (row.Diagnosis.String.toLowerCase()).includes(this.query.toLowerCase());
			},
			deleteConfirm({{{.Id}}}) {
				var url = this.getDeleteLink({{{.Id}}});

				this.$bvModal.msgBoxConfirm("Are you sure that you want to delete?", {
						headerBgVariant: 'danger',
						headerTextVariant: 'light',
						title: "WARNING: About to Delete Data",
						okVariant: "light",
						cancelVariant: "primary",
						okTitle: "Delete",
						cancelTitle: "No",
						footerClass: "p-2",
						hideHeaderClose: false,
						centered: true
					})
					.then(value => {
						if (value == true) {
							window.location.replace(url);
						}
					})
					.catch(err => {
						// An error occurred
					})
			}
		}
	});
</script>
[[end]]
`
	code, err := util.TemplateToString(str, tmplParams)

	if err != nil {
		code = fmt.Sprintf("\n**** ERROR ****\n%s\n********\n", err)
	}

	return code
}
