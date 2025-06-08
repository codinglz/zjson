package zjson

import "encoding/json"

var (
	jsonParser JsonParser = &DefaultParser{}
)

func SetParser(parser JsonParser) {
	jsonParser = parser
}

type JsonParser interface {
	AnyToJsonString(v any) ([]byte, error)
	JsonStringToAny(jsonStr []byte, v any) error
}

type DefaultParser struct {
}

func (p *DefaultParser) AnyToJsonString(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (p *DefaultParser) JsonStringToAny(jsonStr []byte, v any) error {
	return json.Unmarshal(jsonStr, v)
}

func getPointVal[T any](val any) (*T, bool) {
	if t, ok := val.(T); ok {
		return &t, true
	}
	if t, ok := val.(*T); ok {
		return t, true
	}
	return nil, false
}

