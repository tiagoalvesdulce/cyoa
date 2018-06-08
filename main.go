package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>`

// Arc is a struct to hold all info about an Arc
type Arc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option represents a choice
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

// AdventureBook is the map of [string] -> Arc
type AdventureBook map[string]Arc

func decodeJSON(book *AdventureBook) {
	jsonFile, err := ioutil.ReadFile("story.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(jsonFile, &book)
	if err != nil {
		fmt.Println(err)
	}
}

// Handler is the http handler
type Handler struct {
	b    AdventureBook
	t    *template.Template
	path string
}

func getPath(req *http.Request) string {
	path := req.URL.Path
	if path == "/" || path == "" {
		path = "/intro"
	}
	return path[1:]
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.path = getPath(req)
	if arc, ok := h.b[h.path]; ok {
		err := h.t.Execute(w, arc)
		if err != nil {
			fmt.Fprintf(w, "Something went wrong... \n%s\n", err)
		}
		return
	}
	fmt.Fprint(w, "Chapter not found...\n")
}

func main() {
	var book AdventureBook
	decodeJSON(&book)
	h := Handler{book, tpl, ""}
	http.ListenAndServe(":8080", h)
}
