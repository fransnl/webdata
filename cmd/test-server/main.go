package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/fransnl/webdata/webdata"
)

func main() {
	srv := http.Server{
		Addr:    ":42069",
		Handler: routes(),
	}

	log.Println("Listening on port 42069...")

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("===FATAL ERROR MAIN===\n", err)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)

	//htmx
	mux.HandleFunc("POST /webdata", htmxWebdata)

	fs := http.FileServer(http.Dir("./src/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}

func index(w http.ResponseWriter, r *http.Request){
	wd := webdata.WebInfo{}
	type data struct {
		Webdata webdata.WebInfo
	}
	d := data{wd}
	tmpl := template.Must(template.ParseFiles("src/html/base.html", "src/html/index.html"))
	err := tmpl.Execute(w, d)
	if err != nil {
		log.Println("INDEX TEMPLATE: ", err)
	}
}

func htmxWebdata(w http.ResponseWriter, r *http.Request){
	formUrl := r.FormValue("url")
	url := webdata.FixUrl(formUrl)

	wd := webdata.GetWebData(url)

	tmpl := template.Must(template.ParseFiles("src/html/index.html"))
	err := tmpl.ExecuteTemplate(w, "webdata", wd)
	if err != nil {
		log.Println("HTMX WEBDATA HANDLER: ", err)
	}
}
