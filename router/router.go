package router

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/redis/go-redis/v9"

    "github.com/frezik/vgdb/reduce"
)

const REDIS_ADDR = "localhost:6379"
const REDIS_PASSWD = ""
const REDIS_DB = 0


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

var redis_client = redis.NewClient( &redis.Options{
    Addr: REDIS_ADDR,
    Password: REDIS_PASSWD,
    DB: REDIS_DB,
})




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
    write_func := reduce.WriteJsonOutput( w, reduce.EndFunc() )
    encode_func := reduce.EncodeSystemsOutput( write_func )
    systems_list_func := reduce.SystemsListFromDataFiles(
        data_files,
        encode_func,
    )
    
    systems_list_func()
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

    ctx := context.Background()
    redis_key := "system-games-" + system + ":1"
    games := redis_client.LRange( ctx, redis_key, 0, -1 ).Val()

    if len(games) == 0 {
        log.Printf( "Did not fetch system from cache, getting from file" )

        system_data, err := GetSystemData( system_data_file )
        if err != nil {
            http.Error( w, "Internal error", http.StatusInternalServerError )
            log.Println( err )
            return
        }

        games = make( []string, len( system_data ) )
        for k := range system_data {
            games[k] = system_data[k].Name
        }

        err = redis_client.RPush( ctx, redis_key, games ).Err()
        if err != nil {
            log.Printf( "Error setting redis key: %v\n", err )
            // Can still continue since we have the data from outside 
            // the cache
        }
    } else {
        log.Printf( "Fetched system from cache" )
    }

    output := map[string][]string {
        "games": games,
    }

    err := json.NewEncoder( w ).Encode( output )
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
