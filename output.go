package main

// Provides methods for handling output for matches
type OutputHandler interface {
	AddMatches(path string, matches []*Match)
	AddPathError(path string, e error)
}
