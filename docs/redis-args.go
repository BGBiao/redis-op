/**
 * @File Name: redis-args.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-13 17:04:45
 * @Last Modified: 2018-04-13 17:04:46
 * @Description:
 */

package main
import (
    "fmt"
    "github.com/gomodule/redigo/redis"
    redisop "github.com/xxbandy/redis-op/pool"
)

func main() {
    allip := "172.25.60.148 172.25.60.149"
    pool := redisop.NewPool("localhost","32771","123qweasd")
    defer pool.Close()
    conn := pool.Get()
    defer conn.Close()

		lpushO,setErr := redis.Int(conn.Do("lpush",redis.Args{}.Add("ansible-key").AddFlat(allip)...))
		if setErr == nil {
				fmt.Printf("lpush %d ip:%s",lpushO,allip)
		}


}
