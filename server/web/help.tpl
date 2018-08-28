<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
	    <h1>{{.Title}}</h1>
	    <div>
		{{range .Urls}}
		    <div>curl -v {{ . }}</div>
		{{end}}
		</div>
	</body>
</html>