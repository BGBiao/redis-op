/**
 * @File Name: redis-hash.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-04 08:04:55
 * @Last Modified: 2018-04-09 11:04:20
 * @Description:
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
func main() {
    c,connErr := redisConn("localhost","32771","123qweasd")
    errCheck("connErr",connErr)

    getV,_ := redis.String(c.Do("get","name"))
    getV2,_ := redis.Int(c.Do("get","id"))
    fmt.Println(getV,getV2)

    //mset mget
    _,setErr := c.Do("mset","name","biaoge","url","http://xxbandy.github.io")
    errCheck("setErr",setErr)
    if r,mgetErr := redis.Strings(c.Do("mget","name","url")); mgetErr == nil {
        for _,v := range r {
            fmt.Println("mget ",v)
            }
        
    }
}
