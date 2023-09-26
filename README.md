# Tool-SDK

# 项目介绍

本项目是一个工具类的 SDK，给工具赋能，可以让工具更加关注本身能力的优秀实现，其他的功能操作由 SDK 提供！

# 功能介绍

- [x] Restful 接口风格、统一的参数校验功能
- [x] 提供基础任务数据模型、可快速调用数据库
- [x] MySQL、Redis、Elasticsearch 连接和初始化功能
- [x] 任务增删改查、暂停、重启
- [x] 多方式（API、MQ、gRPC）任务结果回调、进度回调
- [x] 全局任务并发量控制、单任务（生产者消费者模型）并发量控制
- [ ] Docker 容器化
- [ ] 自适应任务并发量调度



## 开始使用

```
go get -u gitlab.example.com/zhangweijie/tool-sdk
```

## 注意事项

### config文件

- server_name为工具名称，也是数据模型work中WokrType的值，入库时会直接提取配置中的server_name
- 需要加密的配置项字段名称需要包含`password`或者`apikey`,加密数据前缀必须为`ENC~`,满足以上标准会自动进行解密操作



### 任务参数

> 回调方式只有三种:`API、MQ、gRPC`

#### 1、API回调

- 当回调方式为 `API `时，回调地址为 `progressUrl+/progress`（进度回调），示例：http://10.100.40.35/progress，意味着接收方必须实现` /progress` 接口。

  ```
  示例数据：
  {
         "workUUID":   "91D9e3e3-f20f-e382-8528-6efF3dDdC68A",
         "serverName": component,
         "progress":   96,
  }
  ```

- 当回调方式为 `API `时，回调地址为 `callbackUrl+/callback/result`（结果回调），示例：http://10.100.40.35/callback/result，意味着接收方必须实现 `/callback/result` 接口。

  ```
  示例数据：
  {
         "workUUID":   "91D9e3e3-f20f-e382-8528-6efF3dDdC68A",
         "serverName": component,
         "result":   {"aaa":"bbb"},
  }
  ```

#### 2、MQ 回调

当回调方式为 `MQ`时，回调地址会按照`,`切分，一共需要三个数据，`addr,exchange,queue`，示例：`amqp://guest:guest@rabbitmq-server:5672/,component,component`

- 进度回调

  ```
  示例数据：
  {
         "workUUID":   "91D9e3e3-f20f-e382-8528-6efF3dDdC68A",
         "serverName": component,
         "progress":   96,
  }
  ```

- 结果回调

  ```
  示例数据：
  {
         "workUUID":   "91D9e3e3-f20f-e382-8528-6efF3dDdC68A",
         "serverName": component,
         "result":   {"aaa":"bbb"},
  }
  ```



#### 3、gRPC 回调

当回调方式为`gRPC`时，回调地址为`callbackUrl`，示例：http://10.100.40.35:5001

- 进度回调

  ```
  示例数据：
  {
         "workUUID":   "91D9e3e3-f20f-e382-8528-6efF3dDdC68A",
         "serverName": component,
         "progress":   96,
  }
  ```

- 结果回调

  ```
  示例数据：
  {
         "workUUID":   "91D9e3e3-f20f-e382-8528-6efF3dDdC68A",
         "serverName": component,
         "result":   {"aaa":"bbb"},
  }
  ```

  