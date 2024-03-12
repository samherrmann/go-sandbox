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

	mux := http.NewServeMux()

	todo := &ToDo{
		path:   path,
		todos:  &models.ToDo{},
		logger: logger,
		tpl:    tpl,
		mux:    mux,
	}

	mux.Handle("GET /{$}", todo.read())
	mux.Handle("POST /create", todo.create())
	mux.Handle("POST /{id}/update", todo.update())
	mux.Handle("POST /{id}/delete", todo.delete())

	return todo, nil
}

type ToDo struct {
	path   string
	tpl    *view.Template
	logger *slog.Logger
	todos  *models.ToDo
	mux    *http.ServeMux
}

func (h *ToDo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *ToDo) read() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.tpl.Execute(w, h.newView()); err != nil {
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

func (h *ToDo) newView() *view.View {
	v := &view.View{
		Title: "To Do",
		Path:  h.path,
		Data:  h.todos,
	}
	v.AddStyleSheet(todoFS, "todo.css")
	return v
}
