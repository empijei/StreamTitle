package main

import (
	"io"
	"net/http"
	"text/template"
)

const pageHead = `
<html>
<head>
	<link rel="stylesheet" href="/static/styles.css">
	<link rel="preconnect" href="https://fonts.gstatic.com">
	<link href="https://fonts.googleapis.com/css2?family=Roboto&display=swap" rel="stylesheet">
	<link href="https://fonts.googleapis.com/css2?family=VT323&display=swap" rel="stylesheet">
</head>
	<body><div class="content">`
const pageFoot = `
	</div></body>
</html>
`
const loginForm = `
		<h1> Welcome to SecretBook </h1>
		<p>
			<form action="/login" method="POST">
			<label for="uname">Username:</label><br>
			<input type="text" id="uname" name="uname" value=""><br>
			<label for="password">Password:</label><br>
			<input type="password" id="password" name="password"><br><br>
			<input type="submit" value="Submit">
			</form> 
		</p>
		`

func renderPage(w http.ResponseWriter, body string) {
	io.WriteString(w, pageHead)
	io.WriteString(w, body)
	io.WriteString(w, pageFoot)
}

var listNotesTpl = template.Must(template.New("notes page").Parse(`
<h3>Viewing notes of "{{.ViewedUser}}" as "{{ if .VisitingUser }}{{.VisitingUser}}{{else}}Guest{{end}}":</h3>
<p>
{{ if ne 0 (len .Notes) }}
<ul>
{{ range .Notes }}
	<li>{{.Title}}</li>
{{ end }}
</ul>
{{ else }}
There are no notes.
{{ end }}
</p>
`))

var myNotesTpl = template.Must(template.New("notes page").Parse(`
<h1> Welcome {{.User}} </h1>
<h2> Create a note </h2>
<form action="/addnote" id="addnote_form">
  Title: <input type="text" name="title">
  <input type="checkbox" name="private" value="Private">
  <label for="private">Private</label><br>
  <input type="submit">
</form>
<textarea class="typing" id="textarea_form" rows="4" cols="50" name="text" form="text">
</textarea>
<br>
<button onclick='addSignature("{{.User}}")'>add signature</button>
<div>
{{ if ne 0 (len .Notes) }}
<p>
Or see your notes:
</p>
{{ range .Notes }}
	<h2>{{.Title}}</h2>
	<p class="typing">
	{{.Text}}
	</p>
{{ end }}
{{ else }}
<p>
There are no notes.
</p>
{{ end }}
</div>
<script>
function addSignature(uname){
	const textarea = document.getElementById("textarea_form");
	textarea.value = textarea.value + "\nYours: " + uname;
}
</script>
`))
