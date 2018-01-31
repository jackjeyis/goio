# goio
## 高性能、轻量级 io 框架
### echo server sample

```go
package main

import (
  "github.com/jackjeyis/goio/application"
  "github.com/jackjeyis/goio/msg"
)

func main() {
  app := &application.GenericApplication{}
  app.SetOnStart(func() error {
      app.RegisterService("ping_handler",func(msg msg.Message) {
          msg.Channel().EncodeMessage(msg)
      })
      return nil
  }).Run()
}
```
控制账户鉴权和请求 PING 开放平台后端获取 token
各 APP 端登录鉴权后，通过中间层获取 PING 开放平台后端通信 token，由中间层拼装带前缀的用户唯一标识
控制群聊创建
即使两个人沟通也以订单为依据创建群聊
...
存储订单和群组关联关系，允许单个订单绑定多个群聊会话
进行业务系统通知


同一个城市所有人需要同时打开
骑手强制升级
按城市打开灰度
创建群聊会话时，Ping 开放平台后端检查如果根据账户信息查询不到用户，即用户一次未激活过，返回错误信息给中间层，不允许建立群聊会话
灰度关闭可以由中间层控制群会话创建



容量
每天 3000w 消息量
产品需求消息存储 90 天
压测
高峰期 200w 连接数
高峰期 1h 内 100w*10 条消息
详细数据后边会再重新给到
离线消息通知到业务推送服务（可用中间层做一层处理）
  离线消息多长时间内是否需要再推送到客户端，待产品再确认
  需要考虑多活场景，例如跨机房的长连接问题
  本期可以考虑由中间层禁止跨机房建立连接，让骑手和用户仍然采用电话等方式直接沟通



定位
PING SDK 对标网易云信，提供最基本的 IM 能力
处理
PING SDK 做消息去重
消息历史客户端本地存储
...
消息格式
提供最基本的文本格式
富媒体格式由业务方由 APP 自行解析消息体，类型区分开，富媒体的上传和展示由各端负责
允许保留自定义消息格式，PING 不需要理解消息体内容
PING 的系统消息和业务系统消息分开，PING SDK 需要识别系统消息并做操作，例如群增减以及消息已读回执

gather read|scatter write buf 页面 4k
