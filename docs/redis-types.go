/**
 * @File Name: redis-args.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-13 17:04:45
 * @Last Modified: 2018-05-06 18:05:21
 * @Description:
 */

package main
import (
    "fmt"
    "github.com/gomodule/redigo/redis"
    redisop "github.com/xxbandy/redis-op/pool"
)

func main() {
    pool := redisop.NewPool("localhost","32771","")
    defer pool.Close()
    conn := pool.Get()
    defer conn.Close()

    //redis.Values return a []interface {}
    //redis.Values个人理解会用在复杂场景的hash类型?
    /*
    hmgetV,_ := redis.Values(conn.Do("hgetall","user"))
    //注意v是一个interface{}类型的
    for k,v := range hmgetV {
        //v is a interface{} == []uint8
        //convert v(a interface) to type. v.([]uint8)
        fmt.Println(k)
        for _,v1 := range v.([]uint8) {
            fmt.Println(string(v1))
            }
    }
    */


    //redis.Strings return a []string
		lrangeO,_ := redis.Strings(conn.Do("lrange",redis.Args{}.Add("ansible-key").AddFlat(0).AddFlat(-1)...))
    mgetString,_ := redis.Strings(conn.Do("mget","name","key","url"))
    hmget,_ := redis.Strings(conn.Do("hgetall","user"))
    fmt.Println(lrangeO,mgetString,hmget)

    //redis.String return a string. eg:ok,string value,nil...
    //set 一般返回OK字符串
    setStr,_ := redis.String(conn.Do("set","key","123"))
    fmt.Println(setStr)
    //get 一般会返回字符串，但如果类型为int也可以返回整型
    getStr,_ := redis.String(conn.Do("get","num"))
    getStrint,_ := redis.Int(conn.Do("get","num"))
    fmt.Println(getStr,getStrint)

    //redis.Int return a int
    //del 返回值为1表示删除成功，0表示不存在
    delID,_ := redis.Int(conn.Do("del","ansible-key"))
    fmt.Println(delID)

    //redis.Bool return a bool(redis内部其实会把bool值转化成0、1来去内部判断)
    keyExist,_ := redis.Bool(conn.Do("exists","url"))
    fmt.Println(keyExist)


}
