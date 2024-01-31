package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/frezik/vgdb/router"
)

func TestHealthEndpoint( t *testing.T ) {
    req, err := http.NewRequest( "GET", "/", nil )

    if err != nil {
        t.Errorf( "Error creating new request: %v", err )
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc( router.HeartBeat )
    handler.ServeHTTP( rr, req )

    if status := rr.Code; status != http.StatusOK {
        t.Errorf( "Handler returned wrong status code. Expected: %d, got: %d",
            http.StatusOK, status )
    }
}
