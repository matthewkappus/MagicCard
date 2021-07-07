package main

import (
	"log"
	"net/http"

	"github.com/matthewkappus/MagicCard/src/roster"
)

func main() {
	// stus, err := synergy.GetStu401s("e204920", "AZS2209h", archiveStude401s, time.Minute*2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	sv, err := roster.NewView()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/search", sv.Search)
	http.HandleFunc("/addComment", sv.Add)
	http.HandleFunc("/card", sv.Card)
	http.HandleFunc("/", sv.Home)

	http.ListenAndServe(":8080", nil)
}
