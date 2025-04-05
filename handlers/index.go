package handlers

import (
	"html/template"
	"net/http"
)

func HandleGETIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")

	tmpl := template.Must(template.New("template").Funcs(CustomTemplateFuncs).ParseFiles("components/pages/index.html", "components/blocks/search_form.html", "components/blocks/results.html", "components/blocks/results_items.html", "components/base.html"))

	data := PageData{
		Title: "Go Games Browser",
		Intro: "Main page for the Go Games Browser app",
	}

	tmpl.ExecuteTemplate(w, "base", data)
}
