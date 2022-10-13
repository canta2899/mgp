package model

type OutputHandler interface {
	AddMatches(path string, matches []*Match)
	AddPathError(path string, e error)
}
