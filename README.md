## 转换腾讯视频到 MP4 格式
安装: `make install`
使用方法:
`sudo -alsologtostderr=true -log_dir=log`
会在当前目录下的 out 目录生成对应的文件

腾讯视频位置
`/Users/admin/Library/Containers/com.tencent.tenvideo/Data/Library/Application Support/Download/`

## 安装依赖
`brew install ffmpeg`

## 主要命令
```
ffmpeg -i "concat:1.ts|2.ts|3.ts|4.ts|5.ts" -bsf:a aac_adtstoasc -c copy -vcodec copy 1.mp4

ffmpeg -f concat -i file.txt -c copy output.mp4
file 'path/to/file001.ts' 
file 'path/to/file002.ts'
```



