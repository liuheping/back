package handler

import (
	"encoding/json"
	"pinjihui.com/pinjihui/loader"
	"github.com/graph-gophers/graphql-go"
	"net/http"
    "github.com/nicksnyder/go-i18n/i18n"
	"strings"
)

type GraphQL struct {
	Schema  *graphql.Schema
	Loaders loader.LoaderCollection
}

func (h *GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := h.Loaders.Attach(r.Context())

	response := h.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
    TranslateResponse(response, r)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func TranslateResponse(response *graphql.Response, r *http.Request) {
    lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	T, _ := i18n.Tfunc(lang, accept, "zh-CN")
    if response != nil {
        for _, c := range response.Errors {
        	var ms []string
			for _, v := range strings.Split(c.Message, ";") {
				ms = append(ms, T(v))
			}
            c.Message = strings.Join(ms, ";")
        }
    }
}
