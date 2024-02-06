package router

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

const data_dir = "data"
var data_files = map[string]string{
    "nes": "nes.json",
    "snes": "snes.json",
    "sega_genesis": "sega_genesis.json",
    "sega_master": "sega_master.json",
}


func Routes() *chi.Mux {
    r := chi.NewRouter()
    r.Use( middleware.Logger )

    r.Get( "/", HeartBeat )
    r.Get( "/systems", ListSystems )

    return r
}

func HeartBeat(
    w http.ResponseWriter,
    r *http.Request,
) {
    w.Write([]byte("OK"))
}

func ListSystems(
    w http.ResponseWriter,
    r *http.Request,
) {
    systems := make( []string, len( data_files ) )

    i := 0
    for k := range data_files {
        systems[i] = k
        i++
    }

    output := map[string][]string {
        "systems": systems,
    }

    err := json.NewEncoder( w ).Encode( output )
    if err != nil {
        http.Error( w, "Internal error", http.StatusInternalServerError )
        return
    }
}
