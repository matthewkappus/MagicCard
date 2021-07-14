package roster

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// type ClaimIdentity struct {
// 	Aud           string `json:"aud"`
// 	Azp           string `json:"azp"`
// 	Email         string `json:"email"`
// 	EmailVerified bool   `json:"email_verified"`
// 	Exp           int    `json:"exp"`
// 	FamilyName    string `json:"family_name"`
// 	GivenName     string `json:"given_name"`
// 	Hd            string `json:"hd"`
// 	Iat           int    `json:"iat"`
// 	Iss           string `json:"iss"`
// 	Jti           string `json:"jti"`
// 	Name          string `json:"name"`
// 	Nbf           int    `json:"nbf"`
// 	Picture       string `json:"picture"`
// 	Sub           string `json:"sub"`
// }

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

	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(credential, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	cs, _ := token.Claims.(jwt.MapClaims)

	email, ok := cs["email"].(string)
	if !ok {
		return "", fmt.Errorf("could not get email from JWT claim")
	}
	return strings.ToLower(email), nil

}
func (sv *StaffView) TeacherLock(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sv.isTeacher(r) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

func (sv *StaffView) GetTeacher(r *http.Request) string {
	teacherCookie, err := r.Cookie("teacher")
	if err != nil {
		return ""
	}

	return string(teacherCookie.Value)
}

// looks for teacher and key cookie, matches with staff db
func (sv *StaffView) isTeacher(r *http.Request) bool {

	keyCookie, err := r.Cookie("key")
	if err != nil {
		return false
	}

	dbKey, err := sv.store.GetKeyByTeacher(sv.GetTeacher(r))
	if err != nil {
		fmt.Printf("error looking up key %v", err)
		return false
	}

	return keyCookie.Value == dbKey
}

// Login takes JWT from Google Sign In Button and sets the name, email and token values
func (sv *StaffView) Login(w http.ResponseWriter, r *http.Request) {

	email, err := ParseIdentity(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// todo: check if student number
	teacher, name, key, err := sv.store.TeacherNameFromEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "name", Value: name, Domain: "/"})
	http.SetCookie(w, &http.Cookie{Name: "teacher", Value: teacher})
	http.SetCookie(w, &http.Cookie{Name: "key", Value: key})

	http.Redirect(w, r, "/classes", http.StatusTemporaryRedirect)

}

func Matty(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "name", Value: "Matt Kappus"})
	http.SetCookie(w, &http.Cookie{Name: "teacher", Value: "Kappus, Matthew D."})
	http.SetCookie(w, &http.Cookie{Name: "key", Value: "1830a69c-a641-4832-9b38-77320de25756"})

	http.Redirect(w, r, "/classes", http.StatusTemporaryRedirect)
}
