## 云服务器安装mysql
以阿里云为例：
```
apt update
apt install mysql-server
mysql
```
现在已经进入 mysql 了，下一步：
```
>>use mysql；
>>select user, plugin from mysql.user;  #root用户plugin为auth_socket，之后会出现错误
>>update mysql.user set plugin='mysql_native_password' where user='root';  #修改plugin
>>update user set host = '%' where user = 'root';  #给root用户授权使之可以在任何网络中访问
>>FLUSH PRIVILEGES;  #更新配置
>>alter user 'root'@'%' IDENTIFIED WITH mysql_native_password BY '修改的密码';  #修改密码
>>FLUSH PRIVILEGES;  #更新配置
>>exit
$service mysql restart  #重启mysql服务
```
同时，记得开启云服务器的安全组-入方向3306端口以及22端口
## Navicat连接云服务器上的mysql
[新建连接]-[ssh]处填写好公网ip以及用户、密码；同时在[常规]填写localhost以及3306，连接名随便，但是数据库用户一定是服务器上的用户名，像上面只有root就写root，这个root由于ssh换了环境就不是本地的root了，最后填好数据库用户的密码，应该没有任何问题。
