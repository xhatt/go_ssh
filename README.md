#### 一款简洁的能够快速使用ssh连接服务器的命令行工具

![demo](./screenshot/demo.gif)

#### 安装：

1. 下载源码手动编译或者直接下载我编译好的二进制文件
2. linux或mac可自己配置一个命令别名写入环境文件中`alias ss="~/go_ssh"`

#### 使用：

1. 首次执行命令会在文件所在的目录生成一个go_ssh.yaml的配置文件，内容如下：

   ```yaml
   # 欢迎使用Go SSH 工具
   # 字段说明
   # name         ： 自定义的服务器名字 可不填
   # user         ： 服务器名 不填默认 root
   # host         ： 服务器域名或ip  ！！！必填！！！  不填的话，列表不会出现这条配置
   # port         ： 端口号  不填默认 22
   # password     ： 密码    不填默认用秘钥
   # key          ： 私钥    不填默认  ~/.ssh/id_rsa
   # passphrase   ： 私钥的密码  可不填
   # children     ： 子服务器，可不填，如果填了，这个配置会变成一个分组
   # jump         ： 跳板机 配置同上
   
   - { user: appuser, port: 22, password: 123456 }
   - { name: dev server with key path, user: appuser, host: 192.168.8.35, port: 22}
   - { name: dev server with passphrase key, user: appuser, host: 192.168.8.35, port: 22, passphrase: abcdefghijklmn}
   - { name: dev server without port, user: appuser, host: 192.168.8.35 }
   - { name: dev server without user, host: 192.168.8.35 }
   - { name: dev server without password, host: 192.168.8.35 }
   - { name: ⚡️ server with emoji name, host: 192.168.8.35 }
   - name: server with jump
     user: appuser
     host: 192.168.8.35
     port: 22
     password: 123456
     jump:
     - user: appuser
       host: 192.168.8.36
       port: 2222
   
   
   # server group 1
   - name: server group 1
     children:
     - { name: server 1, user: root, host: 192.168.1.2 }
     - { name: server 2, user: root, host: 192.168.1.3 }
     - { name: server 3, user: root, host: 192.168.1.4 }
   
   # server group 2
   - name: server group 2
     children:
     - { name: server 1, user: root, host: 192.168.2.2 }
     - { name: server 2, user: root, host: 192.168.3.3 }
     - { name: server 3, user: root, host: 192.168.4.4 }
   
   ```

2. 根据自己的需求，编写配置文件。

3. 保存之后重新执行命令即可。

