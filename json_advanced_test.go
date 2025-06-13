package zjson

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试边界条件
func TestJsonObject_EdgeCases(t *testing.T) {
	obj := NewJsonObject()

	// 测试空键
	obj.Put("", "empty key")
	assert.Equal(t, "empty key", obj.Get(""))

	// 测试特殊字符键
	specialKey := "!@#$%^&*()_+"
	obj.Put(specialKey, "special")
	assert.Equal(t, "special", obj.Get(specialKey))

	// 测试nil值
	obj.Put("nil_value", nil)
	assert.Nil(t, obj.Get("nil_value"))

	// 测试大数值
	obj.Put("big_int", 9223372036854775807)
	val, err := obj.GetInt("big_int")
	assert.NoError(t, err)
	assert.Equal(t, 9223372036854775807, val)
}

// 测试类型转换的边界情况
func TestJsonObject_TypeConversionEdgeCases(t *testing.T) {
	obj := NewJsonObject()

	// 测试各种数字格式
	obj.Put("int_str", "123")
	obj.Put("float_str", "123.45")
	obj.Put("invalid_num", "not a number")

	val, err := obj.GetInt("int_str")
	assert.NoError(t, err)
	assert.Equal(t, 123, val)

	val, err = obj.GetInt("invalid_num")
	assert.Error(t, err)

	floatVal, err := obj.GetFloat("float_str")
	assert.NoError(t, err)
	assert.Equal(t, 123.45, floatVal)

	// 测试布尔值转换
	obj.Put("true_str", "true")
	obj.Put("false_str", "false")
	obj.Put("invalid_bool", "not a bool")

	boolVal, err := obj.GetBool("true_str")
	assert.NoError(t, err)
	assert.True(t, boolVal)

	boolVal, err = obj.GetBool("invalid_bool")
	assert.Error(t, err)
}

// 测试嵌套结构
func TestJsonObject_NestedStructures(t *testing.T) {
	obj := NewJsonObject()

	// 创建嵌套对象
	nestedObj := NewJsonObject()
	nestedObj.Put("name", "nested")
	obj.Put("nested", nestedObj)

	// 创建嵌套数组
	nestedArr := NewJsonArray()
	nestedArr.Add("item1")
	nestedArr.Add("item2")
	obj.Put("array", nestedArr)

	// 测试多层嵌套
	deepNested := NewJsonObject()
	deepNested.Put("level", 3)
	nestedObj.Put("deep", deepNested)

	res, err := obj.GetJsonObject("nested")
	assert.NoError(t, err)
	deep, err := res.GetJsonObject("deep")
	assert.NoError(t, err)
	level, err := deep.GetInt("level")
	assert.NoError(t, err)
	assert.Equal(t, 3, level)
}

// 测试并发安全性
func TestJsonObject_Concurrency(t *testing.T) {
	obj := NewJsonObject()
	var wg sync.WaitGroup
	iterations := 1000

	// 并发写入
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			obj.Put("key1", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			obj.Put("key2", i)
		}
	}()

	wg.Wait()

	// 验证结果
	assert.NotNil(t, obj.Get("key1"))
	assert.NotNil(t, obj.Get("key2"))
}

// 测试性能
func BenchmarkJsonObject_Operations(b *testing.B) {
	obj := NewJsonObject()

	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obj.Put("key", i)
		}
	})

	b.Run("Get", func(b *testing.B) {
		obj.Put("key", "value")
		for i := 0; i < b.N; i++ {
			obj.Get("key")
		}
	})

	b.Run("ToJsonStr", func(b *testing.B) {
		obj.Put("key1", "value1")
		obj.Put("key2", 123)
		for i := 0; i < b.N; i++ {
			obj.ToJsonStr()
		}
	})
}

// 测试JSON解析错误处理
func TestJsonObject_ParseErrors(t *testing.T) {
	// 测试无效的JSON字符串
	invalidJSON := `{"key": "value", "unclosed: true}`
	obj, err := ParseToJsonObject(invalidJSON)
	assert.Error(t, err)
	assert.Nil(t, obj)

	// 测试类型不匹配
	obj = NewJsonObject()
	obj.Put("key", make(chan int)) // 不支持的类型
	jsonStr := obj.ToJsonStr()
	assert.Equal(t, "{}", jsonStr)
}

// 测试数组边界条件
func TestJsonArray_EdgeCases(t *testing.T) {
	arr := NewJsonArray()

	// 测试空数组
	assert.Equal(t, 0, arr.Length())
	assert.Nil(t, arr.Get(0))

	// 测试负数索引
	assert.Nil(t, arr.Get(-1))

	// 测试越界索引
	arr.Add("item")
	assert.Nil(t, arr.Get(1))

	// 测试移除不存在的索引
	arr.Remove(-1)
	arr.Remove(1)
	assert.Equal(t, 1, arr.Length())
}

// 测试数组类型转换
func TestJsonArray_TypeConversion(t *testing.T) {
	arr := NewJsonArray()

	// 添加各种类型的值
	arr.Add(123)
	arr.Add("456")
	arr.Add(789.0)
	arr.Add(true)
	arr.Add(nil)

	// 测试类型转换
	val, err := arr.GetInt(0)
	assert.NoError(t, err)
	assert.Equal(t, 123, val)

	val, err = arr.GetInt(1)
	assert.NoError(t, err)
	assert.Equal(t, 456, val)

	floatVal, err := arr.GetFloat(2)
	assert.NoError(t, err)
	assert.Equal(t, 789.0, floatVal)

	strVal, err := arr.GetString(3)
	assert.NoError(t, err)
	assert.Equal(t, "true", strVal)

	// 测试nil值
	assert.Nil(t, arr.Get(4))
}

// 测试数组并发安全性
func TestJsonArray_Concurrency(t *testing.T) {
	arr := NewJsonArray()
	var wg sync.WaitGroup
	iterations := 1000

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			arr.Add(i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if arr.Length() > 0 {
				arr.Remove(0)
			}
		}
	}()

	wg.Wait()
}

// 测试数组性能
func BenchmarkJsonArray_Operations(b *testing.B) {
	arr := NewJsonArray()

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arr.Add(i)
		}
	})

	b.Run("Get", func(b *testing.B) {
		arr.Add("value")
		for i := 0; i < b.N; i++ {
			arr.Get(0)
		}
	})

	b.Run("ToJsonStr", func(b *testing.B) {
		arr.Add("item1")
		arr.Add(123)
		for i := 0; i < b.N; i++ {
			arr.ToJsonStr()
		}
	})
}
