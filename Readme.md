# 目录说明
## http
用来模拟客户端请求

## proxy
用来模拟clash代理

## server
用来模拟http-rest服务器

# 流程

按照常规的来说 代理是这样的
client -> clash -> server

但是在clash中对请求做了一层包装
(left,right) := clash(net.Pipe())
client -> (clash.left->clash.right) -> server

这样就会导致内存泄漏