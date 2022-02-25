package data

type Int *int

func IntValue(i Int) int {
	if i == nil {
		return 0
	}
	return *i
}

func NewInt(i int) Int {
	return &i
}
