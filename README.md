# 简介
简单的go template, 后面模块开发根据此template 示例起服务
## 第三方使用
```
Router 框架：https://github.com/gin-gonic/gin
orm: https://github.com/jinzhu/gorm
单元测试：使用go 自带的testsing、http/httptest以及github.com/stretchr/testify
```
## 相关规范
```
规范相关参考：https://github.com/uber-go/guide/blob/master/style.md
       中文版：https://github.com/xxjwxc/uber_go_guide_cn
配置文件：
  如果服务参数比较多，可以将参数放在配置中，配置文件格式为yaml
```
### 目录结构说明
```
build：放dockerfile、配置文件、部署脚本等相关内容；
cmd: main 函数入口主要作为启动服务逻辑；
pkg:  主要功能实现逻辑都在此目录中；
  | common: 业务模块实现逻辑中比较通用的代码都可以放在此文件中，如：不同请求逻辑的handler 共用的某些方法，多个handler 需要请求一些共同的服务等都放在此文件中；
  | handler: router 中处理方法的入口，请求来到后经过router会映射到handler里面router中定义的方法；
  | model：一些request、数据库表字段、响应等结构体等定义；
  | router：注册路由、实现具体的middleware；
  | utils： 一些和业务逻辑不是非常耦合的方法实现，如client、json的coder，encoder等；
test：mock 测试具体实现在此目录中；

注: template中只列出比较典型的目录内容，如果功能模块有特殊的实现，可以根据自己需求进行相应定义；

```
## 环境说明
```
1 go 版本 >= 1.13
2 统一使用go mod 进行依赖管理
3 代理设置
  公司Goproxy使用方法：
  go env -w GO111MODULE="on"
  go env -w GOPROXY="https://goproxy.cn,direct"
  go env -w GOPRIVATE=""
  go env -w GONOSUMDB="https://goproxy.cn"
```
