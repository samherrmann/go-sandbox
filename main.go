package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/samherrmann/go-sandbox/graphql"
)

func main() {
	if err := app(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func app() error {

	config, err := readConfig("config.json")
	if err != nil {
		return err
	}

	ordersTpl, err := template.ParseFiles("orders.tpl.html")
	if err != nil {
		return err
	}

	http.HandleFunc("/", ordersHandler(config, ordersTpl))

	return http.ListenAndServe(":8080", nil)
}

func ordersHandler(config *Config, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		viewData := &ViewData{}

		query, err := os.ReadFile("orders.graphql")
		if err != nil {
			viewData.Error = err.Error()
			tpl.Execute(w, viewData)
			return
		}

		apiReq := &graphql.Request{
			Query: string(query),
		}

		res, err := graphql.Send(config.Shop, config.AccessToken, apiReq)
		if err != nil {
			viewData.Error = err.Error()
			tpl.Execute(w, viewData)
			return
		}

		m := map[string]any{}
		if err := json.NewDecoder(res).Decode(&m); err != nil {
			viewData.Error = err.Error()
			tpl.Execute(w, viewData)
			return
		}

		viewData.Main = m
		tpl.Execute(w, viewData)
	}
}

type ViewData struct {
	Error string
	Main  map[string]any
}
