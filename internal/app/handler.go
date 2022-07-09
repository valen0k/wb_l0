package app

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func (a *App) NewHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/view", a.view)
	mux.HandleFunc("/show/", a.show)

	return mux
}

func (a *App) view(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("StatusMethodNotAllowed")
		return
	}
	tmpl, err := template.
		New("index.html").
		ParseFiles("web/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Can not expand template")
		return
	}
	err = tmpl.Execute(w, make(map[int]struct{}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Can not expand template")
		return
	}
}

func (a *App) show(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("StatusMethodNotAllowed")
		return
	}
	id := strings.Replace(r.URL.Path, "/show/", "", -1)
	model, ok := a.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(model)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
