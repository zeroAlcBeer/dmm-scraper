# better-av-tool

## 特性
- 批量抓取资料和封面的影片刮削器
- 以日站原始内容为主（DMM，MGStage，FC2等）
- 优化了生成的 nfo 部分标签以更好适配 emby

## 基本用法
- 移动程序到影片目录执行，程序会扫描查询*同目录*影片并生成影片的nfo文件以及封面
- (新)生成的nfo文件和封面路径可以自定义了
- 避免查询失败，建议配置`config.toml`中的代理地址

```toml
## 范例
[output]
# 支持输出项目 {year} {maker} {num} {actor}
# actor按第一位输出
path = 'output/{year}/{actor}'

[proxy]
## proxy [socks5://][127.0.0.1:]<1-65535>; 代理
## 协议：socks4, socks5, http, https
socket = "socks5://127.0.0.1:7891"
## 设置是否启用代理
enable = true
```

## 支持来源
- 通用番号（xxx-000）查询 DMM，搜索二次跳转获取的详情页，自动排除特典页面
- DMM 自有影片特征（xxx00000）
- 支持 MGStage 部分特征番号（000xxx-000）（siro-000）
- 支持 FC2 特征番号（fc2-000000/fc2-ppv-000000），并可查询部分已下架影片（通过 WebArchive 二次查询）
- 支持 Heyzo 特征番号
- 支持 一些自制影片（S*x Friend/S*x Syndrome）

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

## 进度
- 目前版本稳定了，没有重大bug不再更新~
