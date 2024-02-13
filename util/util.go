package util

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/redis/go-redis/v9"
)


const data_dir = "data"
var data_files = map[string]string{
    "nes": "nes.json",
    "snes": "snes.json",
    "sega_genesis": "sega_genesis.json",
    "sega_master": "sega_master.json",
}


const REDIS_ADDR = "localhost:6379"
const REDIS_PASSWD = ""
const REDIS_DB = 0

var redis_client = redis.NewClient( &redis.Options{
    Addr: REDIS_ADDR,
    Password: REDIS_PASSWD,
    DB: REDIS_DB,
})


type SystemData struct {
    Name string
    Publisher string
    FirstRelease string
    JPRelease string
    NARelease string
    EURelease string
    BRRelease string
}


func DataFiles() map[string]string {
    return data_files
}

func SystemsListFromDataFiles(
    data_files map[string]string,
) []string {
    systems := make( []string, len( data_files ) )

    i := 0
    for k := range data_files {
        systems[i] = k
        i++
    }

    return systems
}

func FormatSystemsOutput(
    systems []string,
) map[string][]string {
    output := map[string][]string {
        "systems": systems,
    }
    return output
}

func FormatGamesList(
    games []string,
) map[string][]string {
    output := map[string][]string {
        "games": games,
    }
    return output
}

func WriteJsonOutput(
    w http.ResponseWriter,
    output map[string][]string,
) {
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

func FetchGamesFromRedis(
    system string,
) []string {
    ctx := context.Background()
    redis_key := "system-games-" + system + ":1"
    games := redis_client.LRange( ctx, redis_key, 0, -1 ).Val()
    return games
}

func SetGamesOnRedis(
    games []string,
    system string,
) error {
    ctx := context.Background()
    redis_key := "system-games-" + system + ":1"
    err := redis_client.RPush( ctx, redis_key, games ).Err()
    return err
}
