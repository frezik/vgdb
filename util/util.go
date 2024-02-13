package util

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"
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
