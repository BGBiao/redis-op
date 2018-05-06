/**
 * @File Name: redis-args.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-13 17:04:45
 * @Last Modified: 2018-05-06 15:05:29
 * @Description:
 */

package main
import (
    "fmt"
    "github.com/gomodule/redigo/redis"
    redisop "github.com/xxbandy/redis-op/pool"
)

func main() {
    allip := "10.0.0.1"
    pool := redisop.NewPool("localhost","32771","")
    defer pool.Close()
    conn := pool.Get()
    defer conn.Close()

		lpushO,setErr := redis.Int(conn.Do("lpush",redis.Args{}.Add("ansible-key").AddFlat(allip)...))
		if setErr == nil {
				fmt.Printf("lpush %d ip:%s",lpushO,allip)
		}


}
