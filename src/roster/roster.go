package roster

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/matthewkappus/MagicCard/src/db"
	"github.com/matthewkappus/Roster/src/synergy"
)

type View struct {
	// 	Sessions map[SID]*Session
	// teacher name or student perm
	User string
	// F L-formatted name
	Name  string
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

// sessionCookies returns name, user, sid and scope of cookies
func sessionCookies(r *http.Request) (name string, user string, sid string, scope Scope, err error) {
	// return string(r.Cookie("name").Value), string(r.Cookie("sid").Value), getSessionType(r)
	nameCookie, err := r.Cookie("name")
	if err != nil {
		return "", "", "", 0, err
	}
	name = string(nameCookie.Value)

	userCookie, err := r.Cookie("user")
	if err != nil {
		return "", "", "", 0, err
	}
	user = string(userCookie.Value)

	sidCookie, err := r.Cookie("sid")
	if err != nil {
		return "", "", "", 0, err
	}
	sid = string(sidCookie.Value)

	scopeCookie, err := r.Cookie("scope")
	if err != nil {
		return "", "", "", 0, err
	}

	switch scopeCookie.Value {
	case "0":
		scope = Guest
	case "1":
		scope = Student
	case "2":
		scope = Teacher
	case "3":
		scope = Admin
	default:
		scope = Guest
	}
	return name, user, sid, scope, nil
}

// HF registers handler to provided path and provides handler
// with  Nav and authentication
func (v *View) HF(path string, h http.HandlerFunc) {

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		name, user, sid, scope, _ := sessionCookies(r)
		// todo: What happens if err != nil? (lacks user cookies)

		v.User = user
		sesh, err := v.GetSession(sid)
		if err != nil {
			// todo; create new session
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if sesh.Expires.Before(time.Now()) {
			// todo; create new session
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		if sesh.User != user {
			http.Error(w, "User not in a session "+user, http.StatusUnauthorized)
			return
		}

		// check if user has enough scope
		if scope < v.Type {
			fmt.Printf("HF: v.Type %v != cookie.Type %v\n", v.Type, scope)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// 0 guest < 1 student < 2 teacher < 3 admin

		v.N, err = v.MakeNav(user, path, pathToTitle(path), name, scope, w, r)
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
	name, user, _, scope, _ := sessionCookies(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	var err error
	v.N, err = v.MakeNav(user, "/", "home", name, scope, w, r)
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
