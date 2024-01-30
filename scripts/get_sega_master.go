package main

import (
    "fmt"
    "github.com/gocolly/colly/v2"
    "encoding/json"
)

type Game struct {
    Name string
    Publisher string
    FirstRelease string
    JPRelease string
    NARelease string
    EURelease string
    BRRelease string
}

func main() {
    c := colly.NewCollector(
        colly.AllowedDomains( "en.wikipedia.org" ),
    )

    var results []Game
    c.OnHTML( "table.wikitable:nth-of-type(1) tbody tr", func( e *colly.HTMLElement ) {
        var row_texts []string = e.ChildTexts( "td" )
        if len( row_texts ) == 0 {
            return 
        }

        game := Game{
            Name: row_texts[0],
            Publisher: row_texts[1],
            FirstRelease: row_texts[2],
            JPRelease: row_texts[3],
            NARelease: row_texts[4],
            EURelease: row_texts[5],
            BRRelease: row_texts[6],
        }
        results = append( results, game )
    })

    c.OnScraped( func( r *colly.Response ) {
        b, err := json.Marshal( results )
        if err != nil {
            fmt.Println( "Error encoding to JSON:", err )
        } else {
            fmt.Println( string(b) )
        }
    })

    c.Visit( "https://en.wikipedia.org/wiki/List_of_Master_System_games" )
}
