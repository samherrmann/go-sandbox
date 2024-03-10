package pages

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/samherrmann/go-sandbox/pages/internal"
)

//go:embed home.html
var homeFS embed.FS

func NewHome(logger *slog.Logger) (*Home, error) {
	tpl, err := internal.ParseTemplate(homeFS, "home.html")
	if err != nil {
		return nil, err
	}

	return &Home{
		todos:  []string{},
		logger: logger,
		tpl:    tpl,
	}, nil
}

type Home struct {
	tpl    *internal.Template
	logger *slog.Logger
	todos  []string
}

func (h *Home) GetToDos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := internal.Page{
			Title: "To Do",
			Data:  h.todos,
		}
		if err := h.tpl.Execute(w, page); err != nil {
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
