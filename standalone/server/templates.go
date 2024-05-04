package server

import "html/template"

var home = template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://cdn.simplecss.org/simple.min.css">
	<script src="https://unpkg.com/htmx.org@1.9.12"></script>
    <style>
		div.center {
			text-align: center;
		}

		div#results {
			height: 60dvh;
		}

		div#results > ul {
			list-style-type: none;
			overflow-y: auto;
			height: 100%;
		}
    </style>
    <title>Anagram and Word Search</title>
</head>
<body>
  <header>
    <h1>Anagram and Word Search</h1>
  </header>

  <main>
	<form action="/search" method="post" hx-boost="true" hx-target="#results" hx-push-url="false">
	  <div class="center">
	    <input id="pattern" type="text" name="pattern" required autofocus></input>
	    <button name="mode" value="match">Match</button>
	    <button name="mode" value="anagrams">Anagrams</button>
		<button type="reset" onclick="getfocus()">Clear</button>
	  </div>
    </form>
	<div id="results">
	</div>
  </main>

</body>
<script>
function getfocus() {
  document.getElementById("pattern").focus();
}
</script>
</html>
`))

type ResultParams struct {
	Preamble string
	Results  []string
}

var results = template.Must(template.New("results").Parse(`
{{ with .Preamble }}
<p>{{ . }}</p>
{{ end }}
<ul>
{{ range .Results }}
  <li>
    <a href="https://dicoweb.gnu.org.ua/?q={{.}}&db=gcide&define=1" target="defn">{{.}}</a>
  </li>
{{ end }}
</ul>
`))
