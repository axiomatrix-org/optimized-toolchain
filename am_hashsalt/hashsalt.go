package am_hashsalt

import "golang.org/x/crypto/bcrypt"

/*
* 哈希加盐加密
* 参数：
* originData string：明文
 */
func HashData(originData string) (string, error) {
	// 哈希加盐加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(originData), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

/*
* 比对密文和明文
* 参数：
* 1. hashedData string：哈希加盐加密过后的数据
* 2. data string：要比对的明文数据
 */
func CompareData(hashedData, data string) bool {
	// 比对加密后数据和所需数据
	err := bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(data))
	return err == nil
}
