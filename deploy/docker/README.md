# 打开日志
参考 https://github.com/nfs-ganesha/nfs-ganesha/wiki/NFS-Ganesha-Debug-Logging

## 调试日志开启
1. 在 /export/vfs.conf 中增加下面的
```
LOG {
    COMPONENTS {
        NFS4 = FULL_DEBUG;
        ALL = DEBUG;
    }
}
```

2. 在 /proc 下面找 ganesah 进程，fd 大概在第 5 个， 通过 cat /proc/fd/cmdline 能看出来

3. 发送加载配置的信号： kill -SIGHUP fd, 上一步找到的 fd