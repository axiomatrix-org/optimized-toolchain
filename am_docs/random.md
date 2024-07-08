# 随机生成字串
用于随机产生指定位数的纯数字字串或数字/大写字母混合字串。

## 使用方法
```go
// 产生6位随机纯数字
var digits = am_random.CreateRandomDigits(6)
// 产生6位随机大写字母和数字混合
var fusion = am_random.CreateRandomString(6)
```