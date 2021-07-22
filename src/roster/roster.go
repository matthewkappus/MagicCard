package roster

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/matthewkappus/MagicCard/src/db"
	"github.com/matthewkappus/Roster/src/synergy"
)

type StaffView struct {
	tmpls *template.Template
	store *db.Store
	// store sql.DB
}

// ClassList return unique CourseIDAndTitle sections for "teacher" cookie
// Returns error if no teacher cookie
func (sv *StaffView) ClassList(r *http.Request) ([]*synergy.Stu415, error) {
	teacherCookie, err := r.Cookie("teacher")

	if err != nil {
		return nil, err
	}
	return sv.store.ListClasses(teacherCookie.Value)

}

// NewView takes a roster db and tmpl path and returns handler object
func NewView(store *db.Store) (*StaffView, error) {
	tmpls, err := template.ParseGlob("tmpl/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	return &StaffView{
		tmpls: tmpls,
		store: store,
	}, nil

}

func (sv *StaffView) Search(w http.ResponseWriter, r *http.Request) {
	teacher := sv.GetTeacher(r)
	nav, err := sv.MakeNav(teacher, "students", "Student Search")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cr, err := sv.MakeSchoolClassroom(teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sv.tmpls.Lookup("search").Execute(w, TD{N: nav, C: cr})
}

func (sv *StaffView) Home(w http.ResponseWriter, r *http.Request) {

	teacher := sv.GetTeacher(r)
	nav, err := sv.MakeNav(teacher, "home", "Magic Card")

	// show sign in if no teacher cookie
	if err != nil || teacher == "" {
		nav.Title = "Sign In To Magic Card"
	} else {
		nav.Title = "Magic Card"
	}
	sv.tmpls.Lookup("home").Execute(w, TD{N: nav})
}

func UpdateRoster() ([]*synergy.Stu415, error) {
	// todo: get from synergy
	f, err := os.Open("data/stu415.csv")
	if err != nil {
		log.Fatal(err)
	}

	return synergy.ReadStu415sFromCSV(f)

}
