package interpreter

import "github.com/cstuartroe/teko/src/checker"

type Set struct {
	elements []*TekoObject
}

func (s Set) getUnderlyingType() checker.TekoType {
	return nil // TODO
}

func (s Set) getFieldValue(name string) *TekoObject {
	switch name {
	default:
		panic("Attributes haven't been implemented for sets yet")
	}
}
