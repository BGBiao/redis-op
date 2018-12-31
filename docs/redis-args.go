/**
 * @File Name: redis-args.go
 * @Author: xxbandy @http://xxbandy.github.io
 * @Email:
 * @Create Date: 2018-04-13 17:04:45
 * @Last Modified: 2018-12-31 16:12:02
 * @Description:
 */

package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	redisop "github.com/xxbandy/redis-op/pool"
)

func main() {
	pool := redisop.NewPool("localhost", "32772", "")
	defer pool.Close()
	conn := pool.Get()
	defer conn.Close()

	// sort uid desc by user_day_* get user_name_* get user_day_*
	// lpushO, setErr := redis.Int(conn.Do("lpush", redis.Args{}.Add("ansible-key").AddFlat(allip)...))
	// args := "desc by user_day_* get user_name_* get user_day_*"
	/*
	  args := map[string]string{
	      "by":    					"user_day_*",
	      "get":            "get user_name_*",
	  }
	*/
	//sortv,sortErr := redis.Strings(conn.Do("sort", redis.Args{}.Add("uid desc").AddFlat(args)...))
	sortv, sortErr := redis.Strings(conn.Do("sort", "uid", "desc", "by", "user_day_*", "get","#","get", "user_name_*", "get", "user_day_*"))
	if sortErr != nil {
		fmt.Printf("%s\n", sortErr)
	}
	fmt.Println(sortv)
}
