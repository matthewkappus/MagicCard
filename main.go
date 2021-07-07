package main

import (
	"fmt"
	"log"

	"github.com/matthewkappus/MagicCard/src/db"
)

func main() {

	store, err := db.OpenStore("data/cards.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = store.UpdateStu415("data/stu415.csv"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Store loaded & ready")
	// stus, err := synergy.GetStu401s("e204920", "AZS2209h", archiveStude401s, time.Minute*2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// sv, err := roster.NewView()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// http.HandleFunc("/search", sv.Search)
	// http.HandleFunc("/addComment", sv.Add)
	// http.HandleFunc("/card", sv.Card)
	// // login
	// http.HandleFunc("/", sv.Home)

	// http.ListenAndServe(":8080", nil)
}
