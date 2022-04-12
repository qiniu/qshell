package provider

type Item interface {
	Equal(item Item) bool
}
