package data

type String *string

func StringValue(s String) string {
	if s == nil {
		return  ""
	}
	return *s
}

func NewString(s string) String {
	return &s
}
