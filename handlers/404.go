package handlers

import (
	"html/template"
	"net/http"
)

func Handle404(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	w.WriteHeader(http.StatusNotFound)

	tmpl := template.Must(template.New("template").Funcs(CustomTemplateFuncs).ParseFiles("components/pages/404.html", "components/base.html"))

	tmpl.ExecuteTemplate(w, "base", nil)
}
