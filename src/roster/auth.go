package roster

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Scopes: 0 Guest 1 Student 2 Teacher 3 Admin
type Scope int

const (
	Guest Scope = iota
	Student
	Teacher
	Admin
)

type Session struct {
	// student perm or teacher name
	User    string
	Expires time.Time
	Scope   Scope
}

func (v *View) StartSession(id string, name string, s Scope, w http.ResponseWriter) {
	sid := uuid.NewString()
	exp := time.Now().Add(8 * time.Hour)
	var scope string
	switch s {
	case Guest:
		scope = "0"
	case Student:
		scope = "1"
	case Teacher:
		scope = "2"
	case Admin:
		scope = "3"
	}

	fmt.Println("starting session with sid", sid)
	Sessions[sid] = &Session{
		User:    id,
		Expires: exp,
		Scope:   s,
	}

	http.SetCookie(w, &http.Cookie{Name: "sid", Value: sid, Expires: exp})
	http.SetCookie(w, &http.Cookie{Name: "name", Value: name, Expires: exp})
	http.SetCookie(w, &http.Cookie{Name: "user", Value: id, Expires: exp})
	http.SetCookie(w, &http.Cookie{Name: "scope", Value: scope, Expires: exp})

}

// todo: refresh session by removing expired

func ParseIdentity(r *http.Request) (email string, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	u, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}
	// content is the JWT containing the claims, including email
	credential := u.Get("credential")

	// todo: change jwt library
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(credential, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("D0UBel*V"), nil
	})

	cs, _ := token.Claims.(jwt.MapClaims)

	email, ok := cs["email"].(string)
	if !ok {
		return "", fmt.Errorf("could not get email from JWT claim")
	}
	return strings.ToLower(email), nil

}

// Login takes JWT from Google Sign In Button and sets the name, email and token values
// Login sets the view Scope and session Identity on view and User cookies
func (v *View) Login(w http.ResponseWriter, r *http.Request) {

	email, err := ParseIdentity(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if studentFormat, _ := regexp.MatchString("[0-9]+@aps.edu", email); studentFormat {
		stu, err := v.store.StudentFromEmail(email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		v.StartSession(stu.PermID, formatName(stu.StudentName), Student, w)
	} else {
		// teacher format
		teacher, name, err := v.store.TeacherNameFromEmail(email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		v.StartSession(teacher, name, Teacher, w)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func (v *View) DevAdminLogin(w http.ResponseWriter, r *http.Request) {
	v.StartSession("Madison Admin", "Madison Admin", Admin, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (v *View) DevTeacherLogin(w http.ResponseWriter, r *http.Request) {

	v.StartSession("Kappus, Matthew D.", "Matt Kappus", Teacher, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (v *View) DevStudentLogin(w http.ResponseWriter, r *http.Request) {
	v.StartSession("980016917", "Abbas, Malak", Student, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
