package pages

import (
	"embed"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/samherrmann/go-sandbox/httputil"
	"github.com/samherrmann/go-sandbox/models"
	"github.com/samherrmann/go-sandbox/view"
)

//go:embed todo.html todo.css
var todoFS embed.FS

func NewTodoHandler(path string, logger *slog.Logger) (http.Handler, error) {
	tpl, err := view.ParseTemplate(todoFS, "todo.html")
	if err != nil {
		return nil, err
	}

	todo := &ToDo{
		path:   path,
		todos:  &models.ToDo{},
		logger: logger,
		tpl:    tpl,
	}

	return todo, nil
}

type ToDo struct {
	path   string
	tpl    *view.Template
	logger *slog.Logger
	todos  *models.ToDo
}

func (h *ToDo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.Handle("GET /", httputil.ExactPathHandler("/", h.read()))
	mux.Handle("POST /create", h.create())
	mux.Handle("POST /{id}/update", h.update())
	mux.Handle("POST /{id}/delete", h.delete())
	mux.ServeHTTP(w, r)
}

func (h *ToDo) read() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view := &view.View{
			Title: "To Do",
			Path:  h.path,
			Data:  h.todos,
		}
		view.AddStyleSheet(todoFS, "todo.css")
		if err := h.tpl.Execute(w, view); err != nil {
			h.logger.Error(err.Error())
		}
	})
}

func (h *ToDo) create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.todos.Append(r.Form.Get("value"))

		httputil.SeeOther(w, h.path)
	})
}

func (h *ToDo) update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		httputil.SeeOther(w, h.path)
	})
}

func (h *ToDo) delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		httputil.SeeOther(w, h.path)
	})
}
