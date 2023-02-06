## fs

## description
一个用来向目的端传输文件的工具，灵感来自于scp命令对比scp命令，可以提效降本，简单配置下，这个工具就可以大幅提高你的传输效率

## 安装
### 创建 ~/.zsg/config.json配置文件，补全以后信息
```shell
{
    "user": "",
    "password": "",
    "host": "",
    "port": 22,
    "uploadDir": "",
    "downloadDir": ""
}
```

### 使用
+ fs -h - help
+ fs upload filename - 上传
+ fs download filename - 下载
+ fs use group - 使用哪个组

## ROADMAP
+ [x] 上传文件加个在当前目录下find
+ [ ] 记录命令历史
+ [ ] 多文件下载|上传
+ [ ] 配置增加|删除
