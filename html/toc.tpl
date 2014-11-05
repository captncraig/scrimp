<table class="table-bordered table-striped table-condensed">
	<thead><th>Module</th><th>Services</th><th>Data types</th><th>Constants</th></thead>
	<tr>
		<td>{{.Name}}</td>
		<td>{{range .Services}}
			<a href="#Svc_{{.Name}}">{{.Name}}</a><br/>
			<ul>
				{{$svc := .Name}}
				{{range .Functions}}
					<li><a href="#Fn_{{$svc}}_{{.Name}}">{{.Name}}</a></li>
				{{end}}
			</ul>
		{{end}}</td>
		<td>
			{{range .Typedefs}}<a href="#Typedef_{{.Name}}">{{.Name}}</a><br/>{{end}}
			{{range .Enums}}<a href="#Enum_{{.Name}}">{{.Name}}</a><br/>{{end}}
			{{range .Objects}}<a href="#Struct_{{.Name}}">{{.Name}}</a><br/>{{end}}
		</td>
		<td>{{range .Consts}}<code><a href="#Const_{{.Name}}">{{.Name}}</a></code><br/>{{end}}</td>
	</tr>
</table>