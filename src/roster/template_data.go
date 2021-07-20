package roster

import (
	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

type CardData struct {
	Student *Student
	Class   *Classroom
}

type TD struct {
	S Student
	C Classroom
	N Nav
}
type Nav struct {
	ClassList []*synergy.Stu415
	Path      string
	Title     string
	Status    string
}

// Classroom is teachers class info
type Classroom struct {
	Stu415s   []*synergy.Stu415
	MyStars   []*comment.StarStrike
	MyStrikes []*comment.StarStrike
	ClassList []*synergy.Stu415
	Teacher   string
	ClassName string
}

// Student holds the Stu415 and Starbars for a template
//  My and Strikes are mapped to their title for counting
type Student struct {
	S415 *synergy.Stu415
	// Star.Title to SS
	StarMap map[string][]*comment.StarStrike
	// StrikeMap
	StrikeMap map[string][]*comment.StarStrike
}
