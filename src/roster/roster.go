package roster

import (
	"fmt"
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
	nav, err := sv.MakeNav(teacher, "students", "Student Search", Teacher)

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

	scope, teacher, student, err := sv.GetSessionUser(r)
	if err != nil {
		scope = -1
	}

	nav := &Nav{Title: "Sign In To Magic Card"}
	switch scope {
	case Teacher:
		nav, err = sv.MakeNav(teacher, "home", "Magic Card", scope)
	case Admin:
		nav, err = sv.MakeNav(teacher, "home", "Magic Card", scope)
	case Student:
		nav, err = sv.MakeNav(student, "home", "Magic Card", scope)
	default:
		fmt.Println("no scope found")
	}

	if err != nil {
		fmt.Printf("error making nav: %v", err)
	}
	sv.tmpls.Lookup("home").Execute(w, TD{N: nav})

}

// GetSessionType returns Admin, Teacher, or Student scope if exists or -1 for no session scope
func (sv *StaffView) GetSessionType(r *http.Request) Scope {
	scopeCookie, err := r.Cookie("scope")
	if err != nil {
		return -1
	}

	switch scopeCookie.Value {
	case "0":
		return Student
	case "1":
		return Teacher
	case "2":
		return Admin
	default:
		return -1
	}
}

// GetSessionUser returns scope and the "student" perm or "teacher" name
func (sv *StaffView) GetSessionUser(r *http.Request) (s Scope, teacher, student string, err error) {
	s = sv.GetSessionType(r)
	if s == -1 {
		return -1, "", "", err
	}
	teacher = sv.GetTeacher(r)
	if teacher == "" {
		student = sv.GetStudent(r)
	}
	return s, teacher, student, nil

}
func UpdateRoster() ([]*synergy.Stu415, error) {
	// todo: get from synergy
	f, err := os.Open("data/stu415.csv")
	if err != nil {
		log.Fatal(err)
	}

	return synergy.ReadStu415sFromCSV(f)

}
