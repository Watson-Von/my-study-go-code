package main

import(
	_ "encoding/gob"
	"fmt"
	_ "io"
	_ "os"
	"sync"
	"time"
)

type Item struct{
	// 数据源
	Object interface{}
	// 有效时间
	Expiration int64
}

// 判断数据项是否已经过期
func (item Item) Expired() bool{

	if item.Expiration == 0 {
		return false
	}

	// 当前时间纳秒
	return time.Now().UnixNano() > item.Expiration
}

const (
	// 没有过期时间标志
	NoExpiration time.Duration = -1
	
	// 默认的过期时间
	DefaultExpiration time.Duration = 0

	Nanosecond  time.Duration = 1
    Microsecond          = 1000 * Nanosecond
    Millisecond          = 1000 * Microsecond
    Second               = 1000 * Millisecond
    Minute               = 60 * Second
    Hour                 = 60 * Minute
)

// 缓存结构体
type Cache struct {
    defaultExpiration   time.Duration
    itemMap             map[string]Item // 缓存数据项存储在 map 中
    mu                  sync.RWMutex    // 读写锁
    gcInterval          time.Duration   // 过期数据项清理周期
    stopGc              chan bool
}

// 设置缓存数据项，如果数据项存在则覆盖
func (c *Cache) Set(k string, v interface{}, d time.Duration) {
    var e int64
    if d == DefaultExpiration {
        d = c.defaultExpiration
    }
    if d > 0 {
        e = time.Now().Add(d).UnixNano()
    }
    c.mu.Lock()
    defer c.mu.Unlock()
    c.itemMap[k] = Item {
        Object:     v,
        Expiration: e,
    }
}

// 获取数据项，如果找到数据项，还需要判断数据项是否已经过期
func (c *Cache) get(k string) (interface{}, bool) {
    item, found := c.itemMap[k]
    if !found {
        return nil, false
    }
    if item.Expired() {
        return nil, false
    }
    return item.Object, true
}

// 获取数据项
func (c *Cache) Get(k string) (interface{}, bool) {
    c.mu.RLock()
    item, found := c.itemMap[k]
    if !found {
        c.mu.RUnlock()
        return nil, false
    }
    if item.Expired() {
        return nil, false
    }
    c.mu.RUnlock()
    return item.Object, true
}


// 删除缓存数据项
func (c *Cache) delete(k string) {
	// 内建函数
    delete(c.itemMap, k)
}

// 删除一个数据项
func (c *Cache) Delete(k string) {
    c.mu.Lock()
    c.delete(k)
    c.mu.Unlock()
}

// 删除过期数据项
func (c *Cache) DeleteExpired() {
    now := time.Now().UnixNano()
    c.mu.Lock()
    defer c.mu.Unlock()

    for k, v := range c.itemMap {
        if v.Expiration > 0 && now > v.Expiration {
			fmt.Println("k=", k, " 被删除了")
            c.delete(k)
        }
    }
}

// 过期缓存数据项清理
func (c *Cache) gcLoop() {

	// 定时器
    ticker := time.NewTicker(c.gcInterval)
    for {
        select {
		case <-ticker.C:
			 fmt.Println(time.Now(), " 定时器开始检查过期的key.....")
             c.DeleteExpired()
		case <-c.stopGc:
			 fmt.Println(time.Now(), " 定时器关闭......")
             ticker.Stop()
             return
        }
    }
}

// 停止过期缓存清理
func (c *Cache) StopGc() {
    c.stopGc <- true
}

// 创建一个缓存系统
func NewCache(defaultExpiration, gcInterval time.Duration) *Cache {
    c := &Cache{
        defaultExpiration: defaultExpiration,
        gcInterval:        gcInterval,
        itemMap:           map[string]Item{},
        stopGc:            make(chan bool),
    }
    // 开始启动过期清理 goroutine
    go c.gcLoop()
    return c
}

type Student struct{
	name string
}

func main() {

	cache := NewCache(1 * Second, 3 * Second)

	cache.Set("a",Student{name:"watson"}, 5 * Second)

	a,result := cache.Get("a")
	fmt.Println(a)
	fmt.Println(result)

	time.Sleep(4*time.Second)
	fmt.Println(a)
	fmt.Println(result)

	time.Sleep(10*time.Second)
	cache.StopGc()

	time.Sleep(50*time.Second)

}
