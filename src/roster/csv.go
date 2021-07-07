package roster

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path"
)

// rowToStu415 takes a slice of 11 from a csv row and retuns a student
func rowToStu401s(r []string) (s *Stu401) {

	return &Stu401{
		PermID:          r[0],
		Gender:          r[1],
		Grade:           r[2],
		BirthDate:       r[3],
		StudentName:     r[4],
		Track:           r[5],
		StudentGU:       r[6],
		PrimaryLanguage: r[7],
		HomeLanguage:    r[8],
		FirstName:       r[9],
		LastName:        r[10],
	}
}

// GetFullPath returns $HOME/.data/$fileName"
// File is shared among roster apps to avoid repeated Downloads
func GetFullPath(fileName string) string {
	var base string

	u, err := user.Current()
	if err == nil {
		base = u.HomeDir
	}

	base = path.Join(base, ".data", "Roster")

	if _, err := os.Stat(base); os.IsNotExist(err) {

		if err = os.MkdirAll(base, os.ModePerm); err != nil {
			log.Fatal(err)
		}

	}
	path := path.Join(base, fileName)

	return path

}

// ReadStu401s parses Stu401s from provided reader
func ReadStu401sFromCSV(csvReport io.Reader) (students []*Stu401, err error) {
	rows, err := csv.NewReader(csvReport).ReadAll()
	if err != nil {
		return nil, err
	}

	students = make([]*Stu401, len(rows))
	fmt.Printf("ReadStu415sFromCSV found %d rows", len(rows))

	for i, r := range rows {
		students[i] = rowToStu401s(r)
	}
	return students, nil
}
