# dmm-scraper

[![Go Report Card](https://goreportcard.com/badge/github.com/CheerChen/dmm-scraper)](https://goreportcard.com/report/github.com/CheerChen/dmm-scraper)
[![Downloads](https://img.shields.io/github/downloads/CheerChen/dmm-scraper/total.svg)](https://github.com/CheerChen/dmm-scraper/releases)
[![Release](https://img.shields.io/github/release/CheerChen/dmm-scraper.svg?label=Release)](https://github.com/CheerChen/dmm-scraper/releases)

## 特性

- 批量抓取资料和封面的影片刮削器
- 以日站原始内容为主（DMM，MGStage，FC2等）
- 优化了生成的 nfo 部分标签以更好适配 emby
- 可查询部分已下架影片（仅限FC2，Gyutto）

## 应用场景

在使用之前请对照以下场景：

- 我有大量已下载的命名混乱的电影文件，我**不在乎源文件的命名和目录**，我只想要在 Emby 中正常显示信息和封面
	-  请使用 [Emby.Plugins.JavScraper](https://github.com/JavScraper/Emby.Plugins.JavScraper)
- 我有大量已下载的命名混乱的电影文件，我**需要一次性批量采集，修改命名并且移动到指定目录**，我**无法在 Emby 中正常使用插件**
	- 请使用 [dmm-scraper](https://github.com/CheerChen/dmm-scraper/releases)
- 我希望监控 NAS 的下载，我需要下载好的电影文件**自动采集信息，修改命名并且移动到指定目录**
	- 请使用 [dmm-scraper-docker](https://github.com/CheerChen/dmm-scraper-docker)

## 基本用法

- 移动程序到影片目录执行，程序会扫描查询*同目录*影片并生成影片的nfo文件以及封面
- (新)封面可以配置是否剪裁
- 避免查询失败，建议配置`config.toml`中的代理地址

```toml
## 范例
[output]
# 支持输出项目 {year} {maker} {num} {actor}
# actor按第一位输出
path = 'output/{year}/{num}'
# 是否要切封面
needCut = true

[proxy]
## proxy [socks5://][127.0.0.1:]<1-65535>; 代理
## 协议：socks4, socks5, http, https
socket = "socks5://127.0.0.1:7891"
## 设置是否启用代理
enable = true

[DMMApi]
# 有 DMM affiliate 可以填，加快查询速度
apiId = ""
affiliateId = ""
```

## 支持来源

- 通用番号（xxx-000）依次查询 Fanza、DMM
- 支持 DMM 自有影片特征番号（xxx00000）
- 支持 MGStage 部分特征番号（000xxx-000）
- 支持 FC2 特征番号（fc2-000000/fc2-ppv-000000）
- 支持 Heyzo 特征番号
- 支持 VR 特征番号查询、大部分动画番号（GLOD，ACRN，JDXA）
- 在 Gyutto 上架的自制影片（item000000）

## 影片文件

- 会一并移动并重命名
- 有需要请注意备份

## NFO files 生成

- 按照 [kodi movie 类型的规范](https://kodi.wiki/view/NFO_files/Movies#nfo_Tags)

## Changelog

- **15 Dec 2021 (v1.4.0)** : 增加直接从 dmmapi 刮削的配置（需要开通 dmm affiliate）
- **20 Oct 2021 (v1.3.5)** : 支持 FanzaVR 刮削器
- **21 Sep 2021 (v1.3.4)** : 修复多个分片重命名覆盖问题
- **5 Aug 2021 (v1.3.3)** : 修复标题番号缺失问题；重构 Scraper 包；改良输出日志
- **2 Aug 2021 (v1.3.1)** : 修复 MGStage 查询失败问题；替换（logrus=>golog）
- **30 Jul 2021 (v1.3.0)** : 修复切封面方向；由于识别日文不稳定不再支持 Sx Syndrome 刮削
- **5 Mar 2021 (v1.2.1)** : 重构包；替换多个基础库（grab=>req, cutter=>imaging）
- **7 Aug 2020 (v1.1.0)** : 标题格式正规化为大写番号，便于 emby 搜索；提高 DMM 多个搜索结果时正确率
- **14 Jul 2020 (v1.0.1)** : 修复 DMM 查询失败的问题（需要年龄确认）
- **28 May 2020 (v1.0.0)** : 增加 Sx Syndrome\Sx Friend 的刮削器；可以自定义生成路径；第一个稳定版本
- **3 Apr 2020 (v0.9.3)** : 开始使用配置文件替换传参
- **18 Mar 2020 (v0.9.2)** : 支持 xxx00000 特征
- **16 Mar 2020 (v0.9.1)** : 修复获取多个演员的bug
- **10 Mar 2020 (v0.9.0)** : 增加 DMM 的搜索类型 `digital/videoa`
- **14 Feb 2020 (v0.8.1)** : 修复解析发片日期的bug
- **17 Jan 2020 (v0.8.0)** : 增加 Heyzo 刮削器
