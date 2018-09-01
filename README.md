 **由于国内某种原因，需要golang.org的包拉不下来，所以下面一种解决方案**
 ```sh
    cd ./vendor/golang.org/x
    git clone https://github.com/golang/crypto.git
    git clone https://github.com/golang/sys.git
    
 ```

 #### 技术选型
 * web : gin 
 * database: MySQL
 * config: json
 
 #### 版本 1
 
  * ✔ 完成了基本的sso 流程
  * ✔ 采用cookie存储session 的方式
 
 #### 版本 2, 正在做。。。
 
 - 各种优化
 - 采用dva 来写前端，并管理登录用户等
 - session 存入 redis
 - 用docker部署
 
 ####使用
  - 第一步：去config 中配置自己的MySQL认证信息
  - 第二步：
     ```$xslt
             $ go build -o sso .
             $./sso 
     ``` 
 #### 接入流程
 ##### 第一步：发起请求，认证服务会根据当前服务器的cookie 确认当前用户是否登录。如果登录，会生成ticket并回调redirect_uri，进入第三步，否则进入第二步
 ```
    curl http://authserver/authorize?redirect_uri=http://clientserver/callback
```
##### 第二步：用户登录。登录成功会会生成ticket并回调redirect_uri。
```$xslt
    curl http://clientserver/callback?ticket=5610fd57f6e7668eb31775620e6a671c
```
##### 第三步：clientserver 拿到tiicker 向authserver 发起校验ticket。校验通过后，clientserver,将通过后的信息写入session。整个登录过程完毕。
```$xslt
    curl -XPOST http://authserver/verifyTicket -{ticket: xxx}
``` 


