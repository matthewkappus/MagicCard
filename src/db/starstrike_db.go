package db

import "github.com/matthewkappus/MagicCard/src/comment"

// starstrike table
const (
	createStarStrike                 = `CREATE TABLE IF NOT EXISTS starstrike (id INTEGER PRIMARY KEY, perm_id TEXT, teacher TEXT, comment TEXT, title TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true, FOREIGN KEY(perm_id) REFERENCES stu415(perm_id))`
	insertStarStrike                 = `INSERT INTO starstrike(perm_id, teacher, comment, title, cat, isActive) VALUES(?,?,?,?,?,?);`
	selectStarStrikeByPermID         = `SELECT * FROM starstrike WHERE perm_id = ?`
	selectNewestStarStrikesByTeacher = `SELECT * FROM starstrike WHERE teacher=? LIMIT ?`

	// select count(cat) from starstrike where cat=1 and perm_id="980016917"
	selectStarStrikesByCatCount = `SELECT COUNT(cat) FROM starstrike WHERE cat = ? AND perm_id = ?`
)

func (s *Store) GetStarStrikesByPerm(id string) ([]*comment.StarStrike, error) {
	rows, err := s.db.Query(selectStarStrikeByPermID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := make([]*comment.StarStrike, 0)
	for rows.Next() {
		sb := new(comment.StarStrike)
		if err = rows.Scan(&sb.ID, &sb.PermID, &sb.Teacher, &sb.Comment, &sb.Title, &sb.Created, &sb.Cat, &sb.IsActive); err != nil {
			continue
		}
		ss = append(ss, sb)
	}

	return ss, nil

}

// type StarStrike struct {
// 	ID int `json:"id,omitempty"`
// 	PermID string `json:"perm_id,omitempty"`
// 	// staff(name)
// 	Teacher string `json:"teacher,omitempty"`
// 	Comment string `json:"comment,omitempty"`
// 	// Title is a catagory of the comment
// 	Title string `json:"title,omitempty"`

// 	Created time.Time `json:"created,omitempty"`
// 	// 0 star 1 minor 2 strik 3 major
// 	Cat      Category
// 	IsActive bool `json:"is_active,omitempty"`
// }
