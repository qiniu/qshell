package data

// Bool bool 引用类型
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

// Int int 引用类型
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

// String string 引用类型
type String *string

func StringValue(s String) string {
	if s == nil {
		return ""
	}
	return *s
}

func NewString(s string) String {
	return &s
}
