package types

type Object = any
type ObjectMap = map[string]Object
type StringMap = map[string]string

func MergeObjectMap(this, other *ObjectMap) {
	for k, v := range *other {
		(*this)[k] = v
	}
}

func MergeStringMap(this, other *StringMap) {
	for k, v := range *other {
		(*this)[k] = v
	}
}
