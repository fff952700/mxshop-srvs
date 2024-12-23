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

### protoc 使用

protoc  
1、go install google.golang.org/protobuf/cmd/protoc-gen-go@latest  
2、go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  
3、生成protoc:  
&nbsp;&nbsp;protoc --go_out=. --go-grpc_out=. user.proto

### 前端代码
>https://github.com/OctopusLian/mxshop.git