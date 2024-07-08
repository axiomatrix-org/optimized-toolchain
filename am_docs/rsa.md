# RSA 非对称加密套件
用于进行RSA非对称加密操作。

## 使用方法
### 产生密钥对
密钥对由公钥和私钥组成，均为PEM格式。
```go
// 参数1: 密钥对长度
var privateKey, publicKey, err = am_rsa.GenRSAKeypair(2048)
```

### 公钥加密
```go
var originData = "123456"
// 参数1: 原始数据
// 参数2: 公钥PEM
am_rsa.RsaEncryptBase64(originData, publicKey)
```

### 私钥解密
```go
var encryptedData = "..." // encrypted data
// 参数1: 加密过后的数据
// 参数2: 私钥PEM
am_rsa.RsaDecryptBase64(encryptedData, privateKey)
```