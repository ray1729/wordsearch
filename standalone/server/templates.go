package server

import "html/template"

var home = template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" href="/assets/simple.css">
	<link rel="stylesheet" href="/assets/custom.css">
	<script src="/assets/htmx.min.js"></script>
    <title>Anagram and Word Search</title>
</head>
<body>
  <header>
    <h1>Anagram and Word Search</h1>
  </header>

  <main>
	<form action="/search" method="post" hx-boost="true" hx-target="#results">
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
  <li>{{.}}
    <a href="https://dicoweb.gnu.org.ua/?q={{.}}&db=gcide&define=1" target="defn"><span class="small">GCIDE</span></a>
    <a href="https://dicoweb.gnu.org.ua/?q={{.}}&db=WordNet&define=1" target="defn"><span class="small">WordNet</span></a>
  </li>
{{ end }}
</ul>
`))
