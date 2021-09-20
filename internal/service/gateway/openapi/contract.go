package openapi

type Publisher interface {
	Pub(value ...string) error
}
