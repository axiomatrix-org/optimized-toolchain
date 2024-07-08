package am_rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// 工具错误类型
var (
	KeyPairAlreadyExistError = errors.New("key pair already exist")
)

// KeyPair备份值
var keyPair map[string]string

/*
* 产生Keypair
* 参数：
* bits int：长度，一般为1024或2048
 */
func GenRSAKeypair(bits int) (privateKey, publicKey string, err error) {
	// 防止重复生成
	if keyPair != nil {
		return "", "", KeyPairAlreadyExistError
	}
	// 生成keypair
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 将私钥处理为PEM格式
	pkcs1PrivateKey := x509.MarshalPKCS1PrivateKey(key)
	p := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs1PrivateKey,
	}
	privKeyPEM := pem.EncodeToMemory(p)

	// 将公钥处理成PEM格式
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", err
	}
	p = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	pubKeyPEM := pem.EncodeToMemory(p)

	// 生成最后的PEM字串
	privateKey = string(privKeyPEM)
	publicKey = string(pubKeyPEM)

	// 备份keypair
	keyPair = map[string]string{
		"private": privateKey,
		"public":  publicKey,
	}
	return
}

/*
* RSA公钥加密
* 参数：
* 1. originData string：要加密的明文
* 2. publicKey string：公钥PEM
 */
func RsaEncryptBase64(originalData, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey)) // 解密公钥
	// 加载公钥
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 使用公钥加密数据
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(originalData))
	if err != nil {
		return "", err
	}

	// 产生base64字串
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

/*
* RSA私钥解密
* 参数：
* 1. encryptedData string：加密过后的密文
* 2. privateKey string：私钥PEM
 */
func RsaDecryptBase64(encryptedData, privateKey string) (string, error) {
	// 解密私钥
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(privateKey))

	// 加载私钥
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 使用私钥解密
	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, encryptedDecodeBytes)
	if err != nil {
		return "", err
	}

	// 返回原数据
	return string(originalData), nil
}
