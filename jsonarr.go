package zjson

import (
	"fmt"
	"strconv"
	"sync"
)

type JsonArray struct {
	data []any
	mu   sync.RWMutex // 添加互斥锁以支持并发安全
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
		return nil, fmt.Errorf("failed to parse value: %w", err)
	}

	var arrayVal = make([]any, 0)
	if err := jsonParser.JsonStringToAny(strB, &arrayVal); err != nil {
		return nil, fmt.Errorf("failed to parse JSON array: %w", err)
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
	ja.mu.Lock()
	defer ja.mu.Unlock()
	ja.data = append(ja.data, value)
}

func (ja *JsonArray) Get(index int) any {
	ja.mu.RLock()
	defer ja.mu.RUnlock()
	if index < 0 || index >= len(ja.data) {
		return nil
	}
	return ja.data[index]
}

func (ja *JsonArray) Remove(index int) {
	ja.mu.Lock()
	defer ja.mu.Unlock()
	if index < 0 || index >= len(ja.data) {
		return
	}
	ja.data = append(ja.data[:index], ja.data[index+1:]...)
}

func (ja *JsonArray) Length() int {
	ja.mu.RLock()
	defer ja.mu.RUnlock()
	return len(ja.data)
}

func (ja *JsonArray) ToJsonStr() string {
	ja.mu.RLock()
	defer ja.mu.RUnlock()
	jsonStr, err := jsonParser.AnyToJsonString(ja.data)
	if err != nil {
		return "[]"
	}
	return string(jsonStr)
}

func (ja *JsonArray) ToStruct(s any) error {
	ja.mu.RLock()
	defer ja.mu.RUnlock()
	if err := jsonParser.JsonStringToAny([]byte(ja.ToJsonStr()), s); err != nil {
		return fmt.Errorf("failed to convert to struct: %w", err)
	}
	return nil
}

func (ja *JsonArray) GetJsonObject(index int) (*JsonObject, error) {
	ja.mu.RLock()
	defer ja.mu.RUnlock()

	if index < 0 || index >= len(ja.data) {
		return nil, fmt.Errorf("index %d out of bounds for array of length %d", index, len(ja.data))
	}

	return ParseToJsonObject(ja.data[index])
}

func (ja *JsonArray) GetJsonObjectIgnoreError(index int) *JsonObject {
	jsonObject, _ := ja.GetJsonObject(index)
	return jsonObject
}

func (ja *JsonArray) GetJsonArray(index int) (*JsonArray, error) {
	ja.mu.RLock()
	defer ja.mu.RUnlock()

	if index < 0 || index >= len(ja.data) {
		return nil, fmt.Errorf("index %d out of bounds for array of length %d", index, len(ja.data))
	}

	return ParseToArray(ja.data[index])
}

func (ja *JsonArray) GetJsonArrayIgnoreError(index int) *JsonArray {
	jsonArray, _ := ja.GetJsonArray(index)
	return jsonArray
}

func (ja *JsonArray) GetInt(index int) (int, error) {
	ja.mu.RLock()
	defer ja.mu.RUnlock()

	if index < 0 || index >= len(ja.data) {
		return 0, fmt.Errorf("index %d out of bounds for array of length %d", index, len(ja.data))
	}

	val := ja.data[index]
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		if number, err := strconv.ParseInt(v, 10, 64); err == nil {
			return int(number), nil
		}
	}
	return 0, fmt.Errorf("value at index %d is not an integer", index)
}

func (ja *JsonArray) GetIntIgnoreError(index int) int {
	intVal, _ := ja.GetInt(index)
	return intVal
}

func (ja *JsonArray) GetFloat(index int) (float64, error) {
	ja.mu.RLock()
	defer ja.mu.RUnlock()

	if index < 0 || index >= len(ja.data) {
		return 0, fmt.Errorf("index %d out of bounds for array of length %d", index, len(ja.data))
	}

	val := ja.data[index]
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		if number, err := strconv.ParseFloat(v, 64); err == nil {
			return number, nil
		}
	}
	return 0, fmt.Errorf("value at index %d is not a float", index)
}

func (ja *JsonArray) GetFloatIgnoreError(index int) float64 {
	floatVal, _ := ja.GetFloat(index)
	return floatVal
}

func (ja *JsonArray) GetString(index int) (string, error) {
	ja.mu.RLock()
	defer ja.mu.RUnlock()

	if index < 0 || index >= len(ja.data) {
		return "", fmt.Errorf("index %d out of bounds for array of length %d", index, len(ja.data))
	}

	val := ja.data[index]
	if strVal, ok := val.(string); ok {
		return strVal, nil
	}
	return fmt.Sprint(val), nil
}

func (ja *JsonArray) GetStringIgnoreError(index int) string {
	strVal, _ := ja.GetString(index)
	return strVal
}
