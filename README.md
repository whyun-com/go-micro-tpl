# go-micro-tpl

go 语言版的微服务模板项目。

## 初始化

### Windows

项目中使用了包 `github.com/confluentinc/confluent-kafka-go` , 需要依赖 GCC 进行编译，Windows 下需要安装 mingw 。可以手动下载安装包安装，下载地址为 https://sourceforge.net/projects/mingw-w64/files/latest/download ，安装完后将安装目录下 bin 文件夹添加到环境变量 PATH 中；或者通过 choco 进行安装，命令为 `choco install mingw -y`。

由于项目使用了 Makefile，所以 Windows 下需要使用 gnu make。 可以使用命令 `choco install make` 进行安装；也可以直接下载安装 http://gnuwin32.sourceforge.net/downlinks/make.php ，安装完后将安装目录下 bin 文件夹添加到环境变量 PATH 中。

### Linux

如果系统中之前没有安装过 GCC 等编译工具化，需要手动安装一下，以 Ubuntu 为例，安装命令如下：

`apt-get install build-essential`

## 配置



需要依赖于环境变量 

`HTTP_PORT` http 监听端口号，不设置，则不开启 HTTP 服务
`GRPC_PORT` grpc 监听端口，不设置，则不开启 Grpc 服务
`KAFKA_STARTUP_DISABLED` 是否禁用 kafka，默认启用 kafka ，传 `true` 则禁用 kafka 监听

同时应用需要读取 yaml 格式的配置文件，应用内部会依此从如下路径读取配置文件：环境变量 `CONFIG_PATH` 指定的路径，默认路径 `/etc/micro-config.yaml`，如果这两处都没有读取到，应用会自行退出。

当前支持的配置选项：

```yaml
kafka:
  log:
    brokers: 访问日志 kafka 集群地址，逗号分隔各个broker，例如 1.2.3.4:9092,2.3.4.5:9092
  communication:
    brokers: 通信用 kafka 集群地址，逗号分隔各个broker，例如 1.2.3.4:9092,2.3.4.5:9092
redis:
  hosts:
    - redis 集群地址，逗号分隔各个节点，例如 1.2.3.4:6379,2.3.4.5:6379
```

## 构建

```shell
make build
```
构建完成后，会在 bin/ 目录下生成可执行程序 micro 。
## 运行

运行 `bin/micro` 会启动应用程序，不过要提前设置好环境变量和配置文件，具体参见配置小节。

## 测试
```shell
make test
```

## 构建镜像

## 代码结构

### 请求和响应

micro-tpl 提供了 http grpc kafka 三种通信方式，开发者可以根据需要选择其中一种或者几种来使用。

以 grpc 举例来说明请求和响应的数据字段定义，首先给出一下 protobuf 定义：

```protobuf
syntax = "proto3";
package micro;

message User {
    string user_id = 1;
}
message MessageHeader {
    string msg_type = 1;//消息类型
    uint64 req_ms = 2;//请求触发的时间点
    string req_id = 3;//请求id
    uint32 req_pid = 4;//请求进程id
    string req_ip = 5;//请求者的ip
    map<string,string> extended_fields = 6;//扩展字段
}
message MicroRequest {
    MessageHeader header = 2;
    User user = 3;
    string req_data_json = 4;//请求正文，json 字符串格式
}
message Result {
    uint32 code = 1;
    string msg = 2;
}
message MicroResponse {
    MessageHeader header = 1;
    Result result = 2;//请求结果
    string res_data_json = 3;//响应正文，json 字符串格式
}

service Micro {
    // Sends a greeting
    rpc DoRequest (MicroRequest) returns (MicroResponse) {
    }
}
```



## TODO

- [ ] 配置项更改时动态感知到
- [ ] 错误码定义
- [ ] docker 镜像构建
- [ ] 数据库使用示例