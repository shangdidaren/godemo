
## 编译
`go build -o server.exe main.go server.go`
> 注意： 这里如果不是win，不需要server.exe 直接server就可以；
> 将main.go和server.go编译为server.exe 
## 客户端请求
> git bash here 的curl 工具
`curl --http0.9 127.0.0.1:8888`
> 但是此工具不能发送套接字，所以在3.0中无法测试功能，可以使用linux的nc指令