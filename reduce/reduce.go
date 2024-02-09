package reduce

import (
    "encoding/json"
    "log"
    "net/http"
)

func Exec() {
}


func SystemsListFromDataFiles(
    data_files map[string]string,
    next_func func( []string ),
) (func()) {
    systems := make( []string, len( data_files ) )

    i := 0
    for k := range data_files {
        systems[i] = k
        i++
    }

    return func() {
        next_func( systems )
    }

}

func EncodeSystemsOutput(
    next_func func( map[string][]string ),
) (func( []string )) {
    return func( systems []string ) {
        output := map[string][]string {
            "systems": systems,
        }
        next_func( output )
    }
}

func WriteJsonOutput(
    w http.ResponseWriter,
    next_func func(),
) (func( output map[string][]string )){
    return func( output map[string][]string ) {
        err := json.NewEncoder( w ).Encode( output )
        if err != nil {
            http.Error( w, "Internal error", http.StatusInternalServerError )
            log.Println( err )
            return
        }
    }
}

func EndFunc() (func ()) {
    return func () {}
}
