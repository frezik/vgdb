package main

import (
    "net/http"

    "github.com/frezik/vgdb/router"
)


func main() {
    http.ListenAndServe( ":3000", router.Routes() )
}
