package zjson

import (
	"fmt"
	"strconv"
)

type JsonArray struct {
	data []any
}

func ParseToArray(v any) (*JsonArray, error) {
	var strB []byte
	if strP, ok := getPointVal[JsonArray](v); ok {
		return strP, nil
	} else if strP, ok := getPointVal[[]byte](v); ok {
		strB = *strP
	} else if strP, ok := getPointVal[string](v); ok {
		strB = []byte(*strP)
	} else if strByte, err := jsonParser.AnyToJsonString(v); err == nil {
		return ParseToArray(string(strByte))
	} else {
		return nil, err
	}
	var arrayVal = make([]any, 0)
	if err := jsonParser.JsonStringToAny(strB, &arrayVal); err != nil {
		return nil, err
	}
	return &JsonArray{
		data: arrayVal,
	}, nil
}

func NewJsonArray() *JsonArray {
	return &JsonArray{
		data: make([]any, 0),
	}
}

func (ja *JsonArray) Add(value any) {
	ja.data = append(ja.data, value)
}

func (ja *JsonArray) Get(index int) any {
	if index < 0 || index >= len(ja.data) {
		return nil
	}
	return ja.data[index]
}

func (ja *JsonArray) Remove(index int) {
	if index < 0 || index >= len(ja.data) {
		return
	}
	ja.data = append(ja.data[:index], ja.data[index+1:]...)
}

func (ja *JsonArray) Length() int {
	return len(ja.data)
}

func (ja *JsonArray) ToJsonStr() string {
	jsonStr, err := jsonParser.AnyToJsonString(ja.data)
	if err != nil {
		return "[]"
	}
	return string(jsonStr)
}

func (ja *JsonArray) ToStruct(s any) error {
	if err := jsonParser.JsonStringToAny([]byte(ja.ToJsonStr()), s); err != nil {
		return err
	}
	return nil
}

func (ja *JsonArray) GetJsonObject(index int) (*JsonObject, error) {
	if index < 0 || index >= len(ja.data) {
		return nil, fmt.Errorf("index %d out of bounds for array of length %d", index, len(ja.data))
	}
	return ParseToJsonObject(ja.data[index])
}

func (ja *JsonArray) GetJsonArray(index int) (*JsonArray, error) {
	if index < 0 || index >= len(ja.data) {
		return nil, fmt.Errorf("index %d out of bounds for array of length %d", index, len(ja.data))
	}
	return ParseToArray(ja.data[index])
}

func (ja *JsonArray) GetInt(index int) (int, error) {
	if index < 0 || index >= len(ja.data) {
		return 0, fmt.Errorf("index %d out of bounds for array of length %d", index, len(ja.data))
	}
	val := ja.data[index]
	if val, ok := val.(float64); ok {
		return int(val), nil
	}
	if intVal, ok := val.(int); ok {
		return intVal, nil
	}
	if intStr, ok := val.(string); ok {
		if number, err := strconv.ParseInt(intStr, 10, 64); err == nil {
			return int(number), nil
		}
	}
	return 0, fmt.Errorf("index%d is not an int", index)
}
