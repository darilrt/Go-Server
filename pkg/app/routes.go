package app

import (
    "go-server/pkg/api"
    "net/http"
)

func NewRouter() http.Handler {
    mux := http.NewServeMux();
    
    mux.HandleFunc("/tags", api.Tags)

    return mux
}