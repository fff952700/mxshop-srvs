### mxshop-srvs

1、密码加密规则:采用slat+循环次数+key长度+hash512算法
> https://github.com/anaskhan96/go-password-encoder  
> go get github.com/anaskhan96/go-password-encoder

2、将srv服务注册到consul中,ip地址不能写localhost或127.0.0.1  
3、添加配置中间nacos
4、添加goods srv sql  
5、删除外键影响效率和重写category_brand实现  
6、goods_category_brand添加唯一索引category在前
7、去掉手动初始化改为自动初始化  
8、添加商品需要传入id并且没有唯一键约束不太合理  
9、清洗数据生成新的sql  
10、库存服务高并发时使用锁。需要连接gmp(goroutines,machine,processor三者的对应关系)  
  1) go的sync.RWMutex(并发读) 适用于高并发读，低频写。写操作会阻塞所有的读写操作，导致性能下降(适用于单机)  
  2) go的sync.Mutex 它保证同一时刻只有一个 goroutine 能够访问受保护的资源,适用于写操作(适用于单机)  
  3) 乐观锁 在go中通过版本号或者时间戳等控制。适用于低频读写(适用于分布式系统)
  4) 使用mysql的for update 在where带有索引的情况下变为行锁（无索引会升级为表锁）。适用高并发下的读写  
  5) 使用redis分布式锁redsync: 使用时需要考虑以下几个问题
     1) redsync严重依赖节点时间钟。不同步可能导致锁提前释放
     2) 当集群模式下数据未同步至 node / 2 + 1 并且总耗时小于锁的有效率则可能锁失败
     3) master故障可能导致数据丢失
     4) 单实例情况无高可用并且锁也可能丢失
  6) 需要综合考虑cap。选择redsync和mysql version来选择锁
  >https://github.com/go-redsync/redsync  

### protoc 使用

protoc  
1、go install google.golang.org/protobuf/cmd/protoc-gen-go@latest  
2、go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  
3、生成protoc:  
&nbsp;&nbsp;protoc --go_out=. --go-grpc_out=. user.proto

### 前端代码
>https://github.com/OctopusLian/mxshop.git