package model

type Match struct {
	LineNumber int
	Content    string
}

// Creates a new match given the line number and the
// respective string content
func NewMatch(line int, content string) *Match {
	return &Match{
		LineNumber: line,
		Content:    content,
	}
}
