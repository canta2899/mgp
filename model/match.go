package model

type Match struct {
  LineNumber int 
  Content    string
}

func NewMatch(line int, content string) *Match {
  return &Match{
    LineNumber: line,
    Content:   content, 
  }
}

