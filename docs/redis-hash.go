/**
 * @File Name: redis-hash.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-04 08:04:55
 * @Last Modified: 2018-04-08 23:04:42
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

//构造实际场景的hash结构体
var p1,p2 struct {
    Description  string `redis:"description"`
    Url          string `redis:"url"`
    Author       string `redis:"author"`
}



//主函数
func main() {
    c,cErr := redisConn("localhost","32771","123qweasd")
    errCheck("Conn",cErr)

    defer c.Close()
    p1.Description = "my blog"
    p1.Url = "http://xxbandy.github.io"
    p1.Author = "bgbiao"

    _,hmsetErr := c.Do("hmset",redis.Args{}.Add("hao123").AddFlat(&p1)...)
    errCheck("hmset",hmsetErr)

    m := map[string]string{
        "description":    "oschina",
        "url":            "http://my.oschina.net/myblog",
        "author":         "xxbandy",
    }

    _,hmset1Err := c.Do("hmset",redis.Args{}.Add("hao").AddFlat(m)...)
    errCheck("hmset1",hmset1Err)

    for _,key := range []string{"hao123","hao"} {
        v, err := redis.Values(c.Do("hgetall",key))
        errCheck("hmgetV",err)
        //等同于hgetall的输出类型，输出字符串为k/v类型
        //hashV,_ := redis.StringMap(c.Do("hgetall",key))
        //fmt.Println(hashV)
        //等同于hmget 的输出类型，输出字符串到一个字符串列表
        hashV2,_ := redis.Strings(c.Do("hmget",key,"description","url","author"))
        for _,hashv := range hashV2 {
                fmt.Println(hashv)
            }
        if err := redis.ScanStruct(v,&p2);err != nil {
            fmt.Println(err)
            return
        }
    fmt.Printf("%+v\n",p2)


    }

}

