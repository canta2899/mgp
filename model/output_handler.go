package model

type OutputHandler interface {
	AddMatch(path string, match *Match)
	AddMatches(path string, matches []*Match)
	AddPathError(path string, e error)
}
