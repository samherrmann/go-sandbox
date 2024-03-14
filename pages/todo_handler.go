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

func NewTodoHandler(path string, logger *slog.Logger) (http.Handler, error) {
	tpl, err := view.ParseTemplate(todoFS, "todo.html")
	if err != nil {
		return nil, err
	}

	todo := &ToDoHandler{
		path:   path,
		todos:  &models.ToDo{},
		logger: logger,
		tpl:    tpl,
	}
	return todo, nil
}

type ToDoHandler struct {
	path   string
	tpl    *view.Template
	logger *slog.Logger
	todos  *models.ToDo
}

func (h *ToDoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.read(w)
	case http.MethodPost:
		if err := h.parseForm(w, r); err != nil {
			return
		}
		action := r.Form.Get("action")
		switch action {
		case "create":
			h.create(w, r)
		case "update":
			h.update(w, r)
		case "delete":
			h.delete(w, r)
		default:
			err := fmt.Errorf("unknown action %q", action)
			h.renderViewWithError(w, http.StatusBadRequest, err)
		}
	default:
		http.NotFound(w, r)
	}
}

func (h *ToDoHandler) read(w http.ResponseWriter) {
	h.renderView(w, http.StatusOK)
}

func (h *ToDoHandler) create(w http.ResponseWriter, r *http.Request) {
	v := r.Form.Get("value")
	h.todos.Append(v)
	httputil.SeeOther(w, h.path)
}

func (h *ToDoHandler) update(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		h.renderViewWithError(w, http.StatusBadRequest, err)
		return
	}
	value := r.Form.Get("value")
	if err := h.todos.Update(index, value); err != nil {
		h.renderViewWithError(w, http.StatusBadRequest, err)
		return
	}
	httputil.SeeOther(w, h.path)
}

func (h *ToDoHandler) delete(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		h.renderViewWithError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.todos.Remove(index); err != nil {
		h.renderViewWithError(w, http.StatusBadRequest, err)
		return
	}

	httputil.SeeOther(w, h.path)
}

func (h *ToDoHandler) parseForm(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		h.renderViewWithError(w, http.StatusBadRequest, err)
	}
	return err
}

func (h *ToDoHandler) renderView(w http.ResponseWriter, statusCode int) {
	h.renderViewWithError(w, statusCode, nil)
}

func (h *ToDoHandler) renderViewWithError(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		h.logger.Error(err.Error())
	}

	v := &view.View{
		Title: "To Do",
		Path:  h.path,
		Data:  h.todos,
	}

	if err != nil {
		v.Error = err.Error()
	}
	v.AddStyleSheet(todoFS, "todo.css")

	w.WriteHeader(statusCode)
	if err := h.tpl.Execute(w, v); err != nil {
		h.logger.Error(err.Error())
	}
}
