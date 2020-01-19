# better-av-tool
主要抓取日本原生内容的影片刮削器

## 基本用法
- 目前只能用命令行
- -output 指定输出 nfo , 封面目录
- -path 指定影片目录
- -proxy 指定代理地址(推荐日本)

```bash
# win下
./better-av-tool.exe -output "G:\output" -path "G:\jav" -proxy "socks5://127.0.0.1:7891"
```

## 支持来源
- 通用番号（xxx-000）查询 DMM，搜索二次跳转获取的详情页，自动排除特典页面
- 支持 MGStage 部分特征番号（000xxx-000）（siro-000）
- 支持 FC2 特征番号（fc2-000000/fc2-ppv-000000），并可查询已下架影片（通过 WebArchive 二次查询）
- 支持 Heyzo 特征番号

## 影片文件
- 会一并移动并重命名
- 有需要请注意备份

## nfo 生成
- 目前只支持 emby （kodi 理论通用，未测试）
- 固定生成为 **发行年份/番号.nfo**

## 封面下载
- 有
- fc2 为了清晰度默认抓取内容的第一张图
- dmm 来源自动切封面


## 开发中
- dlsite 来源
- 自定义命名
