## 长安链．chainmaker-docs-MR流程

### 目录
   
   一、注册流程

   二、创建分支（Fork）

   三、使用GitLab界面更新文件的MR流程

   四、使用Git命令行工具更新文件的MR流程

   五、注意事项

       【注意-1】修改页面显示语言为中文

       【注意-2】Gitlab中已fork代码与官方原项目代码同步问题

       【提示-3】关于普通MR与Draft MR

       【注意-4】提交MR，出现“Validate branches Another open merge request already exists for this source branch”报错

       【注意-5】遇见创建MR冲突，关闭无用MR需要在目标仓库close MR ，而不是在源仓库




### 一、注册流程

**1、注册页面**

在长安链平台注册页面，输入手机号码 ，点击 “获取验证码”。注册即表示同意《长安链用户使用协议》《长安链服务条款》《长安链用户隐私保护协议》，填写完成后点击 “注册” 。若已有账号，可点击 “回到登录”。

<img src="./images/contribute-to-docs1.png" alt="contribute-to-docs1.png" style="zoom: 100%;" />

**2、注册成功后，按照要求填写邮箱等信息**

<img src="./images/contribute-to-docs2.png" alt="contribute-to-docs2.png" style="zoom: 100%;" />

<img src="./images/contribute-to-docs3.png" alt="contribute-to-docs3.png" style="zoom: 100%;" />

系统会发送一封邮件到填写的邮箱，邮件中包含激活链接，点击 “Click the link below to confirm your email address” 进行邮箱激活。

<img src="./images/contribute-to-docs4.png" alt="contribute-to-docs4.png" style="zoom: 100%;" />

激活后账号方可正常使用。

<img src="./images/contribute-to-docs5.png" alt="contribute-to-docs5.png" style="zoom: 100%;" />

点击“Update profile settings”，完成信息更新。

<img src="./images/contribute-to-docs6.png" alt="contribute-to-docs6.png" style="zoom: 100%;" />

进入项目主页，可按照进行邮箱填写的完成检查，如果前面步骤均已完成，此步骤可忽略。

<img src="./images/contribute-to-docs7.png" alt="contribute-to-docs7.png" style="zoom: 100%;" />

注册完成后，若想使用 “邮箱 + 密码” 登录，需按照提示，点击 “set a password”，跳转至密码设置页面进行密码设置。设置完成后，在登录页面，输入已激活的邮箱和设置好的密码，即可登录长安链平台。

<img src="./images/contribute-to-docs8.png" alt="contribute-to-docs8.png" style="zoom: 100%;" />

注册流程到此结束，接下来可开展具体的开发任务。



### 二、创建分支（Fork）

**1、项目页面**

登录后，进入长安链项目页面，找到感兴趣的项目（如 chainmaker-docs） ，点击项目名称进入项目详情页。

<img src="./images/contribute-to-docs9.png" alt="contribute-to-docs9.png" style="zoom: 100%;" />

**2、创建分支 **

在项目详情页中，点击 “Fork” 按钮，选择要创建分支的命名空间，点击 “Select” 完成分支创建。

<img src="./images/contribute-to-docs10.png" alt="contribute-to-docs10.png" style="zoom: 100%;" />

**3、查看分支**

创建成功后，可在 “Your projects” 中查看新创建的分支。

<img src="./images/contribute-to-docs11.png" alt="contribute-to-docs11.png" style="zoom: 100%;" />



### 三、使用GitLab界面更新文件的MR流程

**1、进入项目仓库，找到要更新的文件（如 readme.md），点击文件名进入文件详情页**

<img src="./images/contribute-to-docs12.png" alt="contribute-to-docs12.png" style="zoom: 100%;" />

**2、在文件详情页中，点击 “Edit” 按钮，对文件内容进行修改**

<img src="./images/contribute-to-docs13.png" alt="contribute-to-docs13.png" style="zoom: 100%;" />

例如，将 “联系电话” 修改为 “电话”。修改完成后，在 “Commit message” 中填写修改说明（如 “Update readme.md communicate”），选择目标分支（通常为 develop），点击 “Commit changes” 提交修改 。

<img src="./images/contribute-to-docs14.png" alt="contribute-to-docs14.png" style="zoom: 100%;" />

**3、修改提交成功**

<img src="./images/contribute-to-docs15.png" alt="contribute-to-docs15.png" style="zoom: 100%;" />

**4、创建 Merge Request**

提交修改后，点击 “New merge request” 创建合并请求。

<img src="./images/contribute-to-docs16.png" alt="contribute-to-docs16.png" style="zoom: 100%;" />

在 “Source branch” 和 “Target branch” 中选择相应分支。

- **源分支 (Source Branch)**：包含待合并变更的起始分支（通常是您的功能分支或修复分支）
- **目标分支 (Target Branch)**：接收合并变更的基准分支（如发布分支）

<img src="./images/contribute-to-docs17.png" alt="contribute-to-docs17.png" style="zoom: 100%;" />

填写 “Title” 和 “Description”（如 “测试：将‘联系电话’修改为‘电话’” ）。

<img src="./images/contribute-to-docs18.png" alt="contribute-to-docs18.png" style="zoom: 100%;" />

确认无误后点击 “Compare branches and continue” 。

<img src="./images/contribute-to-docs19.png" alt="contribute-to-docs19.png" style="zoom: 100%;" />

<img src="./images/contribute-to-docs20.png" alt="contribute-to-docs20.png" style="zoom: 100%;" />

**5、合并请求详情**

在合并请求详情页中，可查看请求的详细信息，如提交者、提交时间、修改内容等。等待管理员审核并合并请求。

<img src="./images/contribute-to-docs21.png" alt="contribute-to-docs21.png" style="zoom: 100%;" />

【注意】错误：CLA未签署

<img src="./images/contribute-to-docs22.png" alt="contribute-to-docs22.png" style="zoom: 100%;" />

按照提示进入https://chainmaker.org.cn/user/cla，进行CLA签署。

***【！注意！】gitlab id和邮箱信息务必填写正确且对应 否则平台仍然会提示CLA未签署！***

<img src="./images/contribute-to-docs23.png" alt="contribute-to-docs23.png" style="zoom: 100%;" />

其中，“GitLab-ID”为下图“@”后的内容。

<img src="./images/contribute-to-docs24.png" alt="contribute-to-docs24.png" style="zoom: 100%;" />

信息填写完整后确认签署。

<img src="./images/contribute-to-docs25.png" alt="contribute-to-docs25.png" style="zoom: 100%;" />

签署成功。

<img src="./images/contribute-to-docs26.png" alt="contribute-to-docs26.png" style="zoom: 100%;" />

等待管理员合并请求。

<img src="./images/contribute-to-docs27.png" alt="contribute-to-docs27.png" style="zoom: 100%;" />

管理员合并请求。

<img src="./images/contribute-to-docs28.png" alt="contribute-to-docs28.png" style="zoom: 100%;" />

合并成功后可查看文件修改情况。

<img src="./images/contribute-to-docs29.png" alt="contribute-to-docs29.png" style="zoom: 100%;" />

合并成功。

查看修改情况：
<img src="./images/contribute-to-docs30.png" alt="contribute-to-docs30.png" style="zoom: 100%;" />

修改成功。





### 四、使用Git命令行工具更新文件的MR流程

**1、 准备工作**

**(1) Fork 主仓库**

1. 访问主仓库（如 `https://git.chainmaker.org.cn/chainmaker/chainmaker-docs`）
2. 点击 **`Fork`** 创建个人副本仓库（如 `your-username/chainmaker-docs`）

**(2) 克隆你的 Fork 仓库**

```
git clone git@git.chainmaker.org.cn:your-username/chainmaker-docs.git
cd chainmaker-docs
```

<img src="./images/contribute-to-docs31.png" alt="contribute-to-docs31.png" style="zoom: 100%;" />

【注意】需要在gitlab上事先配置SSH

配置流程：
终端获取ssh公钥：

```
cat ~/.ssh/id_rsa.pub
```

获取后复制公钥。

<img src="./images/contribute-to-docs32.png" alt="contribute-to-docs32.png" style="zoom: 100%;" />

<img src="./images/contribute-to-docs33.png" alt="contribute-to-docs33.png" style="zoom: 100%;" />

**(3) 添加上游仓库（主仓库）**

```
git remote add upstream git@git.chainmaker.org.cn:chainmaker/chainmaker-docs.git
git fetch upstream
```

<img src="./images/contribute-to-docs34.png" alt="contribute-to-docs34.png" style="zoom: 100%;" />



**2、创建功能分支**

**(1) 从主仓库的最新 `develop` 分支创建新分支**

```
git checkout -b docs/your-feature-name upstream/develop
```

- **分支命名建议**：`docs/xxx`（文档类）、`fix/xxx`（修复类）、`feat/xxx`（新功能）

<img src="./images/contribute-to-docs35.png" alt="contribute-to-docs35.png" style="zoom: 100%;" />



**3、修改代码并提交**

**(1) 修改文件**

```
vim README.md  # 或其他文件
```

<img src="./images/contribute-to-docs36.png" alt="contribute-to-docs36.png" style="zoom: 100%;" />

<img src="./images/contribute-to-docs37.png" alt="contribute-to-docs37.png" style="zoom: 100%;" />

<img src="./images/contribute-to-docs38.png" alt="contribute-to-docs38.png" style="zoom: 100%;" />

**(2) 提交更改**

```
git add .
git commit -m "fix: 更新 README 中的联系邮箱"
```

<img src="./images/contribute-to-docs39.png" alt="contribute-to-docs39.png" style="zoom: 100%;" />

- **Commit 规范**：
  - `fix:` 表示问题修复
  - `feat:` 表示新功能
  - `docs:` 表示文档更新



**4、推送到你的 Fork 仓库**

```
git push -u origin docs/your-feature-name
```

- `-u` 设置远程跟踪分支，后续可直接 `git push`

<img src="./images/contribute-to-docs40.png" alt="contribute-to-docs40.png" style="zoom: 100%;" />



**5、创建 Merge Request (MR)**

**(1) 访问 GitLab 仓库页面**

进入你的 Fork 仓库（`https://git.chainmaker.org.cn/your-username/chainmaker-docs`）

<img src="./images/contribute-to-docs41.png" alt="contribute-to-docs41.png" style="zoom: 100%;" />

**(2) 点击 `Create Merge Request`**

- **Source branch**: `docs/your-feature-name`（你的分支）
- **Target branch**: `主仓库的分支`（如v2.3.6分支）
- **填写 MR 信息**：
  - **Title**: `fix: 更新 README 中的联系邮箱`
  - **Description**: 描述修改内容（可选附加截图或测试说明）

**(3) 提交 MR**

- 等待维护者审核
- 如需修改，继续在本地提交并 `git push` 更新 MR

**(4) 删除远程分支**

```
git push origin --delete docs/your-feature-name
```

MR流程到此结束。



### 五、注意事项

#### 【注意-1】修改页面显示语言为中文

点击右上角头像→点击“Preferences”。

<img src="./images/contribute-to-docs42.png" alt="contribute-to-docs42.png" style="zoom: 100%;" />

滑动至页面底部，找到“Localization”的“Language”选项，选择熟悉的页面语言。

<img src="./images/contribute-to-docs43.png" alt="contribute-to-docs43.png" style="zoom: 100%;" />



#### 【注意-2】Gitlab中已fork代码与官方原项目代码同步问题

- 通过可视化页面同步最新项目代码

1. 在GitLab上，导航到你的Fork仓库页面。
2. 点击左侧栏中“Settings”选项卡。
3. 找到“Advanced”部分，点击“Expend”后，选择“Delete project”。
4. 重新fork分支。

- 通过Git命令行同步最新项目代码

1. 克隆你的fork到本地：

```
git clone https://github.com/your-username/your-fork.git
```

2. 确认你的remote中已经添加了上游地址，【如果已经存在上游地址，则可以跳过第4步】：

```
git remote -v
```

3. 进入克隆的目录：

```
cd your-fork
```

4. 添加原始项目的远程仓库：

```
git remote add upstream https://github.com/original-owner/original-project.git
```

5. 拉取原始项目的最新更改：

```
git fetch upstream
```

6. 切换到主分支（以main为例）：

```
git checkout main
```

7. 合并原始项目的更改到你的分支（如果有冲突，解决冲突并提交更改）：

```
git merge upstream/main
```

8. 将更新推送到你的fork：

```
git push origin main
```



#### 【提示-3】关于普通MR与Draft MR

​        Merge Request（MR）是 GitLab提供的代码审查与合并机制，用于将**一个分支的变更**合并到另一个分支（如 `main` 或 `develop`）。Draft MR（草稿合并请求）是标记为**“未完成”**的 MR，用于提前发起审查（代码未完成时获取早期反馈）、阻塞自动合并（防止误合并半成品代码）、协作讨论（团队成员可提前评论或提出建议）。

**核心规则：同一分支不能同时存在多个活跃 MR（目标分支相同）。**

Merge requests与Draft merge requests的区别：

| **特性**     | **普通 MR**          | **Draft MR**                 |
| :----------- | :------------------- | :--------------------------- |
| **合并权限** | 可直接合并           | 需手动标记为“就绪”后才能合并 |
| **UI 标识**  | 绿色合并按钮         | 灰色按钮 + "Draft" 标签      |
| **适用场景** | 代码已完成，等待合并 | 代码开发中，需提前审查       |

**如何使用 Draft MR？**

1. 创建 MR 时勾选 "Mark as draft" 复选框。
2. 或在 MR 标题前添加 Draft: 或 WIP:（如 Draft: 用户登录功能）。
3. 将 Draft MR 转为就绪状态：① 在 MR 页面点击 "Mark as ready"。② 或修改标题，移除 Draft:/WIP: 前缀。

**常见问题与注意事项**

**Q：能否对同一分支提交多个 Draft MR？**

- **允许**，但目标分支必须不同（例如：`feature/login → main` 和 `feature/login → staging`）。
- **禁止**同一分支对同一目标分支提交多个 MR（无论是 Draft 还是普通 MR）。



#### 【注意-4】提交MR，出现“Validate branches Another open merge request already exists for this source branch”报错

<img src="./images/contribute-to-docs44.png" alt="contribute-to-docs44.png" style="zoom: 100%;" />

可以看到是因为有提交的MR未审核的原因，可以和社区联络推动加快审核或关闭无效的MR。

<img src="./images/contribute-to-docs45.png" alt="contribute-to-docs45.png" style="zoom: 100%;" />

若是无用请求，点击“Close merge request”，关闭没用的MR。

<img src="./images/contribute-to-docs46.png" alt="contribute-to-docs46.png" style="zoom: 100%;" />

问题解决。



#### 【注意-5】遇见创建MR冲突，关闭无用MR需要在目标仓库close MR ，而不是在源仓库

由于不能同一分支对同一目标分支提交多个 MR，当遇到创建MR冲突时我们需要关闭之前创建的无用MR，这就需要在目标仓库关闭MR。

重复MR报错提醒：

<img src="./images/contribute-to-docs47.png" alt="contribute-to-docs47.png" style="zoom: 100%;" />

为了关闭无用MR，进入目标仓库查看：

<img src="./images/contribute-to-docs48.png" alt="contribute-to-docs48.png" style="zoom: 100%;" />

点击“Close merge request”关闭无用仓库：
<img src="./images/contribute-to-docs49.png" alt="contribute-to-docs49.png" style="zoom: 100%;" />

