package config

import "fmt"

//ModelList represents model
type ModelList struct {
	Models []*Model
}

//Init initialises model list
func (l *ModelList) Init() {
	if len(l.Models) == 0 {
		return
	}
	for i := range l.Models {
		l.Models[i].Init()
	}
}

//Validate validates model list
func (l *ModelList) Validate() error {
	if len(l.Models) == 0 {
		r