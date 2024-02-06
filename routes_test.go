package main

import (
    "encoding/json"
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

func TestSystemsEndpoint( t *testing.T ) {
    req, err := http.NewRequest( "GET", "/systems", nil )

    if err != nil {
        t.Errorf( "Error creating new request: %v", err )
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc( router.ListSystems )
    handler.ServeHTTP( rr, req )

    if status := rr.Code; status != http.StatusOK {
        t.Errorf( "Handler returned wrong status code. Expected: %d, got: %d",
            http.StatusOK, status )
    }

    systems := make( map[string][]string )
    if err := json.NewDecoder( rr.Body ).Decode( &systems ); err != nil {
        t.Errorf( "Error decoding response body: %v", err )
    }
}
