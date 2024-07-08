# Hash Salt 加密件
用于对密码等敏感信息进行加密入库。

## 使用方法
```go
// 加密信息
var originData = "123456" // 明文信息
var encryptedMessage = am_hashsalt.HashData(originData) // 加密操作，获得加密后的字串

// 匹配信息
var result = am_hashsalt.CompareData(encryptedData, originData) // 比对两者是否一致，返回bool值，true为一致，false为不一致
```