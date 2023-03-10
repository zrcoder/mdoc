---
title: 项目介绍
---

_**一魂文档**_ 是一款支持多语言的 Web 文档服务器。

## 项目初衷

项目文档工具已经是一个非常拥挤的领域，各种工具层出不穷，但却没有一个真正专注于个人尤其是开源开发者的成熟方案。无数的静态生成器、文档服务器、SaaS 产品，却没有一个能满足我们自己的实际需求，幸好我们是一群热爱编程的人，可以自己动手实现自己的需求。

由于经过多年的痛苦探索都无法在市面上找到一个合适的产品，促使了 _**一魂文档**_ 的诞生（1.0 之前的版本名为 Peach）。

下表展示了我们所关注的功能在几个主要产品之间的对比（功能点的理解上可能会出现差异）：

|产品/功能                 |_**一魂文档**_|[Mkdocs](https://www.mkdocs.org/)|[Hugo](https://gohugo.io/)|[VuePress](https://v2.vuepress.vuejs.org/)/[VitePress](https://vitepress.vuejs.org/)|[GitBook](https://www.gitbook.com/)|
|:---------------------------:|:-------------:|:----:|:--:|:----------------:|:----:|
|自托管                  | ✅ | ✅ | ✅ | ✅ | ❌ |
|多语言文档<sup>1</sup>   | ✅ | ✅ | ✅ | ✅ | ❌ |
|内置更新同步             | ✅ | ❌ | ❌ | ❌ | ✅ |
|DocSearch              | 🎯 | ❌ | ✅ | ✅ | ❌ |
|内置搜索功能             | 🎯 | ✅ | ❌ | ✅ | ✅ |
|评论系统集成             | ✅ | ❌ | ✅ | ❌ | ❌ |
|多版本                  | 🎯 | ❌ | ❌ | ❌ | ❌ |
|保护资源                | 🎯 | ❌ | ❌ | ❌ | ❌ |
|深色模式                | ✅ | ❌ | ✅ | ✅ | ❌ |
|可定制化<sup>2</sup>     | ✅ | ❌ | ✅ | ❌ | ❌ |
|语言回退<sup>3</sup>    | ✅ | ❌ | ❌ | ❌ | ❌ |

- <sup>1</sup>：目前市面上没有任何一个产品支持在不变更 URL 的情况下支持展现多语言的文档，这导致面向不同用户分享文档时需要使用不同的链接
- <sup>2</sup>：指可定制化的程度让用户无法识别时后端所使用的产品
- <sup>3</sup>：当某个文档在偏好语言中不存在时，回退显示默认语言版本的文档
- 🎯：在产品路线图中的计划功能

## 项目历史

本项目在 1.0 之前的版本名称为 Peach Docs，自 1.0 版本起已更名为 _**一魂文档**_。

项目的技术栈也从 2015 年的热门组合 [Macaron](https://go-macaron.com) 和 [Semantic UI](https://semantic-ui.com/) 升级成为最新的黄金拍档 [Flamego](https://flamego.dev) 和 [Tailwind CSS](https://tailwindcss.com/)。

项目目前也已成为 [A-SOUL 特别兴趣小组](https://github.com/asoul-sig)的一部分（之前所属于 github.com/peachdocs）。

## 开始使用

[下载安装](installation.md)或直接阅读[快速开始](quick-start.md)吧！
