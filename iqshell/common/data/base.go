package data

// Bool bool 引用类型
type Bool bool

func NewBool(b bool) *Bool {
	return (*Bool)(&b)
}

func (b *Bool) Value() bool {
	if b == nil {
		return false
	}
	return bool(*b)
}

func GetNotEmptyBoolIfExist(values ...*Bool) *Bool {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

// Int int 引用类型
type Int int

func (i *Int) Value() int {
	if i == nil {
		return 0
	}
	return int(*i)
}

func NewInt(i int) *Int {
	return (*Int)(&i)
}

func GetNotEmptyIntIfExist(values ...*Int) *Int {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

// Int int 引用类型
type Int64 int64

func (i *Int64) Value() int64 {
	if i == nil {
		return 0
	}
	return int64(*i)
}

func NewInt64(i int64) *Int64 {
	return (*Int64)(&i)
}

func GetNotEmptyInt64IfExist(values ...*Int64) *Int64 {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

// String string 引用类型
type String string

func NewString(s string) *String {
	return (*String)(&s)
}

func (s *String) Value() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

func GetNotEmptyStringIfExist(values ...*String) *String {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func Empty(value interface{}) bool {
	if v, ok := value.(*Bool); ok {
		return v == nil
	} else if v, ok := value.(*Int); ok {
		return v == nil
	} else if v, ok := value.(*Int64); ok {
		return v == nil
	} else if v, ok := value.(*String); ok {
		return v == nil || len(v.Value()) == 0
	} else {
		return value == nil
	}
}

func NotEmpty(value interface{}) bool {
	return !Empty(value)
}
