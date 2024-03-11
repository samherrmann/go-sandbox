package pages

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/samherrmann/go-sandbox/httputil"
	"github.com/samherrmann/go-sandbox/models"
	"github.com/samherrmann/go-sandbox/view"
)

//go:embed todo.html todo.css
var todoFS embed.FS

func NewToDo(mux *http.ServeMux, path string, logger *slog.Logger) (*ToDo, error) {
	tpl, err := view.ParseTemplate(todoFS, "todo.html")
	if err != nil {
		return nil, err
	}

	todo := &ToDo{
		Path:   path,
		todos:  &models.ToDo{},
		logger: logger,
		tpl:    tpl,
	}

	mux.HandleFunc(fmt.Sprintf("GET %s", path), todo.Get())
	mux.HandleFunc(fmt.Sprintf("POST %s/add", path), todo.Add())
	mux.HandleFunc(fmt.Sprintf("POST %s/{id}/update", path), todo.Update())
	mux.HandleFunc(fmt.Sprintf("POST %s/{id}/delete", path), todo.Delete())

	return todo, nil
}

type ToDo struct {
	Path   string
	tpl    *view.Template
	logger *slog.Logger
	todos  *models.ToDo
}

func (h *ToDo) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view := &view.View{
			Title: "To Do",
			Path:  h.Path,
			Data:  h.todos,
		}
		view.AddStyleSheet(todoFS, "todo.css")
		if err := h.tpl.Execute(w, view); err != nil {
			h.logger.Error(err.Error())
		}
	}
}

func (h *ToDo) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.todos.Append(r.Form.Get("value"))

		httputil.SeeOther(w, h.Path)
	}
}

func (h *ToDo) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		indexStr := r.PathValue("id")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		value := r.Form.Get("value")

		if err := h.todos.Update(index, value); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		httputil.SeeOther(w, h.Path)
	}
}

func (h *ToDo) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		indexStr := r.PathValue("id")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		if err := h.todos.Remove(index); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		httputil.SeeOther(w, h.Path)
	}
}
