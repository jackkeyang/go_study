package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"model"
	"net"
	"net/http"
	"regexp"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

type Page model.Page

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page := model.Page{Title: title}
	data, err := page.SelectBlog()
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	p := Page{Title: data.Title, Body: data.Body}

	renderTemplate(w, "view", &p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	page := model.Page{Title: title}
	data, _ := page.SelectBlog()
	p := Page{Title: title, Body: data.Body}
	renderTemplate(w, "edit", &p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := model.Page{Title: title, Body: []byte(body)}
	_, err := p.SelectBlog()
	if err != nil {
		p.InsertBlog()
	} else {
		p.UpdateBlog()
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	p := model.Page{}
	pages := p.SelectAll()
	t, _ := template.ParseFiles("tmpl/list.html")
	t.Execute(w, pages)
}

func main() {
	flag.Parse()
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/list/", listHandler)

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8000", nil)
}
