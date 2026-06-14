package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type server struct {
	contentDir string
	layout		 *template.Template
}

func newServer(contentDir, layoutPath string) (*server, error) {
	tmpl, err := template.ParseFiles(layoutPath)
	if err != nil {
		return nil, err
	}
	return &server{contentDir: contentDir, layout: tmpl}, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Print("Usage: <contents-path>")
		return
	}

	fileinfo, err := os.Stat(os.Args[1])
	if err != nil {
		log.Print("Problem with the paht?")
		return
	}
	if !fileinfo.IsDir() {
		log.Print("A directory must be provided")
		return 
	}

	// provide the absolute due to this folder might be out of the blue.
	contents_path, err := filepath.Abs(os.Args[1])
	if err != nil {
		return
	}	

	s, err := newServer(contents_path, "templates/layout.html")
	if err != nil {
		return
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", s.serveTemplate)

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *server) serveTemplate(w http.ResponseWriter, r *http.Request) {
	// Determine content file: root → about, otherwise use path + .html
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "about"
	}
	contentPath := filepath.Join(s.contentDir, path + ".html")

	raw, err := os.ReadFile(contentPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	content := extractBody(raw)

	data := map[string]interface{}{
		"Content": template.HTML(content),
	}
	s.layout.Execute(w, data)
}

// extractBody grabs the inner HTML from a full document (everything between <body> and </body>).
// as typst generates also the <html> and <header> lables we have to just keep the body, as
// the layout just needs a body.
func extractBody(raw []byte) template.HTML {
	// Find <body> opening tag
	openStart := bytes.Index(bytes.ToLower(raw), []byte("<body"))
	if openStart == -1 {
		return template.HTML(raw)
	}
	openEnd := bytes.IndexByte(raw[openStart:], '>') // we split it in two bcs there is options in here
	if openEnd == -1 {
		return template.HTML(raw)
	}
	bodyStart := openStart + openEnd + 1

	// Find </body> closing tag
	bodyEnd := bytes.LastIndex(bytes.ToLower(raw), []byte("</body>"))
	if bodyEnd == -1 {
		return template.HTML(raw)
	}

	return template.HTML(bytes.TrimSpace(raw[bodyStart:bodyEnd]))
}
