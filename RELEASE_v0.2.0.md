# zjson v0.2.0 发布说明

## 主要更新

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

### 5. 代码质量改进
- 统一了代码风格
- 改进了变量命名
- 添加了适当的空行提高可读性
- 优化了代码结构

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

## 更新日志
- 添加了 `json_advanced_test.go` 文件，包含全面的测试用例
- 修改了 `jsonobj.go` 和 `jsonarr.go` 的实现
- 添加了并发安全性支持
- 改进了错误处理机制
- 优化了类型转换逻辑

## 向后兼容性
此版本保持了与 v0.1.0 的向后兼容性，所有现有代码都可以继续使用。新增的功能和优化不会影响现有代码的行为。

## 贡献者
- [@codinglz](https://github.com/codinglz)

## 下载
- 源代码: [v0.2.0.zip](https://github.com/codinglz/zjson/archive/v0.2.0.zip)
- 源代码: [v0.2.0.tar.gz](https://github.com/codinglz/zjson/archive/v0.2.0.tar.gz) 