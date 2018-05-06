## redisgo使用说明
该doc文档对于[redigo](https://github.com/gomodule/redigo)的第三方包的相关功能进行学习和示例介绍，希望能够方便大家使用golang对redis业务场景的快速开发。

对于redis的使用无非就是以下三个步骤：
- 1.和redis服务端建立TCP链接
- 2.对于redis服务端进行读写操作
- 3.关闭TCP连接并对于后续的业务逻辑进行处理

在这里步骤1和3对于普通的工程师都比较熟悉了，需要注意的是步骤2里面的一些操作，在刚开始使用redis作为业务场景的一个模块时发现忽视了几个问题而导致程序并未达到最终理想的状态，这里需要对redis数据结构以及实际的场景有比较深的了解，这里大概总结下，避免以后重复学习。

需要注意的点:
- 1.对于不同的redis数据结构的不同操作，返回的数据结构也不相同，这个时候要选择合适的数据类型,redigo支持Bool, Int, Bytes, String, Strings and Values几种类型
- 2.一定要使用合适的redis命令来操作目标key，这样可以很大程度上提高性能并且不需要复杂的逻辑(当然你使用goroutine和异步锁也可以，不过并不建议这样做)
- 3.对于多线程或者多客户端针对某一个key的操作建议一定加锁，不然可能造成数据丢失或数据混乱的情况，这样会需要更多精力去清晰数据

对于redis的使用可参考如下文章:
[golang操作redis指南](https://www.jianshu.com/p/89ca34b84101)    
[redis开发指南](https://www.jianshu.com/p/69fc7f73eef1)    

`注意: 可以关注下每种redis命令对应输出的数据结构，这个关系到操作结果使用哪种数据类型进行接收返回结果，比如Bool,Int,Bytes,String,Strings,Values等`

## redisgo中的基本数据类型

`redis.Bool:` 存储bool型结果，比如`EXISTS`命令，虽然在在redis-command中看到是0和1,但其实在内部是转换成true和false来判断key是否存在的     
`redis.Int:` 存储int型结果，比如获取int型的value`get num`;删除操作`del key`,1表示删除成功，0表示key不存在;`append`、`incr`、`incrby`、`ttl`、`lpush\rpush` 、`hincrby`、`sadd` 等写入操作    
`redis.String:` 存储strings结构类型的数据，一般为`set/get`等单strings值的相关的操作，也用在返回结果为`nil`的操作中

`redis.Strings:` 存储返回为[]string类型的数据结构，比如`lrange key start end`或者`mget key1 key2`以及value为string类型的hash值之类的，结果一般为list类型
`redis.Values:` 存储返回为[]interface {}

**注意:**    
虽然在命令行中用户可以`set num 123` 不指定key的类型，但是在redis内部操作时会自动将123转换成`strings`类型进行存储，而用户在操作时可以根据类型来进行操作，也就是同样是`string`类型，我们是可以使用`redis.Int`和`redis.String`来操作`strings`类型的key。所以这个时候获取值的时候需要使用合适的类型，比如`redis.Int`或者`redis.Int64`以及`redis.String`

建议详细阅读[redis命令帮助](http://www.redis.cn/commands.html)，并对照上述几种基本类型进行数据操作

## redisgo的基本操作类型和使用案例
