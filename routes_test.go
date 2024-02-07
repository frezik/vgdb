package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    router "github.com/frezik/vgdb/router"
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

func TestGamesEndpoint( t *testing.T ) {
    chi_router := router.Routes()
    ts := httptest.NewServer( chi_router )
    defer ts.Close()

    req, err := http.NewRequest( "GET", ts.URL + "/snes/games", nil )
    if err != nil {
        t.Errorf( "Error creating new request: %v", err )
    }

    response, err := http.DefaultClient.Do( req )
    if err != nil {
        t.Fatal( err )
    }

    if status := response.StatusCode; status != http.StatusOK {
        t.Errorf( "Handler returned wrong status code. Expected: %d, got: %d",
            http.StatusOK, status )
    }

    games := make( map[string][]string )
    if err := json.NewDecoder( response.Body ).Decode( &games ); err != nil {
        t.Errorf( "Error decoding response body: %v", err )
    }

    first_game := games[ "games" ][0]
    if first_game != "'89 Dennou Kyuusei Uranai" {
        t.Errorf( "Expected different first game, got: " + first_game )
    }
}
