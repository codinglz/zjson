package zjson

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

var (
	errKeyNotExist = errors.New("key does not exist")
	errValueType   = errors.New("value type mismatch")
)

type JsonObject struct {
	data map[string]any
	mu   sync.RWMutex // 添加互斥锁以支持并发安全
}

func ParseToJsonObject(v any) (*JsonObject, error) {
	var strB []byte
	if strP, ok := getPointVal[JsonObject](v); ok {
		return strP, nil
	} else if strP, ok := getPointVal[[]byte](v); ok {
		strB = *strP
	} else if strP, ok := getPointVal[string](v); ok {
		strB = []byte(*strP)
	} else if strByte, err := jsonParser.AnyToJsonString(v); err == nil {
		return ParseToJsonObject(string(strByte))
	} else {
		return nil, fmt.Errorf("failed to parse value: %w", err)
	}

	var mapVal = make(map[string]any)
	if err := jsonParser.JsonStringToAny(strB, &mapVal); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &JsonObject{
		data: mapVal,
	}, nil
}

func NewJsonObject() *JsonObject {
	return &JsonObject{
		data: make(map[string]any),
	}
}

func (jo *JsonObject) Put(key string, value any) {
	jo.mu.Lock()
	defer jo.mu.Unlock()
	jo.data[key] = value
}

func (jo *JsonObject) Get(key string) any {
	jo.mu.RLock()
	defer jo.mu.RUnlock()
	if value, exists := jo.data[key]; exists {
		return value
	}
	return nil
}

func (jo *JsonObject) Remove(key string) {
	jo.mu.Lock()
	defer jo.mu.Unlock()
	delete(jo.data, key)
}

func (jo *JsonObject) ContainsKey(key string) bool {
	jo.mu.RLock()
	defer jo.mu.RUnlock()
	_, exists := jo.data[key]
	return exists
}

func (jo *JsonObject) Length() int {
	jo.mu.RLock()
	defer jo.mu.RUnlock()
	return len(jo.data)
}

func (jo *JsonObject) ToJsonStr() string {
	jo.mu.RLock()
	defer jo.mu.RUnlock()
	jsonStr, err := jsonParser.AnyToJsonString(jo.data)
	if err != nil {
		return "{}"
	}
	return string(jsonStr)
}

func (jo *JsonObject) ToStruct(s any) error {
	jo.mu.RLock()
	defer jo.mu.RUnlock()
	if err := jsonParser.JsonStringToAny([]byte(jo.ToJsonStr()), s); err != nil {
		return fmt.Errorf("failed to convert to struct: %w", err)
	}
	return nil
}

func (jo *JsonObject) GetInt(key string) (int, error) {
	jo.mu.RLock()
	defer jo.mu.RUnlock()

	val, exist := jo.data[key]
	if !exist {
		return 0, fmt.Errorf("%w: key '%s'", errKeyNotExist, key)
	}

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
	return 0, fmt.Errorf("%w: key '%s' is not an integer", errValueType, key)
}

func (jo *JsonObject) GetIntIgnoreError(key string) int {
	val, _ := jo.GetInt(key)
	return val
}

func (jo *JsonObject) GetFloat(key string) (float64, error) {
	jo.mu.RLock()
	defer jo.mu.RUnlock()

	val, exist := jo.data[key]
	if !exist {
		return 0, fmt.Errorf("%w: key '%s'", errKeyNotExist, key)
	}

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
	return 0, fmt.Errorf("%w: key '%s' is not a float", errValueType, key)
}

func (jo *JsonObject) GetFloatIgnoreError(key string) float64 {
	val, _ := jo.GetFloat(key)
	return val
}

func (jo *JsonObject) GetString(key string) (string, error) {
	jo.mu.RLock()
	defer jo.mu.RUnlock()

	val, exist := jo.data[key]
	if !exist {
		return "", fmt.Errorf("%w: key '%s'", errKeyNotExist, key)
	}

	if strVal, ok := val.(string); ok {
		return strVal, nil
	}
	return fmt.Sprint(val), nil
}

func (jo *JsonObject) GetStringIgnoreError(key string) string {
	val, _ := jo.GetString(key)
	return val
}

func (jo *JsonObject) GetBool(key string) (bool, error) {
	jo.mu.RLock()
	defer jo.mu.RUnlock()

	val, exist := jo.data[key]
	if !exist {
		return false, fmt.Errorf("%w: key '%s'", errKeyNotExist, key)
	}

	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		switch v {
		case "true", "True", "TRUE":
			return true, nil
		case "false", "False", "FALSE":
			return false, nil
		}
	}
	return false, fmt.Errorf("%w: key '%s' is not a boolean", errValueType, key)
}

func (jo *JsonObject) GetBoolIgnoreError(key string) bool {
	val, _ := jo.GetBool(key)
	return val
}

func (jo *JsonObject) GetJsonObject(key string) (*JsonObject, error) {
	jo.mu.RLock()
	defer jo.mu.RUnlock()

	val, exist := jo.data[key]
	if !exist {
		return nil, fmt.Errorf("%w: key '%s'", errKeyNotExist, key)
	}

	return ParseToJsonObject(val)
}

func (jo *JsonObject) GetJsonObjectIgnoreError(key string) *JsonObject {
	val, _ := jo.GetJsonObject(key)
	return val
}

func (jo *JsonObject) GetJsonArray(key string) (*JsonArray, error) {
	jo.mu.RLock()
	defer jo.mu.RUnlock()

	val, exist := jo.data[key]
	if !exist {
		return nil, fmt.Errorf("%w: key '%s'", errKeyNotExist, key)
	}

	return ParseToArray(val)
}

func (jo *JsonObject) GetJsonArrayIgnoreError(key string) *JsonArray {
	val, _ := jo.GetJsonArray(key)
	return val
}
