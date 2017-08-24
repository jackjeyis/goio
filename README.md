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
