# split
工具箱：切割文件工具，split



### 使用方法

```

// 切割规则，按行数=16万行
split -f log.txt -l 160000

// 切割规则，按大小=1000byte
split -f log.txt -b 1000
```

