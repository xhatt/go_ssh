#### 一款简洁的能够快速使用ssh连接服务器的命令行工具

**原因**：想自己连接服务器，需要一个一个的找，然后连接，索性自己开发了一个，采用开源的框架promptui，主要是为了简洁。 一开始采用tview准备开发带UI界面的，发现需求的功能其实很简单，有点儿杀鸡用牛刀的感觉，所以采用promptui。 差不多实现了我想要的需求：

1. 快速（命令行打开可以2秒内找到想要的服务器）
2. 高效 （可以搜索）

#### 已知问题：

- 原生的promptui是不支持类似分组的功能的，所以内部的分组，其实是个假分组。这就导致了搜索功能只能搜当前页面的内容，其实promptui他只有一页。所以不能搜组内的内容。
- 配置文件需要手动维护，稍微麻烦了点，之前有考虑过通过命令添加配置或者交互式的输入，但是仔细思考了一下还是没做，有几个原因：
  - 命令加配置有学习成本，不可能每个人都能记得
  - 用交互式输入的形式输入完需要回车确认，然后比如用户名输入错了，最后才发现，就得所有的重新输一遍，体验不好
  - 我要写很多代码（🙄）
- 大家要是有什么好的建议，可以提。

![demo](./screenshot/demo.gif)

#### 安装：

1. 下载源码手动编译或者直接下载我编译好的二进制文件
2. linux或mac可自己配置一个命令别名写入环境文件中`alias ss="~/go_ssh"`

#### 使用：

1. 首次执行命令会在文件所在的目录生成一个go_ssh.yaml的配置文件，启动时可以添加`-p`参数修改配置文件名，配置文件采用yaml格式编辑。内容如下：

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

#### 操作方式：

|     键位      |                            作用                            |
| :-----------: | :--------------------------------------------------------: |
|     ↑ / ↓     |                    控制光标向上/下移动                     |
|     ← / →     |                   控制上下翻页，每页10条                   |
| a-z、A-Z、0-9 | 可直接在当前页面搜索服务器包含字段：序号、名字、用户名、IP |
|    Ctrl+C     |                          退出程序                          |
|     Enter     |                      连接选中的服务器                      |



