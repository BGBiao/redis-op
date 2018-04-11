/**
 * @File Name: redis-hash.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-04 08:04:55
 * @Last Modified: 2018-04-09 15:04:48
 * @Description: 使用redis的连接池封装pub/sub的消息发布者
 */
package main
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
func errCheck(tp string,err error) {
    if err != nil {
        fmt.Printf("sorry,has some error for %s.\r\n",tp,err)
        os.Exit(-1)
    }
}

//构造一个连接池
//url为包装了redis的连接参数ip,port,passwd
func newPool(ip,port,passwd string) *redis.Pool {
    return &redis.Pool{
        MaxIdle:            5,    //定义redis连接池中最大的空闲链接为3
        MaxActive:          18,    //在给定时间已分配的最大连接数(限制并发数)
        IdleTimeout:        240 * time.Second,
        MaxConnLifetime:    300 * time.Second,
        Dial:               func() (redis.Conn,error) { return redisConn(ip,port,passwd) },
    }
}


func publish(rPool *redis.Pool,channel,message string) {
    //从连接池获取一个可用链接
    c := rPool.Get()
    defer c.Close()
    c.Do("publish",channel,message)

}

func main() {
    //使用newPool构建一个redis连接池
    pool := newPool("localhost","32771","123qweasd")
    defer pool.Close()
    publish(pool,"ansible-key","172.25.60.148")

    time.Sleep(1*time.Second)
}
