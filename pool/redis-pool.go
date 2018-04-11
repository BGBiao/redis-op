package pool
import (
    "fmt"
    "time"
    "os"

    "github.com/gomodule/redigo/redis"
)

//构造一个链接函数，如果没有密码，passwd为空字符串
func redisConn(ip,port,passwd string) (redis.Conn, error) {
    c,err := redis.Dial("tcp",
        ip+":"+port,
        redis.DialConnectTimeout(5*time.Second),
        redis.DialReadTimeout(1*time.Second),
        redis.DialWriteTimeout(1*time.Second),
        redis.DialPassword(passwd),
        redis.DialKeepAlive(1*time.Second),
        )
    return c,err
}

//构造一个错误检查函数
func ErrCheck(tp string,err error) {
    if err != nil {
        fmt.Printf("sorry,has some error for %s.\r\n",tp,err)
        os.Exit(-1)
    }
}


//构造一个连接池
//url为包装了redis的连接参数ip,port,passwd
func NewPool(ip,port,passwd string) *redis.Pool {
    return &redis.Pool{
        MaxIdle:            5,    //定义redis连接池中最大的空闲链接为3
        MaxActive:          18,    //在给定时间已分配的最大连接数(限制并发数)
        IdleTimeout:        240 * time.Second,
        MaxConnLifetime:    300 * time.Second,
        Dial:               func() (redis.Conn,error) { return redisConn(ip,port,passwd) },
    }
}


/*
func main() {
    //使用newPool构建一个redis连接池
    pool := NewPool("r2m-proxy.jdfin.local","6379","pe-test")
    defer pool.Close()
    c := pool.Get()
    defer c.Close()
    for {
        llen,_ := redis.Int64(c.Do("llen","ansible-key"))
        if int(llen) == 0 {
            time.Sleep(1*time.Second)
            continue
        }
        switch {
        case int(llen) > 0:
            r,_ := redis.String(c.Do("rpop","ansible-key"))
            fmt.Println("ansible-key: ",r)
        default:
            continue
        }
    }
}
*/
