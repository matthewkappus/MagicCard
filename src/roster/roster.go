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
	User  string
	Type  Scope
	tmpls *template.Template
	store *db.Store
	M     *MagicCard
	C     *Classroom
	N     *Nav
	// store sql.DB
}

// NewView takes a roster db and tmpl path and returns a View that
// scopes requests to users with that viewType
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

// HF registers handler to provided path and provides handler
// with  Nav and authentication
func (v *View) HF(path string, h http.HandlerFunc) {

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		userScope, user, err := v.GetSessionUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if v.Type != userScope {
			fmt.Printf("HF: v.Type %v != cookie.Type %v\n", v.Type, userScope)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		v.User = user
		// todo: normalize path to <title>
		v.N, err = v.MakeNav(user, path, path, v.Type, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h(w, r)
	})

}

// ClassList return unique CourseIDAndTitle sections for "teacher" cookie
// Returns error if no teacher cookie
func (v *View) ClassList(r *http.Request) ([]*synergy.Stu415, error) {
	teacherCookie, err := r.Cookie("teacher")

	if err != nil {
		return nil, err
	}
	return v.store.ListClasses(teacherCookie.Value)

}

func (v *View) Search(w http.ResponseWriter, r *http.Request) {

	cr, err := v.MakeSchoolClassroom(v.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v.tmpls.Lookup("search").Execute(w, TD{N: v.N, C: cr})
}

func (v *View) Home(w http.ResponseWriter, r *http.Request) {
	scope, user, err := v.GetSessionUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v.N, err = v.MakeNav(user, "/", "home", scope, w, r)
	switch scope {
	case Teacher:
		v.M, err = v.MakeTeacherMagicCard(user)
	case Admin:
		// show admin log with all classes
	case Student:
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		v.M, err = v.MakeStudentMagicCard(user)
	default:
		fmt.Println("guest view")
	}

	if err != nil {
		fmt.Printf("error making nav: %v", err)
	}
	v.tmpls.Lookup("home").Execute(w, TD{N: v.N, M: v.M})

}

// GetSessionType returns  3 Admin, 2 Teacher, 1 Student or 0 for guest scope (not signed in)
func (v *View) GetSessionType(r *http.Request) Scope {
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
func (v *View) GetSessionUser(r *http.Request) (s Scope, user string, err error) {
	s = v.GetSessionType(r)
	if s == 0 {
		return 0, "", fmt.Errorf("not signed in")
	}

	userCookie, err := r.Cookie("user")
	if err != nil {
		return 0, "", err
	}

	return s, string(userCookie.Value), nil

}
func UpdateRoster() ([]*synergy.Stu415, error) {
	// todo: get from synergy
	f, err := os.Open("data/stu415.csv")
	if err != nil {
		log.Fatal(err)
	}

	return synergy.ReadStu415sFromCSV(f)

}

func (v *View) SendAlert(w http.ResponseWriter, a *Alert) {

	http.SetCookie(w, &http.Cookie{
		Name:  "alert_msg",
		Value: a.Message,
		Path:  "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "alert_type",
		Value: a.Type,
		Path:  "/",
	})
}

func (v *View) ReadAlert(w http.ResponseWriter, r *http.Request) (*Alert, error) {
	alertCookie, err := r.Cookie("alert_msg")
	if err != nil {
		return nil, err
	}
	if alertCookie.Value == "" {
		return nil, nil
	}
	alertTypeCookie, err := r.Cookie("alert_type")
	if err != nil {
		return nil, err
	}
	a := &Alert{
		Message: alertCookie.Value,
		Type:    alertTypeCookie.Value,
	}

	// expire  cooReadAlertkies
	alertCookie.MaxAge = -1
	alertTypeCookie.MaxAge = -1
	http.SetCookie(w, alertCookie)
	http.SetCookie(w, alertTypeCookie)
	return a, nil

}
