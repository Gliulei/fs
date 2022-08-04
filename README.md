## fs

## description
一个用来向目的端传输文件的工具，灵感来自于scp命令经常记不住，写全效率比较慢，简单配置下这个工具可以大幅提高你的效率

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

## ROADMAP
+ 记录命令历史
+ 优化下载进度条
+ 多文件下载
