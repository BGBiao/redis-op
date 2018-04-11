## 基于redigo封装的redis连接  
最近在研究golang环境如何更好的使用redis，过程中根据以往的运维经验和个人开发爱好封装了[redigo](http://github.com/gomodule/redigo/redis)的连接池，以方便以后进行项目的快速开发和学习。

## 安装和使用
`注意：我个人是在golang1.8.3环境下封装的，其他版本不确定会不会有问题`

```
$ go get -v github.com/xxbandy/redis-op/pool
$ cat test-conn.go
package main
import (
    "fmt"
    "github.com/gomodule/redigo/redis"
    redisop "github.com/xxbandy/redis-op/pool"
)

func main() {
    //创建一个新连接池输入redis的ip,port,passwd
    pool := redisop.NewPool("localhost","6379","passwd")
    defer pool.Close()
    conn := pool.Get()
    defer conn.Close()

    //使用获取到redis连接进行操作
    _,setErr := redis.String(conn.Do("set","bgops","https://xxbandy.github.io"))
    redisop.ErrCheck("set error:",setErr)

    output,_ := redis.String(conn.Do("get","bgops"))
    fmt.Println("blog url:",output)
}
$ go run test-conn.go
blog url: https://xxbandy.github.io
```
