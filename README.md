# autologin

> 用于登陆迅捷网络

> 使用go语言编写

**程序执行的逻辑：**
1. 程序需要从配置文件中读取用户名和密码，配置文件将位置是`C:\Users\***(用户名)\.config\configstore\login.json中`\\如果配置文件创建失败可以手动创建
2. 程序如果没有找到配置文件，会创建一个配置文件，并使用默认文件打开
3. 配置文件是一个json文件，需要把\*\*\*替换为指定内容，userId是用户名，password是密码
4. 在登陆过程中发生错误时，会返回服务器的消息，但只停留2s。并且会打开配置文件，要求你检查一遍
5. 成功登陆时会返回用户序列号userIndex

Release\\
[v0----点击下载](https://github.com/Yuan-byte/autologin/raw/main/login.exe)
