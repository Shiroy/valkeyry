package values

import "container/list"

//go:generate stringer -type=ValueKind

type ValueKind int

const (
	ValueKindString ValueKind = iota
	ValueKindInteger
	ValueKindList
)

type Value interface {
	Kind() ValueKind
}

type ValueList struct {
	data *list.List
}

func (v ValueList) Kind() ValueKind {
	return ValueKindList
}

func NewValueList(data *list.List) *ValueList {
	return &ValueList{
		data: data,
	}
}

type ValueString struct {
	Data string
}

func (v ValueString) Kind() ValueKind {
	return ValueKindString
}

func NewValueString(value string) *ValueString {
	return &ValueString{
		Data: value,
	}
}
