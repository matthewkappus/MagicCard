package roster

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

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
		v.StartSession(stu.PermID, formatName(stu.StudentName), Student, w, r)
	} else {
		// teacher format
		teacher, name, err := v.store.TeacherNameFromEmail(email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		v.StartSession(teacher, name, Teacher, w, r)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func (v *View) DevAdminLogin(w http.ResponseWriter, r *http.Request) {
	v.StartSession("Madison Admin", "Madison Admin", Admin, w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (v *View) DevTeacherLogin(w http.ResponseWriter, r *http.Request) {

	v.StartSession("Kappus, Matthew D.", "Matt Kappus", Teacher, w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (v *View) DevStudentLogin(w http.ResponseWriter, r *http.Request) {
	v.StartSession("980016917", "Abbas, Malak", Student, w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
