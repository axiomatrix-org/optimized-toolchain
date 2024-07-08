package am_ratelimit

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

// 储存用户信息，包括最后一次访问时间和访问计数
type RequestInfo struct {
	LastAccessTime time.Time
	RequestCount   int
}

// 互斥锁，用于同步对共享资源（如requestInfo映射）的访问，确保在多线程环境下的数据一致性
var mutex = &sync.Mutex{}

// 速率限制配置项结构体
type RateLimitConfig struct {
	maxRequests int                     // 最大访问次数
	timeWindow  time.Duration           // 窗口时间
	requestInfo map[string]*RequestInfo // 用户信息
}

/*
* 初始化速率限制配置
* 参数：
* 1. maxRequests int：窗口时间内最大访问数量限制
* 2. timeWindow int：限制的窗口时间，单位为秒
 */
func NewRateLimitConfig(maxRequests int, timeWindow int) *RateLimitConfig {
	return &RateLimitConfig{
		maxRequests: maxRequests,
		timeWindow:  time.Duration(timeWindow) * time.Second,
		requestInfo: make(map[string]*RequestInfo),
	}
}

/*
* 速率限制中间件
 */
func (r1 *RateLimitConfig) RateLimitMiddleware(c *gin.Context) {
	ip := c.ClientIP()   // 获取用户IP
	mutex.Lock()         // 在访问共享资源前加锁，确保只有同一线程才能访问
	defer mutex.Unlock() // 在数据返回前解锁

	info, exists := r1.requestInfo[ip] // 检查requestInfo映射中是否已经存在该IP信息

	// 如果该IP的请求信息不存在，则初始化一个新的RequestInfo实例，并将其添加到requestInfo映射中
	if !exists {
		r1.requestInfo[ip] = &RequestInfo{LastAccessTime: time.Now(), RequestCount: 1}
		return
	}

	// 如果最后一次请求的时间已经超出了时间窗口，重置请求计数并更新最后一次访问时间
	if time.Since(info.LastAccessTime) > r1.timeWindow {
		info.RequestCount = 1
		info.LastAccessTime = time.Now()
		return
	}

	// 增加请求数并检查速率限制
	info.RequestCount++
	if info.RequestCount > r1.maxRequests {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		c.Abort()
		return
	}

	// 更新最后一次访问时间
	info.LastAccessTime = time.Now()
}
