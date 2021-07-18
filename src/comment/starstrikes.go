package comment

import "time"

type Category int

const (
	Star Category = iota
	MinorStrike
	Strik
	MajorStrike
)

type StarStrike struct {
	ID int `json:"id,omitempty"`
	PermID string `json:"perm_id,omitempty"`
	// staff(name)
	Teacher string `json:"teacher,omitempty"`
	Comment string `json:"comment,omitempty"`
	// Title is a catagory of the comment
	Title string `json:"title,omitempty"`

	Created time.Time `json:"created,omitempty"`
	// 0 star 1 minor 2 strik 3 major
	Cat      Category
	IsActive bool `json:"is_active,omitempty"`
}
