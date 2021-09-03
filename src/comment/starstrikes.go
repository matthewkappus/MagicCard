package comment

import (
	"time"

	"github.com/matthewkappus/Roster/src/synergy"
)

type Category int

const (
	Star Category = iota
	MinorStrike
	Strike
	MajorStrike
)

type StarStrike struct {
	ID     int    `json:"id,omitempty"`
	PermID string `json:"perm_id,omitempty"`
	// staff(name)
	Teacher string `json:"teacher,omitempty"`
	Comment string `json:"comment,omitempty"`
	// Title is a catagory of the comment
	Title string `json:"title,omitempty"`

	Created time.Time `json:"created,omitempty"`
	// 0 star 1 minor 2 strik 3 major
	Cat Category

	Icon     string
	IsActive bool `json:"is_active,omitempty"`
}

// Contact logs teacher-parent communications regarding starstrikes
type Contact struct {
	ID int

	Sender *synergy.Staff

	StudentName string

	// student id
	PermID string

	// when contacted
	Sent time.Time

	// Respondent usually is parent of student
	Respondent string

	StarStrike *StarStrike

	Message string
}

func NewContact(ss *StarStrike, sender *synergy.Staff, resp, msg string) *Contact {
	return &Contact{
		StarStrike: ss,
		Sender:     sender,
		Respondent: resp,
		Message:    msg,
	}
}
