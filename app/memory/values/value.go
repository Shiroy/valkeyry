package values

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
}

func (v ValueList) Kind() ValueKind {
	return ValueKindList
}

type ValueString struct {
	Data string
}

func (v ValueString) Kind() ValueKind {
	return ValueKindString
}

func NewValueString(value string) ValueString {
	return ValueString{
		Data: value,
	}
}
