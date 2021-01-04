package rfc3986

const (
	state_scheme = iota
)

type parserState struct {
	index int
	text  string
	state int
}
