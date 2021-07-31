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

type View struct {
	// student/admin/teacher
	UserCookie string
	Type       Scope
	Nav        *Nav
	tmpls      *template.Template
	M          *MagicCard
	C          *Classroom
	N          *Nav
	store      *db.Store
	// store sql.DB
}

// ClassList return unique CourseIDAndTitle sections for "teacher" cookie
// Returns error if no teacher cookie
func (sv *View) ClassList(r *http.Request) ([]*synergy.Stu415, error) {
	teacherCookie, err := r.Cookie("teacher")

	if err != nil {
		return nil, err
	}
	return sv.store.ListClasses(teacherCookie.Value)

}

// NewView takes a roster db and tmpl path and returns handler object
func NewView(store *db.Store, templateGlob string, viewType Scope) (*View, error) {
	tmpls, err := template.ParseGlob(templateGlob)
	if err != nil {
		return nil, err
	}

	return &View{
		Type:  viewType,
		tmpls: tmpls,
		store: store,
	}, nil

}

func (sv *View) Search(w http.ResponseWriter, r *http.Request) {
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

func (sv *View) Home(w http.ResponseWriter, r *http.Request) {

	scope, teacher, student, err := sv.GetSessionUser(r)
	if err != nil {
		// set to guest scope in not signed in
		scope = 0
	}

	nav := &Nav{Title: "Sign In To Magic Card"}
	mc := new(MagicCard)
	switch scope {
	case Teacher:
		nav, err = sv.MakeNav(teacher, "home", "Magic Card", scope)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mc, err = sv.MakeTeacherMagicCard(teacher)

	case Admin:
		nav, err = sv.MakeNav(teacher, "home", "Magic Card", scope)
	case Student:
		nav, err = sv.MakeNav(student, "home", "Magic Card", scope)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mc, err = sv.MakeStudentMagicCard(student)
	case Guest:
		nav, err = sv.MakeNav("guest", "home", "Magic Card", scope)
	default:
		fmt.Println("no scope found")
	}

	if err != nil {
		fmt.Printf("error making nav: %v", err)
	}
	sv.tmpls.Lookup("home").Execute(w, TD{N: nav, M: mc})

}

// GetSessionType returns  3 Admin, 2 Teacher, 1 Student or 0 for guest scope (not signed in)
func (sv *View) GetSessionType(r *http.Request) Scope {
	scopeCookie, err := r.Cookie("scope")
	if err != nil {
		return 0
	}

	switch scopeCookie.Value {
	case "0":
		return Guest
	case "1":
		return Student
	case "2":
		return Teacher
	case "3":
		return Admin
	default:
		return 0
	}
}

// GetSessionUser returns scope and the "student" perm or "teacher" name
func (sv *View) GetSessionUser(r *http.Request) (s Scope, teacher, student string, err error) {
	s = sv.GetSessionType(r)
	if s == 0 {
		return 0, "", "", fmt.Errorf("not signed in")
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
