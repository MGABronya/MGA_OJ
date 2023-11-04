# 需求分析-Beta版本

## 1 引言

### 1.1 编写目的

  本文档的目的是详细地介绍Dream Online Judge其Beta版本所包含的需求，以便客户能够确认产品的确切需求以及开发人员能够根据需求设计编码。

### 1.2 背景

  该项目适为在线评测平台，接口文档可以于https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3.md查看，前端目前仍在开发中。

## 2 任务概述

### 2.1 项目概述

#### 2.1.1 项目来源及背景

  Dream Online Judge是一款针对Online Judge使用者的软件，大家可以在这里找到合适的题单、参加比赛以及分享自己的做题过程与经历。

#### 2.1.2 项目目标

  我们希望该平台是面向新人的。功能以及交互要尽可能的傻瓜式、易懂易学易用，除此之外，我们还需要满足一定的上瘾模型以此来留住目标用户。Beta版本将会更加注重与比赛功能的开发。

#### 2.1.3 系统功能概述

1.关于比赛

(1)比赛的创建：用户可以创建比赛，给比赛增设密码，报名比赛等等。

(2)多种比赛：包含单人比赛、小组赛、OI赛、匹配赛、及时赛、标准赛等。

(3)搜索比赛：可以用文本、标签等搜索比赛。

(4)比赛重判：可以对比赛的提交代码进行重新判断。

(5)分数回滚：可以将比赛的结果清除。

(6)查看比赛结果：用户可以查看比赛的排名、分数、罚时、实时滚榜。

2.关于邮件

(1)发送邮件：用户可以向其他人的邮箱发送邮件，可用于群发比赛通知等。

(2)接收邮件：用户可以写反馈邮件，发送至管理者邮箱。

3.关于行为统计

(1)徽章发布：用户可以发布徽章。

(2)徽章查看：用户可以查看徽章。

## 3 功能需求

### 3.1 功能描述

#### 3.1.1 比赛相关

##### 3.1.1.1 比赛基础功能

   功能描述：用户可以完成各种比赛的增删查改，创建比赛的密码、标签、搜索比赛、对比赛进行重判回滚等。

​		使用接口：

- /competition/create/:id 创建竞赛
- /competition/show/:id 查看竞赛
- /competition/update/:id 更新竞赛
- /competition/delete/:id 删除竞赛
- /competition/list 查看竞赛列表
- /competition/passwd/create/:id 创建比赛密码
- /competition/passwd/delete/:id 删除比赛密码
- /competition/label/create/:id/:label 创建比赛标签
- /competition/label/delete/:id/:label 删除比赛标签
- /competition/label/:id 查看比赛标签
- /competition/search/:text 按文本搜索用比赛
- /competition/search/label 按标签搜索比赛
- /competition/search/with/label/:text 按文本和标签交集搜索比赛
- /competition/rejudge/:id     对指定的提交进行重新判断
- /competition/data/delete/:id 对指定比赛的结果进行清除并回滚分数
- /competition/member/rank/:competition/:member 查看指定竞赛指定成员的排名
- /competition/member/show/:competition/:member 查看指定竞赛成员罚时情况
- /competition/rank/list/:id 查看竞赛排行
- /competition/rolling/list/:id 进行滚榜

##### 3.1.1.2 个人比赛

​		功能描述：用户可以报名个人比赛。

   使用接口：

- /competition/single/submit/:id 创建提交
- /competition/single/show/record/:id 查看id指定提交状态
- /competition/single/search/list/:id 查看某类特定提交列表
- /competition/single/case/list/:id 查看提交的测试通过情况
- /competition/single/case/:id 查看某个测试的情况
- /competition/single/publish/list/:id 订阅提交列表
- /competition/single/hack/:id 黑客指定提交
- /competition/single/enter/:id 报名比赛
- /competition/single/enter/condition/:id 查看报名状态
- /competition/single/cancel/enter/:id 取消报名
- /competition/single/enter/list/:id 查看报名列表
- /competition/single/score/:id 计算比赛分数

##### 3.1.1.3 小组比赛

​		功能描述：用户可以报名小组比赛。

​		使用接口：

- /competition/group/submit/:id 创建提交
- /competition/group/show/record/:id 查看id指定提交状态
- /competition/group/search/list/:id 查看某类特定提交列表
- /competition/group/case/list/:id 查看提交的测试通过情况
- /competition/group/case/:id 查看某个测试的情况
- /competition/group/publish/list/:id 订阅提交列表
- /competition/group/hack/:id 黑客指定提交
- /competition/group/enter/:competition_id/:group_id 报名比赛
- /competition/group/enter/condition/:id 查看报名状态
- /competition/group/cancel/enter/:competition_id/:group_id 取消报名
- /competition/group/enter/list/:id 查看报名列表
- /competition/group/score/:id 计算比赛分数

##### 3.1.1.4 OI比赛

​		功能描述：用户可以报名OI比赛。

​		使用接口：

- /competition/OI/submit/:id 创建提交
- /competition/OI/show/record/:id 查看id指定提交状态
- /competition/OI/search/list/:id 查看某类特定提交列表
- /competition/OI/case/list/:id 查看提交的测试通过情况
- /competition/OI/case/:id 查看某个测试的情况
- /competition/OI/publish/list/:id 订阅提交列表
- /competition/OI/enter/:id 报名比赛
- /competition/OI/enter/condition/:id 查看报名状态
- /competition/OI/cancel/enter/:id 取消报名
- /competition/OI/enter/list/:id 查看报名列表

##### 3.1.1.5 匹配比赛

​		功能描述：用户可以报名匹配比赛。

​		使用接口：

- /competition/match/enter/:id 报名比赛
- /competition/match/enter/condition/:id 查看报名状态
- /competition/match/cancel/enter/:id 取消报名
- /competition/match/enter/list/:id 查看报名列表

##### 3.1.1.6 及时单人比赛

​		功能描述：用户报名及时个人比赛。

​		使用接口：

- /competition/random/single/enter 报名比赛
- /competition/random/single/enter/condition 查看报名状态
- /competition/random/single/cancel/enter 取消报名
- /competition/random/single/enter/list 查看报名列表
- /competition/random/single/enter/publish 实时查看报名情况

##### 3.1.1.7 及时小组比赛

​		功能描述：用户报名及时小组比赛。

​		使用接口：

- /competition/random/group/enter 报名比赛
- /competition/random/group/enter/condition 查看报名状态
- /competition/random/group/cancel/enter 取消报名
- /competition/random/group/enter/list 查看报名列表
- /competition/random/group/enter/publish 实时查看报名情况

##### 3.1.1.8 标准小组比赛

​		功能描述：用户报名标准小组比赛。

​		使用接口：

- /competition/standard/group/enter/:id 报名比赛
- /competition/standard/group/enter/condition/:id 查看报名状态
- /competition/standard/group/cancel/enter/:group_id/:competition_id 取消报名
- /competition/standard/group/enter/list/:id 查看报名列表

##### 3.1.1.9 标准个人比赛

​		功能描述：用户报名标准个人比赛。

​		使用接口：

- /competition/standard/single/enter/:id 报名比赛
- /competition/standard/single/enter/condition/:id 查看报名状态
- /competition/standard/single/cancel/enter/:id 取消报名
- /competition/standard/single/enter/list/:id 查看报名列表

#### 3.1.2 邮件相关

##### 3.1.2.1 邮件收发

​		功能描述：用户可以通过公用邮箱进行邮件的收发。

​		使用接口：

- /email/send/:id 发送邮件
- /email/receive 接收邮件

#### 3.1.3 徽章相关

##### 3.1.3.1 徽章

   功能描述：用户可以发布、查看、更新、删除徽章，同时，用户可以实时查看徽章获取情况，查看行为统计，设置徽章中的变量描述。

​		使用接口：

- /badge/create 徽章发布
- /badge/show/:id 徽章查看
- /badge/update/:id 更新徽章
- /badge/delete/:id 删除徽章
- /badge/list 查看徽章列表
- /badge/user/show/:user/:badge 查看指定用户的指定徽章
- /badge/user/list/:id 查看用户的徽章列表
- /badge/behavior/list 查看变量列表
- /badge/behavior/description/:id 查看变量描述
- /badge/publish 用户连接
- /badge/evaluate/expression/:user/:expression 查看某用户的某行为统计

##### 3.1.3.2 文件

​		功能描述：用户可以上传和下载文件。

​		使用接口：

- /file/upload 上传文件
- /file/friend 下载指定文件

#### 3.1.4 相似度查询

​	功能描述：用户可以查看多组代码的相似度

​	使用接口：

- /ngram/similarity 计算文本相似度
- /ngram/judge/:judge 计算矩阵图连通块

#### 3.1.5 心跳相关

​	功能描述：用户可以查看各个容器的心跳情况

​	使用接口：

- /heart/show/:id/:start/:end 查看指定时间段的心跳情况
- /heart/publish/:id 订阅心跳长连接
- /heart/percentage 查看近10s内的心跳忙碌占比

#### 3.1.6 chatgpt相关

​	功能描述：用户可以使用chatgpt的衍生功能

​	使用接口：

- /chatgpt/generate/code/:language 按照注释生成代码
- /chatgpt/generate/note/:language 根据代码生成注释
- /chatgpt/change/:language1/:language2 代码转换
- /chatgpt/opinion/:language 代码修改意见
- /message/ai/create 设置留言板为AI回复
- /message/ai/delete 删除ai回复模板
- /message/ai/show 查看AI回复模板
- /message/ai/update 更新AI回复模板

## 4 其他需求

### 4.1 验收标准

​	尽早交付，需轮换数个周期，交付后由测试人员依据用户故事进行功能测试，并提出改进方案以及可能存在的bug。

### 4.2 题目、题单等资源建设

​	在交付周期进行时，我们需要制作出至少一份面向新生的QA引导、题单以及各类资源，以及可能的面向大二学生的QA引导、题单以及给各种资源。同时我们需要为已有题目配置输入检测程序等。

### 4.3 云端分压

​	要求通过热更新增加承压设备。

### 4.4 推荐系统

​	推荐系统要求实现输入一组用户的提交记录，并输出一组给用户的推荐题目。

### 4.5 提交判重

​	要求对用户的提交进行相似度判断，过高相似度的提交会被标红（？）

### 4.6 降级/熔断

​	要求经过校内服务器压测后计算降级/熔断点。

### 4.7心跳检测

​	要求完成对容器的心跳检测，使容器拥有坏死后自动重启的功能。