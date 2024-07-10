package am_jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axiomatrix-org/optimized-toolchain/am_redis"
	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 工具返回代码
const (
	OKCode                  = 2000
	InternalServerErrorCode = 5000
	ErrCodeMalformed        = 4001 // token格式错误
	ErrCodeExpired          = 4002 // token过期
	ErrCodeNotValidYet      = 4003 // token尚未生效
	ErrCodeSignatureInvalid = 4004 // token签名无效
	ErrCodeInvalidRole      = 4005 // token权限不符
	ErrCodeNoAuthHeader     = 4006 // 没有对应的请求头
	ErrCodeMalformedReq     = 4007 // 请求头格式不正确
)

// 工具错误类型
var (
	MalformedError         = errors.New("malformed jwt")
	ExpiredAndDiedError    = errors.New("expired jwt")
	ExpiredButCanSaveError = errors.New("expired jwt, but can save it")
	NotValidError          = errors.New("invalid jwt")
	SignatureError         = errors.New("expired jwt")
	InvalidRoleError       = errors.New("invalid role")
	UnknownError           = errors.New("unknown jwt error")
)

// 不同身份权限，数字越大，权限越大。高权限用户可以执行低权限用户的所有操作
const (
	ROOTROLE  = 4 // root权限，最大的权限，整个系统的超级管理权限
	ADMINROLE = 3 // admin权限，只能由root授予，拥有管理用户的权限，但不具备修改root的权限
	USERROLE  = 2 // user权限，普通用户权限
	TEMPROLE  = 1 // temp权限，临时权限，该权限仅用于注册和修改密码的临时使用
)

// 身份字段转换权限等级数字
func ClaimToRole(claim string) int {
	switch claim {
	case "root":
		return ROOTROLE
	case "admin":
		return ADMINROLE
	case "user":
		return USERROLE
	case "temp":
		return TEMPROLE
	default:
		return 0
	}
}

// token claims
type TokenClaims struct {
	Email              string // 用户的email
	Role               string // 用户的身份登记
	Exp                int    // token过期时间，以秒计数
	Issuer             string // 签发人
	SECRET             string // token secret
	jwt.StandardClaims        // standard claims，无需用户设定
}

// token config
type TokenConfig struct {
	Issuer string `json:"issuer"`
	Secret string `json:"secret"`
}

/*
* PRIVATE
* 查验redis储存情况
* 参数：
* 1. token string：要查验的token
 */

func checkRedis(token string) (*TokenClaims, error) {
	// 从redis中查验
	claimsJson, err := am_redis.GetValue(token)
	claims := &TokenClaims{}
	if err != nil {
		if errors.Is(err, am_redis.RedisGetNilError) { // 如果不存在于库中，则证明已经提前过期
			return nil, ExpiredAndDiedError
		} else {
			return nil, err
		}
	}
	err = json.Unmarshal([]byte(claimsJson), claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

/*
* 生成token字串
* 参数：
* 1. claims *TokenClaims token参数
 */
func GenToken(claims *TokenClaims) (string, error) {
	tokenClaims := claims

	// 配置standard claims
	tokenClaims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(claims.Exp)).Unix(), // 过期时间
		IssuedAt:  time.Now().Unix(),                                              // 签名日期
		Issuer:    claims.Issuer,                                                  // 签发人
	}

	// 解析私鑰
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(claims.SECRET))
	if err != nil {
		return "", err
	}

	// 生成token字串
	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	token, err := tokenGenerator.SignedString(privKey) // 生成token
	if err != nil {
		return "", err
	}

	// 储存redis
	value, err := json.Marshal(claims)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	err = am_redis.SetValue(token, string(value), claims.Exp) // 键为token字串，值为原始过期时间
	if err != nil {
		return "", err
	}

	return token, nil
}

/*
* 解析token
* 参数：
* 1. token string：要解析的token字串
* 2. roleRequired int：验证通过需要的权限等级
* 3. secret string：解码密钥
 */
func ParseToken(token string, roleRequired int, secret string) (*TokenClaims, error) {
	// 解析公鑰
	pem, err := jwt.ParseRSAPublicKeyFromPEM([]byte(secret))
	if err != nil {
		return nil, err
	}

	// 解析token
	result, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return pem, nil
	})

	// 解析出现问题
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 { // token格式不正确
				return nil, MalformedError
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 { // token过期
				claims, err := checkRedis(token) // 查验redis
				if err != nil {
					if errors.Is(err, ExpiredAndDiedError) { // 如果redis也过期了
						return nil, ExpiredAndDiedError // 无可救药
					} else {
						return nil, err
					}
				}
				return claims, ExpiredButCanSaveError // 否则，还能救一救（重新生成token）
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 { // token尚未生效
				return nil, NotValidError
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 { // token签名错误
				return nil, SignatureError
			} else { // 其他错误
				return nil, err
			}
		}
	}

	// 解析未出现问题，即此token形式合法
	if claims, ok := result.Claims.(*TokenClaims); ok && result.Valid {
		// 从redis中查验
		_, err := checkRedis(token)
		if err != nil {
			return nil, err
		}

		// 刷新token在redis中的过期时间
		err = am_redis.SetValue(token, strconv.Itoa(claims.Exp), claims.Exp)
		if err != nil {
			return nil, err
		}

		// 验证权限合法性
		if ClaimToRole(claims.Role) >= roleRequired { // 如果提供的token权限验证大于所需权限，初步判断通过
			if roleRequired == 1 && ClaimToRole(claims.Role) == 2 { // user权限无权新增用户
				return nil, InvalidRoleError
			} // user不允許操作temp

			if ClaimToRole(claims.Role) == 1 { // temp权限仅用于注册和重设密码临时使用，一经使用立即灭活
				_, err := Kickoff(token)
				if err != nil {
					return nil, err
				}
			}
			return claims, nil // 返回claims
		} else {
			return nil, InvalidRoleError
		}
	}
	return nil, UnknownError
}

/*
* 灭活token
* 参数：
* 1. token string：需要灭活的token字串
 */
func Kickoff(token string) (bool, error) {
	err := am_redis.DelValue(token)
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
* JWT中间件
 */
func JWTAuthMiddleware(role string, secret string) func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization") // 获取请求头中的Authorization字段
		// 如果没有Authorization字段，直接拦截掉，并返回4006代码
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": ErrCodeNoAuthHeader,
				"msg":  "No Authorization header",
			})
			c.Abort()
			return
		}

		// 请求头中的token必须是Bearer-Token的格式，否则拦截，返回4007代码
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": ErrCodeMalformedReq,
				"msg":  "Wrong Authorization header",
			})
			c.Abort()
			return
		}

		// 开始认证
		claims, err := ParseToken(parts[1], ClaimToRole(role), secret)
		if err != nil {
			if errors.Is(err, ExpiredButCanSaveError) { // 如果还可以救一救
				tokenClaims := TokenClaims{
					Email:  claims.Email,
					Role:   claims.Role,
					Exp:    claims.Exp,
					Issuer: claims.Issuer,
					SECRET: claims.SECRET,
				}

				fmt.Println(tokenClaims) // test print

				token, err := GenToken(&tokenClaims) // 重新生成token
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": InternalServerErrorCode,
						"msg":  err.Error(),
					})
					c.Abort()
					return
				}

				c.Writer.Header().Set("Token", token) // 将token放入请求头后放行
				c.Set("email", claims.Email)
				c.Set("role", claims.Role)
				c.Next()
				return
			} else if errors.Is(err, ExpiredAndDiedError) { // 救不了了
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": ErrCodeExpired,
					"msg":  err.Error(),
				})
				c.Abort()
				return
			} else if errors.Is(err, MalformedError) { // 格式不正确
				c.JSON(http.StatusBadRequest, gin.H{
					"code": ErrCodeMalformed,
					"msg":  err.Error(),
				})
				c.Abort()
				return
			} else if errors.Is(err, InvalidRoleError) { // 权限不符
				c.JSON(http.StatusForbidden, gin.H{
					"code": ErrCodeInvalidRole,
					"msg":  err.Error(),
				})
				c.Abort()
				return
			} else if errors.Is(err, NotValidError) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": ErrCodeNotValidYet,
					"msg":  err.Error(),
				})
				c.Abort()
				return
			} else if errors.Is(err, SignatureError) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": ErrCodeSignatureInvalid,
					"msg":  err.Error(),
				})
				c.Abort()
				return
			} else if errors.Is(err, UnknownError) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": InternalServerErrorCode,
					"msg":  err.Error(),
				})
				c.Abort()
				return
			}
		}
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}
