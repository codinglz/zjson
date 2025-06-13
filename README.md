# zjson

一个简单、高效的 Go JSON 处理库。

## 主要特性

### 1. 并发安全性
- 为 `JsonObject` 和 `JsonArray` 添加了 `sync.RWMutex` 互斥锁
- 所有读写操作都使用适当的锁保护
- 写操作使用 `Lock()`，读操作使用 `RLock()`
- 添加了并发安全性测试用例

### 2. 错误处理改进
- 使用 `fmt.Errorf` 包装错误，提供更多上下文信息
- 统一了错误消息格式
- 添加了更详细的错误描述
- 改进了错误处理测试用例

### 3. 类型转换优化
- 使用 `switch` 语句替代多个 `if` 判断，提高代码可读性
- 增加了对 `int64` 类型的支持
- 改进了字符串到数字的转换逻辑
- 添加了更多类型转换测试用例

### 4. 测试覆盖
- 添加了边界条件测试
- 添加了并发安全性测试
- 添加了性能测试
- 添加了错误处理测试
- 添加了嵌套结构测试

## 使用示例

### 并发安全使用
```go
var wg sync.WaitGroup
obj := zjson.NewJsonObject()

wg.Add(2)
go func() {
    defer wg.Done()
    for i := 0; i < 1000; i++ {
        obj.Put("key1", i)
    }
}()

go func() {
    defer wg.Done()
    for i := 0; i < 1000; i++ {
        obj.Put("key2", i)
    }
}()

wg.Wait()
```

### 改进的错误处理
```go
obj := zjson.NewJsonObject()
obj.Put("age", "25")

age, err := obj.GetInt("age")
if err != nil {
    // 错误信息会包含更多上下文
    log.Printf("获取年龄失败: %v", err)
    return
}
```

### 类型转换
```go
obj := zjson.NewJsonObject()
obj.Put("number", "123.45")

// 支持多种数字类型
val, err := obj.GetInt("number")
if err == nil {
    fmt.Printf("整数: %d\n", val)
}

floatVal, err := obj.GetFloat("number")
if err == nil {
    fmt.Printf("浮点数: %.2f\n", floatVal)
}
```

## 安装

```bash
go get github.com/codinglz/zjson
```

## 版本历史

### v0.2.0
- 添加了并发安全性支持
- 改进了错误处理机制
- 优化了类型转换逻辑
- 添加了全面的测试用例
- 改进了代码质量

### v0.1.0
- 初始版本发布
- 基本的 JSON 对象和数组操作
- 类型转换支持
- 错误处理

## 贡献者
- [@codinglz](https://github.com/codinglz)

## 许可证
MIT License
