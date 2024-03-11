package pages

import (
	"embed"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/samherrmann/go-sandbox/models"
	"github.com/samherrmann/go-sandbox/pages/internal"
)

const (
	homePath = "/todo"
)

//go:embed home.html home.css
var homeFS embed.FS

func NewHome(logger *slog.Logger) (*Home, error) {
	tpl, err := internal.ParseTemplate(homeFS, "home.html")
	if err != nil {
		return nil, err
	}

	return &Home{
		todos:  &models.ToDo{},
		logger: logger,
		tpl:    tpl,
	}, nil
}

type Home struct {
	tpl    *internal.Template
	logger *slog.Logger
	todos  *models.ToDo
}

func (h *Home) GetToDos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := internal.Page{
			Title: "To Do",
			Path:  homePath,
			Data:  h.todos,
		}
		page.AddStyleSheet(homeFS, "home.css")
		if err := h.tpl.Execute(w, page); err != nil {
			h.logger.Error(err.Error())
		}
	}
}

func (h *Home) AddToDo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.todos.Append(r.Form.Get("value"))
		w.Header().Add("Location", homePath)
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (h *Home) UpdateToDo() http.HandlerFunc {
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

		w.Header().Add("Location", homePath)
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (h *Home) RemoveToDo() http.HandlerFunc {
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

		w.Header().Add("Location", homePath)
		w.WriteHeader(http.StatusSeeOther)
	}
}
