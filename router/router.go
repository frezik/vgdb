package router

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func Routes() *chi.Mux {
    r := chi.NewRouter()
    r.Use( middleware.Logger )

    r.Get( "/", HeartBeat )

    return r
}

func HeartBeat(
    w http.ResponseWriter,
    r *http.Request,
) {
    w.Write([]byte("OK"))
}
