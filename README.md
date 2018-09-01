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
 
 #### 版本 2 
 
 - []优化
 - []采用dva 来写前端，并管理登录用户
 - []session 存入 redis

