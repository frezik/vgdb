// Package router sets up and handles the chi routes.
package router

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "github.com/frezik/vgdb/util"
)


// Creates a chi router, sets up the paths on it, and returns the router.
func Routes() *chi.Mux {
    r := chi.NewRouter()
    r.Use( middleware.Logger )

    r.Get( "/", HeartBeat )
    r.Get( "/systems", ListSystems )
    r.Get( "/{system}/games", ListSystemGames )

    return r
}

// A basic heart beat check to show if the endpoints are basically functional.
func HeartBeat(
    w http.ResponseWriter,
    r *http.Request,
) {
    w.Write([]byte("OK"))
}

// Outputs a list of available systems.
func ListSystems(
    w http.ResponseWriter,
    r *http.Request,
) {
    systems_list := util.SystemsListFromDataFiles( util.DataFiles() )
    formatted_systems_list := util.FormatSystemsOutput( systems_list )
    util.WriteJsonOutput( w, formatted_systems_list )
}

// Given a system, outputs all games for it.
func ListSystemGames(
    w http.ResponseWriter,
    r *http.Request,
) {
    system := chi.URLParam( r, "system" )
    data_files := util.DataFiles()
    system_data_file, does_system_exist := data_files[ system ]
    if ! does_system_exist {
        http.Error( w, "System not found", http.StatusNotFound )
        log.Println( "System does not exist: " + system )
        return
    }

    games := util.FetchGamesFromRedis( system )

    if len(games) == 0 {
        log.Printf( "Did not fetch system from cache, getting from file" )

        games, err := util.GetGamesData( system_data_file )
        if err != nil {
            http.Error( w, "Internal error", http.StatusInternalServerError )
            log.Println( err )
            return
        }

        err = util.SetGamesOnRedis( games, system )
        if err != nil {
            log.Printf( "Error setting redis key: %v\n", err )
            // Can still continue since we have the data from outside 
            // the cache
        }
    } else {
        log.Printf( "Fetched system from cache" )
    }

    output := util.FormatGamesList( games )
    util.WriteJsonOutput( w, output )
}
