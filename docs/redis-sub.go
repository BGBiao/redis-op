/**
 * @File Name: redis-hash.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-04 08:04:55
 * @Last Modified: 2018-05-06 18:05:15
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
 //       MaxIdle:            5,    //定义redis连接池中最大的空闲链接为3
 //       MaxActive:          18,    //在给定时间已分配的最大连接数(限制并发数)
 //       IdleTimeout:        240 * time.Second,
 //       MaxConnLifetime:    300 * time.Second,
        Dial:               func() (redis.Conn,error) { return redisConn(ip,port,passwd) },
    }
}


func publish(rPool *redis.Pool,channel,message string) {
    //从连接池获取一个可用链接
    c := rPool.Get()
    defer c.Close()
    c.Do("publish",channel,message)

}

func listenandSubscribe(rPool *redis.Pool,channel string)  {
    c := rPool.Get()
    defer c.Close()
    //使用redis连接初始化一个pubSub链接
    psc := redis.PubSubConn{Conn: c}

    //使用Subscribe订阅一个channel
    //if subErr := psc.Subscribe(redis.Args{}.AddFlat(channel)...); subErr != nil {
    //由于该subscribe命令是阻塞式的，因此需要使用context包进行上下文的锁定
    //ctx, cancel := context.WithCancel(context.Background())

    if subErr := psc.Subscribe(channel); subErr != nil {
        errCheck("subErr",subErr)
        os.Exit(-1)
    }
    //done := make(chan string,1)
    //使用goroutine进行阻塞subscribe
    go func() {
        for {
            ////psc.Receive() 返回的消息会为Subscription, Message, Pong or error
            /*
            func (c PubSubConn) Receive() interface{} {
                  return c.receiveInternal(c.Conn.Receive())
            }
            func (c PubSubConn) receiveInternal(replyArg interface{}, errArg error) interface{} {	
            */
            //n = [message pmessage | "subscribe", "psubscribe", "unsubscribe", "punsubscribe"| "pong"]
            //psc.Receive().(type) 判断psc.Receive()的类型
            switch rRecType := psc.Receive().(type) {
            case redis.Message:
                fmt.Println("messages for channel:%s, data:%s",rRecType.Channel,string(rRecType.Data))
                continue
            case redis.Subscription:
                fmt.Printf("subscribe type:%s\r\n",rRecType.Kind)
                continue
            case redis.Pong:
                fmt.Println("message for Pong:%s",string(rRecType.Data))
                continue
            default:
            //    done <- "error"
                continue

            }
        }
    }()
    time.Sleep(100*time.Minute)
}



func main() {
    //使用newPool构建一个redis连接池
    pool := newPool("localhost","32771","123qweasd")
    defer pool.Close()
    //publish(pool,"ansible-key","10.0.0.1")
    listenandSubscribe(pool,"ansible-key")

}
