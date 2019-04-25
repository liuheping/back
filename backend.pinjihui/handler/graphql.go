package handler

import (
	"encoding/json"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/loader"
	rp "pinjihui.com/backend.pinjihui/repository"
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

	queryfiled, _ := h.Schema.GetQueryFields(params.Query, params.OperationName)
	ctx = context.WithValue(ctx, "query_filed", queryfiled)
	// 插入操作日志
	rp.L("public").(*rp.PublicRepository).CreateOperationLog(ctx)

	response := h.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
