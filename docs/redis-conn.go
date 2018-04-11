/**
 * @File Name: redis-conn.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-04 07:04:57
 * @Last Modified: 2018-04-08 16:04:06
 * @Description:
 */
package main
import (
    "fmt"
    "os"
    "time"
    redis "github.com/gomodule/redigo/redis"
)
//32771
func main() {
    //使用redis封装的Dial进行tcp连接
    c,err := redis.Dial("tcp","localhost:32771",
        redis.DialKeepAlive(1*time.Second),
        redis.DialPassword("123qweasd"),
        redis.DialConnectTimeout(5*time.Second),
        redis.DialReadTimeout(1*time.Second),
        redis.DialWriteTimeout(1*time.Second),
        )
    errCheck(err)

    defer c.Close()

/*
    //对本次连接进行set操作
    _,setErr := c.Do("set","url","xxbandy.github.io")
    errCheck(setErr)
*/
    //使用redis的string类型获取set的k/v信息
    r,getErr := redis.String(c.Do("get","url"))
    errCheck(getErr)
    fmt.Println(r)
    
}

func errCheck(err error) {
    if err != nil {
        fmt.Println("sorry,has some error:",err)
        os.Exit(-1)
        }
    }
