package domain

type Grade string

const (
	Grade5C  Grade = "5C"
	Grade5Cp Grade = "5C+"
	Grade6A  Grade = "6A"
	Grade6Ap Grade = "6A+"
	Grade6B  Grade = "6B"
	Grade6Bp Grade = "6B+"
	Grade6C  Grade = "6C"
	Grade6Cp Grade = "6C+"
	Grade7A  Grade = "7A"
	Grade7Ap Grade = "7A+"
	Grade7B  Grade = "7B"
)

var gradeRatings = map[Grade]int{
	Grade5C:  1000,
	Grade5Cp: 1100,
	Grade6A:  1200,
	Grade6Ap: 1300,
	Grade6B:  1400,
	Grade6Bp: 1500,
	Grade6C:  1600,
	Grade6Cp: 1700,
	Grade7A:  1800,
	Grade7Ap: 1900,
	Grade7B:  2000,
}

func (g Grade) Rating() (int, error) {
	r, ok := gradeRatings[g]
	if !ok {
		return 0, ErrUnknownGrade
	}
	return r, nil
}

func (g Grade) Valid() bool {
	_, ok := gradeRatings[g]
	return ok
}

func AllGrades() []Grade {
	return []Grade{
		Grade5C, Grade5Cp,
		Grade6A, Grade6Ap,
		Grade6B, Grade6Bp,
		Grade6C, Grade6Cp,
		Grade7A, Grade7Ap,
		Grade7B,
	}
}
