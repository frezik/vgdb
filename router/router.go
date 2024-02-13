package router

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "github.com/frezik/vgdb/util"
)


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
    systems_list := util.SystemsListFromDataFiles( util.DataFiles() )
    formatted_systems_list := util.FormatSystemsOutput( systems_list )
    util.WriteJsonOutput( w, formatted_systems_list )
}

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

        system_data, err := util.GetSystemData( system_data_file )
        if err != nil {
            http.Error( w, "Internal error", http.StatusInternalServerError )
            log.Println( err )
            return
        }

        games = make( []string, len( system_data ) )
        for k := range system_data {
            games[k] = system_data[k].Name
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
