# go-micro-tpl

go 语言版的微服务模板项目。

## 初始化

### Windows

项目中使用了包 `github.com/confluentinc/confluent-kafka-go` , 需要依赖 GCC 进行编译，Windows 下需要安装 mingw 。可以手动下载安装包安装，下载地址为 https://sourceforge.net/projects/mingw-w64/files/latest/download ，安装完后将安装目录下 bin 文件夹添加到环境变量 PATH 中；或者通过 choco 进行安装，命令为 `choco install mingw -y`。

由于项目使用了 Makefile，所以 Windows 下需要使用 gnu make。 可以使用命令 `choco install make -y` 进行安装；也可以直接下载安装 http://gnuwin32.sourceforge.net/downlinks/make.php ，安装完后将安装目录下 bin 文件夹添加到环境变量 PATH 中。

### Linux

如果系统中之前没有安装过 GCC 等编译工具化，需要手动安装一下，以 Ubuntu 为例，安装命令如下：

`apt-get install build-essential`

## 配置



### 环境变量 

`HTTP_PORT` http 监听端口号，不设置，则不开启 HTTP 服务  
`GRPC_PORT` grpc 监听端口，不设置，则不开启 Grpc 服务  
`KAFKA_STARTUP_DISABLED` 是否禁用 kafka，默认启用 kafka ，传 `true` 则禁用 kafka 监听  

### 配置文件
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

#### grpc 模式

首先给出一下 protobuf 定义：

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

请求体 `MicroRequest` 分为三部分： `头信息` `用户信息` 和 `请求正文`。其中 `头信息` 中的 `msg_type` 被服务端用来做路径映射使用，如果不存在当前 `msg_type` 对应的路由信息，则代表当前请求不被支持。虽然整个数据结构使用的 protobuf 来定义，但是这里依然把 `请求正文` 字段设计为 JSON 字符串格式 ，而不是纯 protobuf 对象，这样新增请求消息类型时，客户端和服务端不用修改 protobuf 文件定义，显得更加灵活。

响应体 `MicroResponse` 分为三部分：`头信息` `请求结果` `响应正文`。其中 `头信息` 中除了 `msg_type` 字段之外，其他字段和其匹配的请求体中的 `头信息` 中的值保持一致。

#### kafka 模式

kafka 模式下消费者得到的数据包中的 value 字段是一个 JSON 字符串，其数据结构和 grpc 的请求消息体格式类似：

```javascript
{
    user: Object,
    header: Object,
    req_data: Object
}
```

其中 header 属性和 user 属性的类型 和 protobuf 中定义的一致，由于 kafka 模式下消息正文本身就是一个 JSON，所以 req_data 的类型不是字符串而是 JSON Object。

kafka 模式没响应正文数据包，因为这个模式的设计初衷是用来承载异步消息处理，适合请求者发送完请求后，不关心处理结果的场景。如果当前请求处理出错，可以使用报警信息来提示相关人员。

#### http 模式

http 模式使用 POST JSON 的方式提交请求数据，JSON 的格式和 kafka 模式是相同的，同时为了方便后端服务进行路由，每个请求都需要携带请求头 `message-msg-type`，其值为 JSON 数据中的 `header.msg_type` 字段值。 


### 路由

#### 服务器端

服务器端路由位于文件 `src/route/route.go` 中，其格式定义为以 `header.msg_type` 为 key，以 `reqcontext.ServiceConfig` 为 value。下面是一个示例：

```go
var RouteMap = map[string]reqcontext.ServiceConfig{
	"demo_get": {Service: service.DemoGet},
}
```

上述代码指示，对于请求类型为 `demo_get` 的请求使用函数 `service.DemoGet` 进行处理。

### 通信协议
虽然目前仅支持 grpc kafka http 三种通信协议，但是如果想增加新的通信协议的话，需要根据使用的驱动情况做分别处理。
#### 显式返回请求数据
第三方协议的驱动需要显式拿到处理后的结果，在驱动内部做发送，则需要做如下处理：

首先在 src/filter 文件夹下，写一个当前新增协议的解析处理文件，假设命名为 `xx_filter.go`，在其中暴漏一个如下函数：
```go
// XXRequest XXResponse 分别为当前协议特定的请求和响应数据结构
func DoXXFilter(in *XXRequest) (*XXResponse, error) {
    //在里面实现从 XXRequest 到 reqcontext.RequestPacket 的转化，以方便用来调用路由函数，
    //调用完路由函数后，将得到 reqcontext.MicroResponse 结构再转化为 XXResponse 进行返回。
}
```
#### 调用驱动暴漏的发送函数

在 src/filter 中，新增一个文件，假设命名为 `xx_filter.go`，然后在其中实现如下函数：
```go
// XXRequest XXResponse 分别为当前协议特定的请求和响应数据结构
func DoXXFilter(in *XXRequest) {
    //在里面实现从 XXRequest 到 reqcontext.RequestPacket 的转化，以方便用来调用路由函数，
    //并根据 reqcontext.RequestPacket 对象初始化一个 reqcontextXX 对象，
    //将 reqcontext.RequestPacket 对象和  reqcontextXX 对象传入路由函数
}
```

上述伪代码中的 `reqcontextXX` 需要继承自 `src\reqcontext\request_context.go` 下的 `BaseContext` 结构体，且实现接口 `ContextInterface` 中的 `doError` 和 `doSuccess` 两个函数。

### 打点

#### 访问日志

为了让服务的数据更好的做追踪，代码中设计了访问日志写入逻辑，写入的字段定义参见代码 `src\reqcontext\access_log.go` 中的 `AccessLog` 结构。


## LICENSE

[MIT](LICENSE)

## TODO

- [ ] 配置项更改时动态感知到
- [ ] 错误码定义
- [ ] docker 镜像构建
- [ ] 数据库使用示例