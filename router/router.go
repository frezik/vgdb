package router

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"

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

type SystemData struct {
    Name string
    Publisher string
    FirstRelease string
    JPRelease string
    NARelease string
    EURelease string
    BRRelease string
}


func Routes() *chi.Mux {
    r := chi.NewRouter()
    r.Use( middleware.Logger )

    r.Get( "/", HeartBeat )
    r.Get( "/systems", ListSystems )
    r.Get( "/{system}/games", ListSystemGames )

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
        log.Println( err )
        return
    }
}

func ListSystemGames(
    w http.ResponseWriter,
    r *http.Request,
) {
    system := chi.URLParam( r, "system" )
    system_data_file, does_system_exist := data_files[ system ]
    if ! does_system_exist {
        http.Error( w, "System not found", http.StatusNotFound )
        log.Println( "System does not exist: " + system )
        return
    }

    system_data, err := GetSystemData( system_data_file )
    if err != nil {
        http.Error( w, "Internal error", http.StatusInternalServerError )
        log.Println( err )
        return
    }

    games := make( []string, len( system_data ) )
    for k := range system_data {
        games[k] = system_data[k].Name
    }

    output := map[string][]string {
        "games": games,
    }

    err = json.NewEncoder( w ).Encode( output )
    if err != nil {
        http.Error( w, "Internal error", http.StatusInternalServerError )
        log.Println( err )
        return
    }
}

func GetSystemData(
    data_file string,
) ([]SystemData, error) {
    file_path := filepath.Join( data_dir, data_file )
    data, err := os.ReadFile( file_path )
    if err != nil {
        return nil, err
    }

    var system_data []SystemData
    err = json.Unmarshal( data, &system_data )
    if err != nil {
        return nil, err
    }

    return system_data, nil
}
