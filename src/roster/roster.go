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
	stu415s, err := sv.store.SelectStu415s()
	if err != nil {
		http.Error(w, "Could not get students\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	classes, _ := sv.ClassList(r)
	ci := &ClassInfo{
		Stu415s:   stu415s,
		ClassList: classes,
		Title:     "Student Search",
	}
	sv.tmpls.Lookup("search").Execute(w, ci)
}

func (sv *StaffView) Home(w http.ResponseWriter, r *http.Request) {

	classes, _ := sv.ClassList(r)
	teacher := sv.GetTeacher(r)
	ci := &ClassInfo{
		ClassList: classes,
		Teacher:   teacher,
		Title:     "Student Search",
		Path:      "home",
	}
	sv.tmpls.Lookup("home").Execute(w, ci)
}

func (sv *StaffView) Card(w http.ResponseWriter, r *http.Request) {
	// todo: look up student in db: join with comments
	permid := r.FormValue("id")
	if len(permid) != 9 {
		http.NotFound(w, r)
		return
	}

	stu415, err := sv.store.SelectStu415(permid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, _ := sv.ClassList(r)
	teacher := sv.GetTeacher(r)
	stars, bars, _ := sv.store.GetStarBars(teacher)

	ci := &ClassInfo{
		ClassList: classes,
		Stu415s:   []*synergy.Stu415{stu415},
		Stars:     stars,
		Bars:      bars,
		Teacher:   teacher,
		Title:     stu415.StudentName + " Magic Card",
		Path:      "classes",
	}
	// todo: add email and session info
	sv.tmpls.Lookup("card").Execute(w, ci)
}
func UpdateRoster() ([]*synergy.Stu415, error) {
	// todo: get from synergy
	f, err := os.Open("data/stu415.csv")
	if err != nil {
		log.Fatal(err)
	}

	return synergy.ReadStu415sFromCSV(f)

}
