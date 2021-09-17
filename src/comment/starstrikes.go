package comment

import (
	"time"

	"github.com/matthewkappus/Roster/src/synergy"
)

type StarStrike struct {
	ID     int    `json:"id,omitempty"`
	PermID string `json:"perm_id,omitempty"`

	Student *synergy.Stu415

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

type Category int

const (
	Star Category = iota
	MinorStrike
	MajorStrike
)

func ToCat(s string) Category {
	switch s {
	case "0":
		return Star
	case "1":
		return MinorStrike
	case "2":
		return MajorStrike
	}
	return Star
}

func (ss *StarStrike) IsValid() bool {

	// todo: validate the ss
	return true
}

// Contact logs teacher-parent communications regarding starstrikes
type Contact struct {
	ID int

	Sender *synergy.Staff

	// when contacted
	Sent time.Time

	// Respondent usually is parent of student
	Respondent string

	StarStrike *StarStrike

	Message string

	// Is issue resolved
	IsClosed bool
}

func NewContact(ss *StarStrike, sender *synergy.Staff, resp, msg string, isClosed bool) *Contact {
	return &Contact{
		StarStrike: ss,
		Sender:     sender,
		Respondent: resp,
		Message:    msg,
		IsClosed:   isClosed,
	}
}

// BatchCreateStarstrike returns a slice of non-repeating starstrikes for the provided students
func BatchCreateStarstrike(s *StarStrike, permIDs []string) []*StarStrike {

	m := make(map[string]*StarStrike)
	for _, id := range permIDs {
		s.PermID = id
		m[id] = s
	}
	out := make([]*StarStrike, len(m))
	for _, ss := range m {
		out = append(out, ss)
	}
	return out

}
