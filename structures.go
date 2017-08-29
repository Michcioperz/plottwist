package main

import "fmt"

type Series struct {
	Title string
	Alt   string
	Slug  string
}

func (s *Series) NiceTitle() string {
	if s.Alt != "" {
		return fmt.Sprintf("%v (%v)", s.Title, s.Alt)
	}
	return s.Title
}

type Episode struct {
	Number uint
	Source string
}
