package alert

func CannotEmpty(words string, suggest string) string {
	desc := words
	if len(words) > 0 {
		desc += " can't empty"
	}
	return Description(desc, suggest)
}

func Description(desc string, suggest string) string {
	ret := desc
	if len(suggest) > 0 {
		ret += ", you can do like this:\n" + suggest
	}
	return ret
}
