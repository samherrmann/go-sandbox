package pages

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/samherrmann/go-sandbox/httperror"
	"github.com/samherrmann/go-sandbox/models"
	"github.com/samherrmann/go-sandbox/view"
)

//go:embed todo.html todo.css
var todoFS embed.FS

func NewTodoHandler() (httperror.Handler, error) {
	tpl, err := view.ParseTemplate(todoFS, "todo.html")
	if err != nil {
		return nil, err
	}

	todo := &ToDoHandler{
		todos: &models.ToDo{},
		tpl:   tpl,
	}
	return todo, nil
}

type ToDoHandler struct {
	tpl   *view.Template
	todos *models.ToDo
}

func (h *ToDoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	code, err := h.serveHTTP(r)
	switch code {
	case http.StatusNotFound:
		http.NotFound(w, r)
	case http.StatusSeeOther:
		w.Header().Add("Location", r.URL.Path)
		w.WriteHeader(http.StatusSeeOther)
	default:
		h.renderView(w, code, r.URL.Path, err)
	}
	return err
}

func (h *ToDoHandler) serveHTTP(r *http.Request) (int, error) {
	switch r.Method {
	case http.MethodGet:
		return http.StatusOK, nil
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			return http.StatusBadRequest, err
		}
		action := r.Form.Get("action")
		switch action {
		case "create":
			return h.create(r)
		case "update":
			return h.update(r)
		case "delete":
			return h.delete(r)
		default:
			return http.StatusBadRequest, fmt.Errorf("unknown action %q", action)
		}
	default:
		return http.StatusNotFound, errors.New("404 not found")
	}
}

func (h *ToDoHandler) create(r *http.Request) (int, error) {
	v := r.Form.Get("value")
	h.todos.Append(v)
	return http.StatusSeeOther, nil
}

func (h *ToDoHandler) update(r *http.Request) (int, error) {
	index, err := parseIndex(r.Form.Get("index"))
	if err != nil {
		return http.StatusBadRequest, err
	}
	value := r.Form.Get("value")
	if err := h.todos.Update(index, value); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusSeeOther, nil
}

func (h *ToDoHandler) delete(r *http.Request) (int, error) {
	index, err := parseIndex(r.Form.Get("index"))
	if err != nil {
		return http.StatusBadRequest, err
	}
	if err := h.todos.Remove(index); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusSeeOther, nil
}

func (h *ToDoHandler) renderView(w http.ResponseWriter, statusCode int, path string, err error) error {
	v := &view.View{
		Title: "To Do",
		Path:  path,
		Data:  h.todos,
		Error: httperror.String(err),
	}

	// TODO: Figure out a way that the stylesheet doesn't need to be added on
	// every request.
	v.AddStyleSheet(todoFS, "todo.css")

	w.WriteHeader(statusCode)
	if err := h.tpl.Execute(w, v); err != nil {
		return err
	}

	return err
}

func parseIndex(str string) (int, error) {
	index, err := strconv.Atoi(str)
	if err != nil {
		return 0, httperror.New(fmt.Sprintf("invalid index %v", str), err)
	}
	return index, nil
}
