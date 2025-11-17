[toc]





# readthecods docs MD编写规范

本文档主要描述MD文档在sphinx环境下需要注意的事项

地址： <a href="https://docs.chainmaker.org.cn" target="_blank">docs.chainmaker.org.cn</a>



## 文档

**文档目录**：docs/readthedocs/docs

**图片编辑**：统一用drawio画UML等图，其他的可用ppt

**图片目录**：所有图片放在docs/readthedocs/docs/images 文件夹下



## 文档编写流程

1. clone：克隆 docs项目

2. checkout、pull： 切换并拉取到想要修改的分支

3. checkout -b；新建自己的分支

4. git push：提交代码到自己分支，等待jenkins通过后`预览网页`

5. merge request：提交合并请求，将自己分支合并到对应分支


## 编写格式规范

1. **`需在左侧目录显示的新文件`需要添加到索引`index.rst`中**

2.  **若文档在index.rst中则，标题无须手动添加序号**

   直接在总索引文件[docs/readthedocs/docs/index.rst](https://git.code.tencent.com/ChainMaker/docs/blob/readthedocs/readthedocs/docs/index.rst)中加入` :numbered: `即可自动生成索引序号
   
   索引文件中可直接写文件名不带.md/.rst等后缀

   若文档不在index.rst中，则标题序号需要手动添加
   

样例：

```rst
.. toctree::
    :maxdepth: 2
    :caption: 快速入门
    :numbered:

    tutorial/quick_start
```

3. **每个MD文档`有且仅有`一个一级标题**`且与文档名称一致`

4. **图片要用相对路径（必须用`英文`路径和命名）**

   样例如下：

```markdown
<img loading="lazy" src="../images/ManagementAccount.png" style="zoom:100%;" />
```

5. **画图需使用drawoi，图的源文件(.drawio)需放在/readthedocs/drawio文件夹下**
6. 新添加的图片需要提交UI由UI重新做图（可选）
7. 可支持标准html标签，但标准html标签内写入markdown语法可能不会被识别。
8. 为了确保图片性能达到最好
   1. 图片格式优先考虑`png`
   2. 图片尺寸分辨率不高于`1100px`，推荐`1024px`
   3. 必须：img标签增加`loading="lazy"`
   4. 必须：图片文件大小不超过200kb

## 锚点

`跳转`至当前文件中某处：通过html标签id定位

```
[跳转到：这里](#here)

<span id="here"></span>
这里是一段文字
```

`跳转`至其他文件：

```
文件A.md
    [跳转到：文件B](./文件B.html)
    
文件B.md
    我是文件B正文
```

`跳转`至其他文件中某处：

```
文件A.md
    [跳转到：哪里](./文件B.html#where)
    
文件B.md
    <span id="where"></span>
    哪里有一段文字
```

## 数学公式

**数学公式仅支持rst文档格式。**
超链接到该文档需去除后缀：如： 超链接到`math.rst `文档需写成

```
[math](./math) 
```

若为md文档则可在此网站转换为rst：
	https://cloudconvert.com/md-to-rst