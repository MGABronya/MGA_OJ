# 需求分析-Alpha版本

## 1 引言

### 1.1 编写目的

  本文档的目的是详细地介绍Dream Online Judge其Alpha版本所包含的需求，以便客户能够确认产品的确切需求以及开发人员能够根据需求设计编码。

### 1.2 背景

  该项目适为在线评测平台，接口文档可以于https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3.md查看，前端目前仍在开发中。

## 2 任务概述

### 2.1 项目概述

#### 2.1.1 项目来源及背景

  Dream Online Judge是一款针对Online Judge使用者的软件，大家可以在这里找到合适的题单、参加比赛以及分享自己的做题过程与经历。

#### 2.1.2 项目目标

  我们希望该平台是面向新人的。功能以及交互要尽可能的傻瓜式、易懂易学易用，除此之外，我们还需要满足一定的上瘾模型以此来留住目标用户。以下的设计理念请务必贯通。

**1.外部触发**(引导)

(1)对于新生引导，考虑到新生不会编程，应当优先引导新生来到QA、文章或者公告界面，除了通过这些界面介绍该平台的使用规则外，也可以为新生提供一些资源如学习引导、电子书等。

(2)在引导用户刷题时，应当优先向新生用户提供新手题单。这个题单的题目应该相当的简单(参考150，应当比150更简单更少)以让用户刷的尽兴。而对于会编程、初次接触算法的人也应当安排尽可能简单且完善的题单。

(3)在刷题时应当在题目界面提供下一题与上一题的按钮，让用户能够简单且快速的在提单内切换题目而不用回退网页等。

(4)在用户提交代码后也应该直接跳转到代码判断的页面给用户展示代码运行状态，此处应当有一个进度条可以显示代码的运行状态，尽可能简洁并美观以吊起用户的等待状态，此处也需要有回到题目、下一题的按钮来让用户简单的跳转回题目。

(5)在做题页面可以列出一个小型列表展示用户完成的题单题目完成状况，以此避免用户需要回退至题单才能查看(欣赏)自己通过题目的数量情况。

**2.内部触发**(用户情绪)

(1)归属感。给用户营造归属感需要极佳的做题体验，高颜值高自定义的界面，详细的个人信息统计等待。

(2)压抑痛感。压抑用户的恐惧而给用户一定的希望。考虑到我们的用户大多是新生，新生对新生活的恐惧大多来自于未知、悲观。我们需要在平台上给他们灌输一些刷题是正确的、明确的、健康行为的观念。我们可以在文章中引入一些但有算法题目的面经，以此来诱导用户产生如此的想法。也可以将教材上的题目写入题单，以增加刷题带来收获的感觉。

**3.自我酬劳**(自我满足)

(1)提供用户组服务，用户组+题单=表单，表单内可以看到每个用户做题情况与排名。可以尝试在做题时显示用户在用户组中的实时排名。

(2)收集用户的做题信息，为用户汇出六边形面板。

(3)引入成就/徽章系统，除了尽可能的美观、多变、有意思外，还要让初级的成就/徽章对于新人来说极其容易获取，其主要目的是增加用户的探索欲与投入欲。高级的徽章需要精致，其主要用于提醒用户自己在该平台的投入。

(4)当用户完成题目/题单后，应当在该题目/题单上标记为已完成，如填上绿色、标上绿色小勾等，要求设计的简洁干爽。

**4.社交酬劳**(被认可、尊重)

(1)提供社区服务，如QA，文章，题解，讨论区等，同时提供点赞、收藏、游览数、社区等级功能，用于满足用户的社交欲望。

(2)用户的排名、六边形面板、徽章等处了用户自己能看见，也应当可以对外展示。

(3)对于部分用户而言，实名可以帮助他们在校内更好的完成社交。但我们不能强迫用户完成实名，我们可以鼓励实名制，并让某些功能(如参加正式比赛)需要实名后才能使用，用户实名后也可以自己选择是否把实名信息展示出来。

(4)设置用户组聊天室、权限用户组等。

触发、体验、筹赏、投入

#### 2.1.3 系统功能概述

1.关于用户

(1)注册与登录：用户可以通过邮箱进行注册与登录。

(2)找回密码与更新密码：用户可以通过邮箱找回自己的密码并更新自己的密码。

(3)设置信息：用户可以查看或修改头像、博客、地址、邮箱等基础信息。

(4)查看信息：用户可以查看如基本信息以及热度等级、ac题数、ac率、六边形面板、文章总量以及列表、QA总量以及列表、题解总量以及列表、用户组列表、热度、徽章、游览、点赞、收藏、标签、竞赛分数、竞赛次数以及每场分数变化等。

(5)搜索用户：用户可以搜索其它用户。

(6)用户列表：用户可以查看ac排名、热度排名、点赞排名、收藏排名、游览排名等。

(7)用户留言：创建、查看、删除留言。

(8)用户好友：添加、通过、拒绝、删除、拉黑。

2.关于用户组

(1)管理用户组：用户可以创建、更新、查看、删除用户组信息。

(2)申请加入：用户可以申请加入某个用户组，组长也可以通过、拒绝申请。

(3)黑名单：用户组可以拉黑名单，黑名单人员无法加入用户组。

(4)点赞/踩、收藏：用户可以点赞/踩某个用户组。

(5)标签：用户可以给用户组创造、查看、删除标签。

(5)搜索用户组：用户可以搜索其它用户组。

(6)用户组列表：用户可以查看热度排名、点赞排名、收藏排名等。

(7)标准用户组：用户可以创建标准用户组。

3.关于社区

(1)文章：用户可以对文章进行增删查改点赞收藏查看游览数量并搜索、查看文章排行。

(2)QA：用户可以对QA进行增删查改点赞收藏查看游览数量并搜索、查看QA排行。

(3)题解：用户可以对题解进行增删查改查看游览数量并搜索、查看题解排行。

(4)公告：发布、订阅、查看公告。

4.关于聊天

(1)私信：用户可以与其它用户私信。

(2)群聊：用户可以在用户组中发起群聊。

5.关于题目

(1)题目：拥有权限的用户可以创建、查看、更新、删除题目。

(2)测试提交：用户可以在正式提交前进行测试提交并查看提交结果。

(3)提交：用户可以提交代码并实时查看提交状态。同时用户可以查看判题列表查看所有用户的判题状态。用户可以通过指定字段查看某些提交的状态。

(4)黑客功能：用户可以在实现了hack题目的提交中使用hack功能。

(5)题单与表单：用户可以创建、查看、更新、删除题单和表单。

## 3 功能需求

### 3.1 功能描述

#### 3.1.1 用户相关

##### 3.1.1.1 用户注册

   功能描述：用户可以通过邮箱接收验证码的方式完成账号的注册。

​		使用接口：

- /user/regist 用户注册
- /user/verify/:id 验证码获取

##### 3.1.1.2 用户登录

​		功能描述：用户可以通过邮箱以及对应密码完成登录。

   使用接口：

- /user/login 用户登录
- /user/info 返回当前登录的用户

##### 3.1.1.3 找回密码

​		功能描述：用户可以通过邮箱接收验证码的方式找回密码。

​		使用接口：

- /user/ security 找回密码
- /user/verify/:id 验证码获取

##### 3.1.1.4 更新密码

​		功能描述：用户可以使用旧密码更新密码。

​		使用接口：

- /user/update/password 更新密码

##### 3.1.1.5 修改用户信息

​		功能描述：用户可以修改自己的信息。

​		使用接口：

- /user/update 修改用户信息
- /img/upload  上传图片
- /user/label/:label 创建用户标签
- /user/label/:label 删除用户标签
- /real/name/create 创建实名
- /real/name/show/:id 查看实名
- /real/name/update 修改实名
- /real/name/delete 解除实名
- /real/name/list 查看实名列表
- /real/name/upload 上传实名表单文件

##### 3.1.1.6 查看用户信息

​		功能描述：用户可以查看自己、他人的用户信息。

​		使用接口：

- /user/info 返回当前登录的用户

- /user/show/:id 获取某个用户的所有信息

- /user/label/:id 查看用户标签

- /user/accept/rank/show/:id 获取某个用户的ac题目数量排行情况

- /user/accept/num/:id 获取某个用户的ac题目数量

- /user/score/change/:id 查看某一用户竞赛分数变化情况

- /user/hot/:id 查看某一用户今日热度数据

- /user/hot/rank/:id 查看用户热度排名

- /user/score/rank/list 查看用户竞赛分排名

  **注**：此处需要添加用户面板相关接口以及徽章相关接口。当用户获得/解锁徽章时应当有提示。

##### 3.1.1.7 查看各种用户榜单

​		功能描述：用户可以查看ac题目数量排行榜、今日点赞榜单、今日点踩榜单、今日收藏榜单、今日游览榜单、用户热度排行、用户竞赛分排行等用户榜单。

​		使用接口：

- /user/accept/rank/list 获取ac题目数量排行榜
- /user/like/rank 查看今日点赞榜单
- /user/unlike/rank 查看今日点踩榜单
- /user/collect/rank 查看今日收藏榜单
- /user/visit/rank 查看今日游览榜单
- /user/hot/rank/list 查看用户热度排行
- /user/score/rank/list 查看用户竞赛分排行

##### 3.1.1.8 用户权限管理

​		功能描述：用户可以管理比自己管理等级低的用户管理等级。

​		使用接口：

- /user/update/level/:id/:level 修改用户等级

##### 3.1.1.9 搜索用户

​		功能描述：用户可以搜索其它的用户。

​		使用接口：

- /user/search/:text 按文本搜索用户
- /user/search/label 按标签搜索用户
- /user/search/with/label/:text 按文本和标签交集搜索用户

##### 3.1.1.10 用户留言板

​		功能描述：每个用户拥有留言板界面，该界面可以用于给其它用户留言。

​		使用接口：

- /message/create/:id 创建留言
- /message/delete/:id 删除留言
- /message/list/:id 查看留言列表

##### 3.1.1.11 用户相关列表

​		功能描述：这些列表用于展示用户相关的文章、帖子、题解、AC题目、参加的用户组、收藏夹等等。

​		使用接口：

- /article/visits/:id 文章游览历史
- /article/user/list/:id 查看指定用户的文章列表
- /article/likes/:id 查看用户对文章的点赞、点踩列表
- /article/collects/:id 查看用户文章收藏夹
- /comment/user/list/:id 获取指定用户的讨论列表
- /comment/likes/:id 查看用户对讨论的点赞、点踩列表
- /group/likes/:id 查看用户对用户组的点赞、点踩列表
- /group/collects/:id 查看用户用户组收藏夹
- /group/leader/list/:id 查看某一用户创建的用户组列表
- /group/member/list/:id 查看某一用户参加的用户组列表
- /post/likes/:id 查看用户对题解的点赞、点踩列表
- /post/collects/:id 查看用户对于题解的收藏夹
- /post/visits/:id 题解游览历史
- /post/user/list/:id 查看指定用户创建的题解列表
- /problem/likes/:id 查看用户点赞、点踩题目列表
- /problem/collects/:id 查看用户题目收藏夹
- /problem/visits/:id 用户的题目游览历史
- /problem/user/list/:id 查看指定用户创建的题目列表
- /record/list 查看某类特定提交列表
- /remark/likes/:id 查看用户点赞、点踩文章回复列表
- /remark/user/list/:id 查看指定用户创建的文章回复列表
- /remark/user/list/:id 查看指定用户创建的文章回复列表
- /reply/likes/:id 查看用户点赞、点踩讨论回复列表
- /reply/user/list/:id 查看指定用户创建的讨论回复列表
- /reply/user/list/:id 查看指定用户创建的讨论回复列表
- /set/likes/:id 查看用户点赞、点踩表单列表
- /set/collects/:id 查看用户表单收藏夹
- /thread/likes/:id 查看用户点赞、点踩题解回复列表
- /thread/user/list/:id 查看指定用户创建的题解回复列表
- /thread/user/list/:id 查看指定用户创建的题解回复列表
- /topic/likes/:id 查看用户点赞、点踩的主题列表
- /topic/collects/:id 查看用户主题收藏夹
- /topic/visits/:id 用户的主题游览历史
- /topic/user/list/:id 查看指定用户创建的主题列表

##### 3.1.1.12 好友

​		功能描述：用户可以添加、删除、拉黑好友，同时还可以与好友实时聊天等。

​		使用接口：

- /friend/apply/:id 用户申请添加某个好友
- /friend/applying/list 用户查看发出的好友申请
- /friend/applied/list 用户查看接收到的好友申请
- /friend/consent/:id 用户通过好友申请
- /friend/refuse/:id 用户拒绝申请
- /friend/quit/:id 用户删除某个好友
- /friend/block/:id 用户拉黑某用户
- /friend/remove/black/:id 移除某用户的黑名单
- /friend/black/list 查看黑名单
- /letter/send/:id 创建私信
- /letter/link/list 查看获取多篇用户连接
- /letter/list/:id 列出指定用户组的聊天列表
- /letter/remove/link/:id 移除指定连接
- /letter/receive/:id 建立实时接收
- /letter/receivelink 建立连接实时接收
- /letter/block/:id 用户私信拉黑某用户
- /letter/remove/black/:id 移除某用户私信的黑名单
- /letter/black/list 查看私信黑名单
- /letter/read/:id  已读

##### 3.1.1.13 用户组

​		功能描述：用户可以添加、删除、拉黑用户组，同时还可以在用户组内实时聊天等。

​		使用接口：

- /group/create 创建用户组
- /group/show/:id 查看用户组
- /group/update/:id 更新用户组
- /group/delete/:id 删除用户组
- /group/list 查看用户组列表
- /group/apply/:id 用户申请加入某个用户组
- /group/applying/list 用户查看发出的加入用户组的申请
- /group/applied/list/:id 用户组组长查看指定组的申请
- /group/consent/:id 用户组组长通过申请
- /group/refuse/:id 用户组组长拒绝申请
- /group/quit/:id 用户退出某个用户组
- /group/block/:group/:user 用户组组长拉黑某用户
- /group/remove/black/:group/:user 移除某用户的黑名单
- /group/black/list/:id 查看黑名单
- /group/like/:id 点赞或点踩
- /group/cancel/like/:id 取消点赞或点踩
- /group/like/number/:id 查看点赞点踩数量
- /group/like/list/:id 查看点赞、点踩列表
- /group/like/show/:id 查看用户当前点赞状态
- /group/collect/:id 收藏
- /group/cancel/collect/:id 取消收藏
- /group/collect/show/:id 查看收藏状态
- /group/collect/list/:id 查看收藏列表
- /group/collect/number/:id 查看用户组被收藏数量
- /group/label/:id/:label 创建用户组标签
- /group/label/:id/:label 删除用户组标签
- /group/label/:id 查看用户组标签
- /group/search/:text 按文本搜索用户组
- /group/search/label 按标签搜索用户组
- /group/search/with/label/:text 按文本和标签交集搜索用户组
- /group/hot/rank 获取用户组热度排行
- /group/user/list/:id  看某一用户组的用户列表
- /group/standard/create/:id/:num 在用户组内生成num数量的标准用户用于标准测试
- /group/standard/list/:id 标准用户组成员信息包含账号以及密码
- /chat/send/:id 创建群聊消息
- /chat/link/list 查看获取多篇用户组连接
- /chat/list/:id 列出指定用户组的聊天列表
- /chat/remove/link/:id 移除指定连接
- /chat/receive/:id 建立实时接收
- /chat/receivelink 建立连接实时接收

#### 3.1.2 社区相关

##### 3.1.2.1 文章

​		功能描述：用户可以发布文章、查看文章、点赞收藏文章、查看文章热度、文章标签、文章分类等。

​		使用接口：

- /article/create 文章发布
- /article/show/:id 文章查看
- /article/update/:id  更新文章
- /article/delete/:id 删除文章
- /article/list 查看文章列表
- /article/like/:id 点赞或点踩文章
- /article/cancel/like/:id 取消点赞或点踩文章
- /article/like/number/:id 查看点赞点踩数量
- /article/like/list/:id 查看点赞、点踩列表
- /article/like/show/:id 查看用户当前点赞状态
- /article/collect/:id 收藏
- /article/cancel/collect/:id 取消收藏
- /article/collect/show/:id 查看收藏状态
- /article/collect/list/:id 查看收藏列表
- /article/collect/number/:id 查看文章被收藏数量
- /article/visit/:id 游览文章
- /article/visit/number/:id 游览人次
- /article/visit/list/:id 游览人次列表
- /article/search/:text 按文本搜索文章
- /article/search/label 按标签搜索文章
- /article/search/with/label/:text 按文本和标签交集搜索文章
- /article/hot/rank 获取文章热度排行
- /article/label/:id/:label 创建文章标签
- /article/label/:id/:label 删除文章标签
- /article/label/:id 查看文章标签
- /article/category/list/:id 查看指定分类的文章列表

- /category/create 创建分类
- /category/show/:id 查看分类
- /category/update/:id 更新分类
- /category/delete/:id 删除分类
- /category/list 查看分类列表

##### 3.1.2.2 题目讨论区

​		功能描述：用户可以发布讨论、查看讨论、点赞讨论、查看讨论热度等。

​		使用接口：

- /comment/create/:id 创建讨论
- /comment/show/:id 查看讨论
- /comment/update/:id 更新讨论
- /comment/delete/:id 删除讨论
- /comment/list/:id 查看讨论列表
- /comment/like/:id 点赞、点踩讨论
- /comment/cancel/like/:id 取消点赞、点踩状态
- /comment/like/number/:id 查看点赞点踩数量
- /comment/like/list/:id 查看点赞、点踩列表
- /comment/like/show/:id 查看用户当前点赞状态
- /comment/hot/rank/:id 获取讨论热度排行

##### 3.1.2.3 题解

​		功能描述：用户可以发布题解、查看题解、点赞题解、查看题解热度等。

​		使用接口：

- /post/create/:id 题解发布
- /post/show/:id 题解查看
- /post/update/:id 更新题解
- /post/delete/:id 删除题解
- /post/list/:id 查看题解列表
- /post/like/:id 点赞或点踩题解
- /post/cancel/like/:id 取消点赞或点踩题解
- /post/like/number/:id 查看点赞点踩数量
- /post/like/list/:id 查看点赞、点踩列表
- /post/like/show/:id 查看用户当前点赞状态
- /post/collect/:id 收藏
- /post/cancel/collect/:id 取消收藏
- /post/collect/show/:id 查看收藏状态
- /post/collect/list/:id 查看收藏列表
- /post/collect/number/:id 查看题解被收藏数量
- /post/visit/:id 游览题解
- /post/visit/number/:id 游览人次
- /post/visit/list/:id 游览人次列表
- /post/search/:id/:text 按文本搜索题解
- /post/search/label/:id 按标签搜索题解
- /post/search/with/label/:id/:text 按文本和标签交集搜索题解
- /post/hot/rank/:id 获取题解热度排行
- /post/label/:id/:label 创建题解标签
- /post/label/:id/:label 删除题解标签
- /post/label/:id 查看题解标签

##### 3.1.2.4 文章的回复

​		功能描述：用户可以发布文章的回复、查看文章的回复、点赞文章的回复、查看文章的回复等。

​		使用接口：

- /remark/create/:id 创建回复
- /remark/show/:id 查看回复
- /remark/update/:id 更新回复
- /remark/delete/:id 删除回复
- /remark/list/:id 查看回复列表
- /remark/like/:id 点赞或点踩回复
- /remark/cancel/like/:id 取消点赞或点踩回复
- /remark/like/number/:id 查看点赞点踩数量
- /remark/like/list/:id 查看点赞、点踩列表
- /remark/like/show/:id 查看用户当前点赞状态
- /remark/hot/rank/:id 获取回复热度排行

##### 3.1.2.5 讨论区的回复

​		功能描述：用户可以发布讨论区的回复、查看讨论区的回复、点赞讨论区的回复、查看讨论区的回复热度等。

​		使用接口：

- /reply/create/:id 创建回复
- /reply/show/:id 查看回复
- /reply/update/:id 更新回复
- /reply/delete/:id 删除回复
- /reply/list/:id 查看回复列表
- /reply/like/:id 点赞或点踩回复
- /reply/cancel/like/:id 取消点赞或点踩回复
- /reply/like/number/:id 查看点赞点踩数量
- /reply/like/list/:id 查看点赞、点踩列表
- /reply/like/show/:id 查看用户当前点赞状态
- /reply/hot/rank/:id 获取回复热度排行

##### 3.1.2.6 题解的回复

​		功能描述：用户可以发布题解的回复、查看题解的回复、点赞题解的回复、查看题解的回复热度等。

​		使用接口：

- /thread/create/:id 创建回复
- /thread/show/:id 查看回复
- /thread/update/:id 更新回复
- /thread/delete/:id 删除回复
- /thread/list/:id 查看回复列表
- /thread/like/:id 点赞或点踩回复
- /thread/cancel/like/:id 取消点赞或点踩回复
- /thread/like/number/:id 查看点赞点踩数量
- /thread/like/list/:id 查看点赞、点踩列表
- /thread/like/show/:id 查看用户当前点赞状态
- /thread/hot/rank/:id 获取回复热度排行

##### 3.1.2.7 公告

​		功能描述：用户可以发布、查看、更新、删除公告等，公告可能需要拥有一定的引导、跳转能力。

​		使用接口：

- /notice/board/create 公告发布
- /notice/board/show/:id 查看公告
- /notice/board/update/:id 更新公告
- /notice/board/delete/:id 删除公告
- /notice/board/list/:id 查看公告列表

#### 3.1.3 题目相关

##### 3.1.3.1 题目

   功能描述：用户可以发布、查看、更新、删除题目，同时，用户可以为题目创建程序以供题目的特判、标准程序、输入检测程序。除了题目的基本信息外，还包含题目是否可以黑客（是否拥有标准程序或特判程序以及输入检测程序）、是否特判、更加详细的题目提交详细。用户可以查看题目列表，同时列出每道题的ac次数、提交次数、ac率、用户是否已通过。完成题目的批量上传与题目的标签自动生成。

​		使用接口：

- /problem/create 创建题目
- /problem/show/:id 题目查看
- /problem/update/:id 更新题目
- /problem/delete/:id 删除题目
- /problem/list 查看题目列表
- /problem/like/:id 点赞或点踩题目
- /problem/cancel/like/:id 取消点赞或点踩题目
- /problem/like/number/:id 查看点赞点踩数量
- /problem/like/list/:id 查看点赞、点踩列表
- /problem/like/show/:id 查看用户当前点赞状态
- /problem/collect/:id 收藏
- /problem/cancel/collect/:id 取消收藏
- /problem/collect/show/:id 查看收藏状态
- /problem/collect/list/:id 查看收藏列表
- /problem/collect/number/:id 查看题目被收藏数量
- /problem/visit/:id 游览题目
- /problem/visit/number/:id 游览人次
- /problem/visit/list/:id 游览人次列表
- /problem/hot/rank 获取题目热度排行
- /problem/search/:text 按文本搜索题目
- /problem/search/label 按标签搜索题目
- /problem/search/with/label/:text 按文本和标签交集搜索题目
- /problem/label/:id/:label 创建题目标签
- /problem/label/:id/:label 删除题目标签
- /problem/label/:id 查看题目标签
- /problem/test/num/:id 查看题目测试样例数量
- /program/create 创建程序
- /program/show/:id 查看id指定程序
- /program/update/:id 更新程序
- /program/delete/:id 删除程序
- /program/list 查看程序列表
- /record/list 查看提交列表，可以使用这些信息做一个信息整理，如提交数量、各提交状态占比。
- /problem/create/by/text 通过xml文本创建题目
- /problem/createby/file 通过xml文件创建题目
- /tag/auto 自动生成标签
- /tag/create/:tag 创建自动标签
- /tag/delete/:tag 删除自动标签
- /tag/show/:tag 查看自动标签
- /tag/list 查看自动标签列表

##### 3.1.3.2 提交

​		功能描述：用户可以进行题目的测试提交，并查看测试提交的结果，此处需要设置帮用户自动填入测试样例的功能，并自动进行输出比对，当然用户也可以对输入进行修改。用户可以对题目进行正式提交，并能够实时查看提交的状态，除此之外，要求用户可以看到先前自己对于该题目的提交以及提交状态，并查看提交的详细信息。用户可以hack某个指定的提交。注意，提交必须标注是否可以提交hack、是否已经被hack，如果已经被hack，应当可以看到具体的hack信息。提交是否可以hack取决于提交的题目是否有输入检测程序，如果提交已经被hack，应当不允许再对该提交进行hack。

​		使用接口：

- /test/create（注意，该接口前缀与先前不同）
- /record/create 创建提交
- /record/show/:id 查看id指定提交状态
- /record/list 查看某类特定提交列表
- /record/list/case/:id 查看提交的测试通过情况
- /record/case/:id 查看某个测试的情况
- /record/publish/list 订阅提交列表
- /record/publish/:id 订阅某个提交
- /record/hack/:id 黑客指定提交
- /hack/show/:id  查看黑客

##### 3.1.3.3 题单

​		功能描述：用户可以进行创建、查看、修改、删除题单。

​		使用接口：

- /topic/create 创建主题
- /topic/show/:id 查看主题
- /topic/update/:id 更新主题
- /topic/delete/:id 删除主题
- /topic/list 查看主题列表
- /topic/like/:id 点赞、点踩主题
- /topic/cancel/like/:id 取消点赞、点踩状态
- /topic/like/number/:id 查看点赞点踩数量
- /topic/like/list/:id 查看点赞、点踩列表
- /topic/like/show/:id 查看用户当前点赞状态
- /topic/collect/:id 收藏
- /topic/cancel/collect/:id 取消收藏
- /topic/collect/show/:id 查看收藏状态
- /topic/collect/list/:id 查看收藏列表
- /topic/collect/number/:id 查看主题被收藏数量
- /topic/visit/:id 游览主题
- /topic/visit/number/:id 游览人次
- /topic/visit/list/:id 游览人次列表
- /topic/search/:text 按文本搜索主题
- /topic/search/label 按标签搜索主题
- /topic/search/with/label/:text 按文本和标签交集搜索主题
- /topic/hot/rank 获取文章热度排行
- /topic/label/:id/:label 创建主题标签
- /topic/label/:id/:label 删除主题标签
- /topic/label/:id 查看主题标签
- /topic/problem/list/:id 查看某一主题的题目列表

##### 3.1.3.4 表单

​		功能描述：用户可以进行创建、查看、修改、删除表单。表单=用户组+题单。

​		使用接口：

- /set/create 创建表单
- /set/show/:id 查看表单
- /set/update/:id 更新表单
- /set/delete/:id 删除表单
- /set/list 查看表单列表
- /set/apply/:id 用户组组长申请加入某个表单
- /set/applying/list/:id 用户组组长查看发出的表单的申请
- /set/applied/list/:id 表单创建者查看用户组的申请
- /set/consent/:id 表单创建者通过申请
- /set/refuse/:id 表单创建者拒绝申请
- /set/quit/:set/:group 用户组退出某个表单
- /set/block/:set/:group 表单拉黑某用户组
- /set/remove/black/:set/:group 移除某用户组的黑名单
- /set/black/list/:id 查看黑名单
- /set/like/:id 点赞或点踩
- /set/cancel/like/:id 取消点赞或点踩
- /set/like/number/:id 查看点赞点踩数量
- /set/like/list/:id 查看点赞、点踩列表
- /set/like/show/:id 查看用户当前点赞状态
- /set/collect/:id 收藏
- /set/cancel/collect/:id 取消收藏
- /set/collect/show/:id 查看收藏状态
- /set/collect/list/:id 查看收藏列表
- /set/collect/number/:id 查看表单被收藏数量
- /set/label/:id/:label 创建表单标签
- /set/label/:id/:label 删除表单标签
- /set/label/:id 查看表单标签
- /set/search/:text 按文本搜索表单
- /set/search/label 按标签搜索表单
- /set/search/with/label/:text 按文本和标签交集搜索表单
- /set/hot/rank 获取表单热度排行
- /set/topic/list/:id 查看某一表单的主题列表
- /set/group/list/:id 查看某一表单的用户组列表
- /set/user/group/:user/:set 查看某用户在某表单在哪一组中
- /set/rank/list/:id 查看表单内用户排行
- /set/rank/update/:id 更新表单排行

## 4 其他需求

### 4.1 验收标准

​	交付周期为**一周一交付**，交付后由测试人员依据用户故事进行功能测试，并提出改进方案以及可能存在的bug。Alpha版本不包括exam系统、competition系统以及file系统。**迭代次数不得超过6次**。

### 4.2 题目、题单等资源建设

​	在交付周期进行时，我们需要制作出至少一份面向新生的QA引导、题单以及各类资源，以及可能的面向大二学生的QA引导、题单以及给各种资源。同时我们需要为已有题目配置输入检测程序等。

### 4.3 对外联通

​	我们会尝试对以下OJ进行联通处理：UVa、POJ、HDU、CF、AtCoder。联通即获得其题目列表、题目具体信息、代码提交等接口并组装入我们的OJ中，我们将在我们的OJ中为这些接口提供热度、信息统计以及推荐系统服务等。

### 4.4 推荐系统

​	推荐系统要求实现输入一组用户的提交记录，并输出一组给用户的推荐题目。