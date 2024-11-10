package _type

// StringObject stores a string object
type StringObject struct {
	*BaseObject
	Value []byte
}

// GetType returns redis object type
func (o *StringObject) GetType() string {
	return StringType
}
