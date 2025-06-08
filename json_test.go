package zjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonObject_PutAndGet(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("name", "Alice")
	assert.Equal(t, "Alice", obj.Get("name"))
	assert.Nil(t, obj.Get("age"))
}

func TestJsonObject_Remove(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("name", "Alice")
	obj.Put("age", 30)
	obj.Remove("name")
	assert.Nil(t, obj.Get("name"))
	assert.NotNil(t, obj.Get("age"))
}

func TestJsonObject_ContainsKey(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("name", "Alice")
	assert.True(t, obj.ContainsKey("name"))
	assert.False(t, obj.ContainsKey("age"))
}

func TestJsonObject_Length(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("name", "Alice")
	obj.Put("age", 30)
	assert.Equal(t, 2, obj.Length())
}

func TestJsonObject_ToJsonStr(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("name", "Alice")
	jsonStr := obj.ToJsonStr()
	assert.JSONEq(t, `{"name":"Alice"}`, jsonStr)
}

func TestJsonObject_ToStruct(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}
	obj := NewJsonObject()
	obj.Put("name", "Alice")
	var user User
	err := obj.ToStruct(&user)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
}

func TestJsonObject_GetInt(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("age", 30)
	obj.Put("height", 175.5)
	obj.Put("weight_str", "65")

	val, err := obj.GetInt("age")
	assert.NoError(t, err)
	assert.Equal(t, 30, val)

	val, err = obj.GetInt("height")
	assert.NoError(t, err)
	assert.Equal(t, 175, val)

	val, err = obj.GetInt("weight_str")
	assert.NoError(t, err)
	assert.Equal(t, 65, val)

	_, err = obj.GetInt("name")
	assert.Error(t, err)
}

func TestJsonObject_GetFloat(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("price", 9.99)
	obj.Put("discount_str", "0.85")

	val, err := obj.GetFloat("price")
	assert.NoError(t, err)
	assert.Equal(t, 9.99, val)

	val, err = obj.GetFloat("discount_str")
	assert.NoError(t, err)
	assert.Equal(t, 0.85, val)

	_, err = obj.GetFloat("name")
	assert.Error(t, err)
}

func TestJsonObject_GetString(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("name", "Alice")
	obj.Put("age", 30)

	val, err := obj.GetString("name")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", val)

	val, err = obj.GetString("age")
	assert.NoError(t, err)
	assert.Equal(t, "30", val)

	_, err = obj.GetString("gender")
	assert.Error(t, err)
}

func TestJsonObject_GetBool(t *testing.T) {
	obj := NewJsonObject()
	obj.Put("is_student", true)
	obj.Put("is_teacher_str", "True")

	val, err := obj.GetBool("is_student")
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = obj.GetBool("is_teacher_str")
	assert.NoError(t, err)
	assert.True(t, val)

	_, err = obj.GetBool("is_admin")
	assert.Error(t, err)
}

func TestJsonObject_GetJsonObject(t *testing.T) {
	obj := NewJsonObject()
	nestedObj := NewJsonObject()
	nestedObj.Put("city", "Beijing")
	obj.Put("address", nestedObj)

	res, err := obj.GetJsonObject("address")
	assert.NoError(t, err)
	assert.Equal(t, "Beijing", res.Get("city"))
}

func TestJsonObject_GetJsonArray(t *testing.T) {
	obj := NewJsonObject()
	arr := NewJsonArray()
	arr.Add("Apple")
	obj.Put("fruits", arr)

	res, err := obj.GetJsonArray("fruits")
	assert.NoError(t, err)
	assert.Equal(t, "Apple", res.Get(0))
}

func TestJsonArray_AddGetRemove(t *testing.T) {
	arr := NewJsonArray()
	arr.Add("Apple")
	arr.Add("Banana")

	assert.Equal(t, "Apple", arr.Get(0))
	assert.Equal(t, "Banana", arr.Get(1))

	arr.Remove(0)
	assert.Equal(t, "Banana", arr.Get(0))
}

func TestJsonArray_Length(t *testing.T) {
	arr := NewJsonArray()
	arr.Add("Apple")
	arr.Add("Banana")
	assert.Equal(t, 2, arr.Length())
}

func TestJsonArray_ToJsonStr(t *testing.T) {
	arr := NewJsonArray()
	arr.Add("Apple")
	jsonStr := arr.ToJsonStr()
	assert.JSONEq(t, `["Apple"]`, jsonStr)
}

func TestJsonArray_ToStruct(t *testing.T) {
	arr := NewJsonArray()
	arr.Add(map[string]any{"name": "Alice"})
	var users []struct {
		Name string `json:"name"`
	}
	err := arr.ToStruct(&users)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", users[0].Name)
}

func TestJsonArray_GetJsonObject(t *testing.T) {
	arr := NewJsonArray()
	nestedObj := NewJsonObject()
	nestedObj.Put("city", "Beijing")
	arr.Add(nestedObj)

	res, err := arr.GetJsonObject(0)
	assert.NoError(t, err)
	assert.Equal(t, "Beijing", res.Get("city"))

	arrStr := `{"city":"Shanghai"}`
	arr.Add(arrStr)
	res2, err := arr.GetJsonObject(1)
	assert.NoError(t, err)
	assert.Equal(t, "Shanghai", res2.Get("city"))
}

func TestJsonArray_GetJsonArray(t *testing.T) {
	arr := NewJsonArray()
	nestedArr := NewJsonArray()
	nestedArr.Add("Apple")
	arr.Add(nestedArr)

	res, err := arr.GetJsonArray(0)
	assert.NoError(t, err)
	assert.Equal(t, "Apple", res.Get(0))

	arrStr := `["Banana"]`
	arr.Add(arrStr)
	res2, err := arr.GetJsonArray(1)
	assert.NoError(t, err)
	assert.Equal(t, "Banana", res2.Get(0))
}

func TestJsonArray_GetInt(t *testing.T) {
	arr := NewJsonArray()
	arr.Add(30)
	arr.Add(175.5)
	arr.Add("65")

	val, err := arr.GetInt(0)
	assert.NoError(t, err)
	assert.Equal(t, 30, val)

	val, err = arr.GetInt(1)
	assert.NoError(t, err)
	assert.Equal(t, 175, val)

	val, err = arr.GetInt(2)
	assert.NoError(t, err)
	assert.Equal(t, 65, val)

	_, err = arr.GetInt(3)
	assert.Error(t, err)
}
