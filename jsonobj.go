package zjson

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	errKeyNotExist = errors.New("key does not exist")
	errValueType   = errors.New("value type mismatch")
)

type JsonObject struct {
	data map[string]any
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
		return nil, err
	}
	var mapVal = make(map[string]any)
	if err := jsonParser.JsonStringToAny(strB, &mapVal); err != nil {
		return nil, err
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
	jo.data[key] = value
}

func (jo *JsonObject) Get(key string) any {
	if value, exists := jo.data[key]; exists {
		return value
	}
	return nil
}

func (jo *JsonObject) Remove(key string) {
	delete(jo.data, key)
}

func (jo *JsonObject) ContainsKey(key string) bool {
	_, exists := jo.data[key]
	return exists
}

func (jo *JsonObject) Length() int {
	return len(jo.data)
}

func (jo *JsonObject) ToJsonStr() string {
	jsonStr, err := jsonParser.AnyToJsonString(jo.data)
	if err != nil {
		return "{}"
	}
	return string(jsonStr)
}

func (jo *JsonObject) ToStruct(s any) error {
	if err := jsonParser.JsonStringToAny([]byte(jo.ToJsonStr()), s); err != nil {
		return err
	}
	return nil
}

func (jo *JsonObject) GetInt(key string) (int, error) {
	val, exist := jo.data[key]
	if !exist {
		return 0, errKeyNotExist
	}
	if intVal, ok := val.(int); ok {
		return intVal, nil
	}
	if int64Val, ok := val.(int64); ok {
		return int(int64Val), nil
	}
	if floatVal, ok := val.(float64); ok {
		return int(floatVal), nil
	}
	if intStr, ok := val.(string); ok {
		if number, err := strconv.ParseInt(intStr, 10, 64); err == nil {
			return int(number), nil
		}
	}
	return 0, errValueType
}

func (jo *JsonObject) GetFloat(key string) (float64, error) {
	val, exist := jo.data[key]
	if !exist {
		return 0, errKeyNotExist
	}
	if floatVal, ok := val.(float64); ok {
		return floatVal, nil
	}
	if floatStr, ok := val.(string); ok {
		if number, err := strconv.ParseFloat(floatStr, 64); err == nil {
			return number, nil
		}
	}
	return 0, errValueType
}

func (jo *JsonObject) GetString(key string) (string, error) {
	val, exist := jo.data[key]
	if !exist {
		return "", errKeyNotExist
	}
	if strVal, ok := val.(string); ok {
		return strVal, nil
	}
	return fmt.Sprint(val), nil
}

func (jo *JsonObject) GetBool(key string) (bool, error) {
	val, exist := jo.data[key]
	if !exist {
		return false, errKeyNotExist
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal, nil
	}
	if strVal, ok := val.(string); ok {
		if strVal == "true" || strVal == "True" {
			return true, nil
		} else if strVal == "false" || strVal == "False" {
			return false, nil
		}
	}
	return false, errValueType
}

func (jo *JsonObject) GetJsonObject(key string) (*JsonObject, error) {
	val, exist := jo.data[key]
	if !exist {
		return nil, errKeyNotExist
	}
	return ParseToJsonObject(val)
}

func (jo *JsonObject) GetJsonArray(key string) (*JsonArray, error) {
	val, exist := jo.data[key]
	if !exist {
		return nil, errKeyNotExist
	}
	return ParseToArray(val)
}
