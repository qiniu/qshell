package data

type Bool *bool

func BoolV(b Bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func NewBool(b bool) Bool {
	return &b
}