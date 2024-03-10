package pages

import (
	"log/slog"
	"net/http"

	"github.com/samherrmann/go-sandbox/templates"
)

func NewHome(logger *slog.Logger, tpls *templates.Templates) *Home {
	tpls.Add(embedded, "home.html")

	return &Home{
		todos:  []string{},
		logger: logger,
		tpls:   tpls,
	}
}

type Home struct {
	tpls   *templates.Templates
	logger *slog.Logger
	todos  []string
}

func (h *Home) GetToDos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := templates.Page{
			Title: "To Do",
			Data:  h.todos,
		}
		if err := h.tpls.ExecuteTemplate(w, "home.html", page); err != nil {
			h.logger.Error(err.Error())
		}
	}
}

func (h *Home) AddToDo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Error(err.Error())
			return
		}
		value := r.Form.Get("value")
		if value != "" {
			h.todos = append(h.todos, value)
		}
		w.Header().Add("Location", r.URL.Path)
		w.WriteHeader(http.StatusSeeOther)
	}
}
