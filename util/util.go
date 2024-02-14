// Package util provides most of the heavy lifting of getting and modifying data
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


// Directory for holding the JSON files for each system. This is relative 
// to the top level dir of the project. The files in data_files will be in 
// this dir.
const data_dir = "data"
// Maps short system names (like "nes") to a JSON file that has the data for 
// that system. The files are under data_dir.
var data_files = map[string]string{
    "nes": "nes.json",
    "snes": "snes.json",
    "sega_genesis": "sega_genesis.json",
    "sega_master": "sega_master.json",
}


// Redis addr
const REDIS_ADDR = "localhost:6379"
// Redis pass
const REDIS_PASSWD = ""
// Redis DB num
const REDIS_DB = 0

// Connection to redis
var redis_client = redis.NewClient( &redis.Options{
    Addr: REDIS_ADDR,
    Password: REDIS_PASSWD,
    DB: REDIS_DB,
})


// SystemData holds information about a single game.
//
// All dates are in "Monthname YYYY" format, and can be null.
type SystemData struct {
    Name string // Name of game
    Publisher string // Publisher
    FirstRelease string // Date of first release in any territory
    JPRelease string // Japan release date
    NARelease string // North America release date
    EURelease string // Europe release date
    BRRelease string // Brazil release date
}


// Returns map of system names (like "nes") to data files in the data_dir
// directory.
func DataFiles() map[string]string {
    return data_files
}

// Transform the mapping of data files into a list of systems.
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

// Wraps a list of systems into a map with "output" as the single key holding 
// the list.
func FormatSystemsOutput(
    systems []string,
) map[string][]string {
    output := map[string][]string {
        "systems": systems,
    }
    return output
}

// Wraps a list of games into a map with "output" as the single key holding 
// the list.
func FormatGamesList(
    games []string,
) map[string][]string {
    output := map[string][]string {
        "games": games,
    }
    return output
}

// Transforms the give output into JSON and writes it out.
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

// From the given data_file path, returns all the games inside.
func GetGamesData(
    data_file string,
) ([]string, error) {
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

    games := make( []string, len( system_data ) )
    for k := range system_data {
        games[k] = system_data[k].Name
    }

    return games, nil
}

// Gets the list of games from Redis
func FetchGamesFromRedis(
    system string,
) []string {
    ctx := context.Background()
    redis_key := "system-games-" + system + ":1"
    games := redis_client.LRange( ctx, redis_key, 0, -1 ).Val()
    return games
}

// Sets the list of games to Redis
func SetGamesOnRedis(
    games []string,
    system string,
) error {
    ctx := context.Background()
    redis_key := "system-games-" + system + ":1"
    err := redis_client.RPush( ctx, redis_key, games ).Err()
    return err
}
