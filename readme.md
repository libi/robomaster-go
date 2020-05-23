RoboMaster EP SDK (Golang) 
---

### 简介
robomasterEP SDK Golang版，可实现通过Go操控EP。初始化过程可选传入机器人IP,SDK会自动扫描局域网内存在的机器人进行连接。
可直接调用RunCmd函数发送指令并返回指令执行结果。

音视频流推送事件等连接在初始化时会根据配置自动建立连接，开发者根据需要直接使用响应的连接句柄即可。
 
### 使用说明

#### 初始化


````go
roboMasterConn, err := NewRoboMasterConn(&Option{EnableVideo: true})
````

#### 配置

```go
type Option struct {
	IP          string   //EP的IP地址 如果不传则会自动扫描
	EnableVideo bool     // 是否开启视频流
	EnableAudio bool     // 是否开启音频流
	ScanTimeout time.Duration  //未传入IP时，扫描局域网IP的超时时间
	CtrlTimeOut time.Duration  //控制指令超时时间
}
```

#### 属性



#### 调用控制指令
```go
// arg 控制指令
// return 响应结果 / 错误信息
// 该方法为同步方法
ret,err := roboMasterConn.RunCmd("robot battery ?")
```

#### 接收视频流举例

```go
for{
    buff := make([]byte,1024)
    n,err := roboMasterConn.VideoConn.Read(buff)
    fmt.Println("rec video stream",buff[0:n],err)
}
```

