# PNG图片格式首部字段生成
- 协议: https://www.w3.org/TR/PNG
- 可通过php.getimagesize()函数验证
- 后续利用漏洞,用php解释该图片,利用中国蚁剑等工具暴露文件系统
- 当然最简单的是: `"append some php code" >> aRealImage.png`

# 文件上传类型限制
## 1.前端限制
等于没有限制
## 2.Content-type限制
修改request的content-type即可
## 3.文件后缀名限制
修改文件后缀名即可
## 4.检查图片首部字段验证图片信息
手动构造图片首部,或在真正图片后面追加代码

> 其中,3,4都只能上传一个图片格式的文件,无论是其本身是纯php代码,或者是图片嵌入了php代码,都无法让web服务器直接用php来解释这些图片,我们需要其他的漏洞才能达到用php解释这些图片的目的