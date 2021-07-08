package roster

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"text/template"

	"github.com/matthewkappus/MagicCard/src/db"
	"github.com/matthewkappus/Roster/src/synergy"
	"golang.org/x/net/publicsuffix"
)

type StudentView struct {
	jar   *cookiejar.Jar
	tmpls *template.Template
	store *db.Store
	// store sql.DB
}

// New view takes a roster db and tmpl path and returns handler object
// todo: put auth of object for sessions
func NewView(store *db.Store) (*StudentView, error) {
	tmpls, err := template.ParseGlob("tmpl/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	// todo: load roster from sql

	return &StudentView{
		jar:   jar,
		tmpls: tmpls,
		store: store,
	}, nil

}

func (sv *StudentView) Search(w http.ResponseWriter, r *http.Request) {
	stu401s, err := UpdateRoster()
	if err != nil {
		http.Error(w, "Could not get students\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	sv.tmpls.Lookup("search.tmpl.html").Execute(w, stu401s)
}

func (sv *StudentView) Home(w http.ResponseWriter, r *http.Request) {

}

func (sv *StudentView) Card(w http.ResponseWriter, r *http.Request) {
	// todo: look up student in db: join with comments
	id := r.FormValue("id")
	if len(id) != 9 {
		http.NotFound(w, r)
		return
	}

	stu415s, _ := UpdateRoster()
	var found *synergy.Stu415

	for _, stu := range stu415s {
		if stu.PermID == id {
			found = stu
			break
		}

	}

	// todo: add email and session info
	sv.tmpls.Lookup("card.tmpl.html").Execute(w, found)
}
func UpdateRoster() ([]*synergy.Stu415, error) {
	// todo: get from synergy
	f, err := os.Open("data/stu415.csv")
	if err != nil {
		log.Fatal(err)
	}

	return synergy.ReadStu415sFromCSV(f)

}
