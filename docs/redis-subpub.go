/**
 * @File Name: abc.go
 * @Author: xxbandy @http://xxbandy.github.io 
 * @Email:
 * @Create Date: 2018-04-09 23:04:16
 * @Last Modified: 2018-04-09 23:04:01
 * @Description:
 */
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/gomodule/redigo/redis"
)

//
func listenPubSubChannels(ctx context.Context, redisServerAddr string,
    onStart func() error,
    onMessage func(channel string, data []byte) error,
    channels ...string) error {

    //一分钟的心跳检测
    const healthCheckPeriod = time.Minute

    //构建一个redis链接
    c, err := redis.Dial("tcp", redisServerAddr,
        redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
        redis.DialWriteTimeout(10*time.Second))
    if err != nil {
        return err
    }
    defer c.Close()

    //构建一个pubsub链接
    psc := redis.PubSubConn{Conn: c}

    //订阅channels
    if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
        return err
    }
    //构造chan来检测通知状态
    done := make(chan error, 1)

    //启动一个goroutine来接受来自server端的通知
    go func() {
        for {
            //使用interface{}.(type) 来获取对应的类型，并借助switch和case进行interface{}的类型判断
            switch n := psc.Receive().(type) {
            case error:
                done <- n
                return
            case redis.Message:
                if err := onMessage(n.Channel, n.Data); err != nil {
                    done <- err
                    return
                }
            case redis.Subscription:
                switch n.Count {
                case len(channels):
                    // Notify application when all channels are subscribed.
                    if err := onStart(); err != nil {
                        done <- err
                        return
                    }
                case 0:
                    // Return from the goroutine when all channels are unsubscribed.
                    done <- nil
                    return
                }
            }
        }
    }()

    ticker := time.NewTicker(healthCheckPeriod)
    defer ticker.Stop()
loop:
    for err == nil {
        select {
        case <-ticker.C:
            // Send ping to test health of connection and server. If
            // corresponding pong is not received, then receive on the
            // connection will timeout and the receive goroutine will exit.
            if err = psc.Ping(""); err != nil {
                break loop
            }
        case <-ctx.Done():
            break loop
        case err := <-done:
            // Return error from the receive goroutine.
            return err
        }
    }

    // Signal the receiving goroutine to exit by unsubscribing from all channels.
    psc.Unsubscribe()

    // Wait for goroutine to complete.
    return <-done
}

func publish(redisServerAddr string) {
    c, err := redis.Dial("tcp",redisServerAddr)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer c.Close()

    c.Do("PUBLISH", "c1", "hello")
}

// This example shows how receive pubsub notifications with cancelation and
// health checks.
func main() {
    redisServerAddr := "localhost:32771"

    ctx, cancel := context.WithCancel(context.Background())

    //ctx和start callback很好的解决了丢失的消息的填充(使用goroutine来占住message类型的消息)
    listenErr := listenPubSubChannels(ctx,
        redisServerAddr,
        func() error {
            // The start callback is a good place to backfill missed
            // notifications. For the purpose of this example, a goroutine is
            // started to send notifications.
            go publish(redisServerAddr)
            return nil
        },
        func(channel string, message []byte) error {
            fmt.Printf("channel: %s, message: %s\n", channel, message)

            // For the purpose of this example, cancel the listener's context
            // after receiving last message sent by publish().
            if string(message) == "goodbye" {
                cancel()
            }
            return nil
        },
        "ansible-key")

    if listenErr != nil {
        fmt.Println(listenErr)
        return
    }

}
