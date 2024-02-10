package util

import (
    "encoding/json"
    "log"
    "net/http"
)


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
