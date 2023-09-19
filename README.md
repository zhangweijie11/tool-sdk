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
