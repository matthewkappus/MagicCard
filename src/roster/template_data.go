package roster

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

type TD struct {
	M *MagicCard
	C *Classroom
	N *Nav
}

// Alert is a struct for use in <head> / <nav>
type Alert struct {
	// Types: primary secondary success warning danger info light dark
	Type    string
	Message string
	Link    string
}

type Nav struct {
	ClassList []*synergy.Stu415
	Path      string
	Title     string
	Status    string
	// permID or teacher(name)
	User string
	// 0 Guest 1 Student 2 Teacher 3 Admin
	Type  Scope
	Alert *Alert
}

// formatMame from L, F MI. to F L
func formatName(name string) string {
	n := strings.Split(name, " ")

	// If no  first & last name, return original
	if len(n) < 2 {
		return name
	}
	// return
	return fmt.Sprintf("%s %s", n[1], n[0][:len(n[0])-1])
}

// Classroom is teachers class info
type Classroom struct {
	// todo: remove stu415
	Stu415s []*synergy.Stu415
	// StarStrikes map students to StarStrike list for button data
	StarStrikes map[*synergy.Stu415][]*comment.StarStrike
	Teacher     string
	ClassName   string
}

// MagicCard holds the Stu415 and Starbars for a template
//  My and Strikes are mapped to their title for counting
type MagicCard struct {
	Name string
	ID   string

	S415 *synergy.Stu415
	// Star.Title to SS
	StarMap map[string][]*comment.StarStrike
	// StrikeMap
	StrikeMap map[string][]*comment.StarStrike
}

func (v *View) MakeClassroom(teacher, section string) (*Classroom, error) {
	s415s, err := v.store.ListStudents(section)
	if err != nil {
		return nil, err
	}
	// sortStudents by fomatting StudentName to F Mi L and ascending order by L
	sortStudents(s415s)

	myss, err := v.store.GetMyStarStrikes(teacher)

	if err != nil {
		return nil, err
	}

	// starstrikes takes generic starstrikes and puts each student perm in for use with buttons
	ss := make(map[*synergy.Stu415][]*comment.StarStrike)
	for _, stu := range s415s {
		s := make([]*comment.StarStrike, 0)
		for _, mss := range myss {
			mss.PermID = stu.PermID
			mss.Teacher = teacher
			s = append(s, mss)
		}
		ss[stu] = s
	}
	return &Classroom{
		Stu415s:     s415s,
		StarStrikes: ss,
		Teacher:     teacher,
	}, nil

}

// sortStudents by last name
// todo: fix Jrs:  John C. Jr Anderson
func sortStudents(s415s []*synergy.Stu415) {
	sort.Slice(s415s, func(i, j int) bool {
		return s415s[i].StudentName < s415s[j].StudentName
	})

	for _, c := range s415s {
		flm := strings.Split(c.StudentName, ", ")
		flm[1] = strings.Split(flm[1], " ")[0]
		if len(flm) > 1 {

			c.StudentName = fmt.Sprintf("%s %s", flm[1], flm[0])
		}
	}

}

// MakeSchoolClassroom returns list of every unique stu415
func (v *View) MakeSchoolClassroom(teacher string) (*Classroom, error) {
	s415s, err := v.store.ListAllStudents()
	if err != nil {
		return nil, err
	}

	myss, err := v.store.GetMyStarStrikes(teacher)
	if err != nil {
		return nil, err
	}

	// sortStudents by fomatting StudentName to F Mi L and ascending order by L
	sortStudents(s415s)

	// starstrikes takes generic starstrikes and puts each student perm in for use with buttons
	ss := make(map[*synergy.Stu415][]*comment.StarStrike)
	for _, stu := range s415s {
		s := make([]*comment.StarStrike, 0)
		for _, mss := range myss {
			mss.PermID = stu.PermID
			mss.Teacher = teacher
			s = append(s, mss)
		}
		ss[stu] = s
	}
	return &Classroom{
		Stu415s:     s415s,
		StarStrikes: ss,
		Teacher:     teacher,
	}, nil

}

func (v *View) MakeTeacherMagicCard(name string) (*MagicCard, error) {

	stars, strikes, err := v.store.SelectTeacherStarStrikes(name, 100)
	if err != nil {
		return nil, err
	}

	starsM := make(map[string][]*comment.StarStrike)
	for _, ss := range stars {
		match, exists := starsM[ss.Title]
		if exists {
			starsM[ss.Title] = append(match, ss)
			continue
		}
		starsM[ss.Title] = []*comment.StarStrike{ss}
	}

	stikesM := make(map[string][]*comment.StarStrike)
	for _, ss := range strikes {
		match, exists := stikesM[ss.Title]
		if exists {
			stikesM[ss.Title] = append(match, ss)
			continue
		}
		stikesM[ss.Title] = []*comment.StarStrike{ss}
	}

	return &MagicCard{

		Name:      name,
		ID:        name,
		StarMap:   starsM,
		StrikeMap: stikesM}, nil
}

func (v *View) MakeStudentMagicCard(perm string) (*MagicCard, error) {

	stu, err := v.store.SelectStu415(perm)
	if err != nil {
		return nil, err
	}

	ss, err := v.store.GetStarStrikesByPerm(perm)
	if err != nil {
		return nil, err
	}

	stars := make(map[string][]*comment.StarStrike)
	strikes := make(map[string][]*comment.StarStrike)

	for _, s := range ss {
		// 0 star 1 minor 2 strik 3 major
		if s.Cat == comment.Star {
			tStars, exists := stars[s.Title]
			if !exists {
				stars[s.Title] = []*comment.StarStrike{s}
				continue
			}
			// title exists: add to list
			stars[s.Title] = append(tStars, s)
		} else {
			// ss.Cat is 1+: A strike
			tStrikes, exists := strikes[s.Title]
			if !exists {
				strikes[s.Title] = []*comment.StarStrike{s}
				continue
			}
			// strike exists: add to rest
			strikes[s.Title] = append(tStrikes, s)
		}
	}

	return &MagicCard{

		Name:      formatName(stu.StudentName),
		ID:        stu.PermID,
		S415:      stu,
		StarMap:   stars,
		StrikeMap: strikes}, nil

}

func formatClassList(classlist []*synergy.Stu415) {
	for _, c := range classlist {
		idT := strings.Split(c.CourseIDAndTitle, "- ")
		c.CourseIDAndTitle = fmt.Sprintf("%s %s", c.Per, idT[1])
	}

	sort.Slice(classlist, func(i, j int) bool {
		return classlist[i].Per < classlist[j].Per
	})
}

// MakeNav returns data struct for use in <head> / <nav>
func (v *View) MakeNav(user, path, title string, stype Scope, w http.ResponseWriter, r *http.Request) (n *Nav, err error) {
	// user is teacher name: try getting classes
	classlist, err := v.store.ListClasses(user)
	if err != nil {
		fmt.Printf("Nav couldn't make classlist for %s: %v", user, err)
	}

	formatClassList(classlist)
	a, _ := v.ReadAlert(w, r)

	n = &Nav{
		User:      user,
		ClassList: classlist,
		Path:      path,
		Title:     title,
		Type:      stype,
		Alert:     a,
	}
	return n, nil
}
