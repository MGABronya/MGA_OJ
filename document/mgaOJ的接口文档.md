# mgaOJ的接口文档

## 简介

在这个文档中，你将看到mgaoj所包含的所有接口以及使用方法，文档所采用的格式为：某模型实现了某些接口，为了便于理解，这里将要用到的接口类型给列举出来：

### ApplyInterface

定义：申请

描述：实现了该接口的模型将拥有接受申请、通过申请、拒绝申请、退出等关于申请的方法

该接口包含的具体方法：

- ApplyingList	发出请求方看到的请求列表
- AppliedList      接收请求方看到的请求列表
-  Apply               发出请求
- Consent           通过请求
- Refuse              拒绝申请
- Quit                   退出/删除

### BlockInterface

定义：黑名单

描述：实现了该接口的模型将拥有拉黑、查看黑名单、移除黑名单等关于黑名单的方法

该接口包含的具体方法：

- Block		       拉黑某用户
- BlackList         查看黑名单
- RemoveBlack 移除黑名单

### CollectInterface

定义：收藏

描述：实现了该接口的模型将拥有收藏、取消收藏、查看收藏状态等关于收藏的方法

该接口包含的具体方法：

- Collect				收藏
- CancelCollect    取消收藏
-  CollectShow     查看收藏状态
- CollectList          查看收藏用户列表
- CollectNumber 查看收藏用户数量
- Collects              查看用户收藏夹

### HotInterface

定义：热度

描述：实现了该接口的模型将拥有查看热度排行的方法

该接口包含的具体方法：

- HotRanking			热度排行

### LabelInterface

定义：标签

描述：实现了该接口的模型将拥有增设标签、删除标签、查看标签等关于标签的方法

该接口包含的具体方法：

- LabelCreate	增设标签
- LabelDelete    删除标签
- LabelShow      查看标签

### LikeInterface

定义：点赞

描述：实现了该接口的模型将拥有点赞、点踩、查看点赞数量等关于点赞的方法

该接口包含的具体方法：

- Like				点赞或点踩
- CancelLike    取消点赞、点踩状态
- LikeNumber 点赞、点踩数量
- LikeList          查看点赞、点踩列表
- LikeShow      点赞状态查询
- Likes              查看用户点赞、点踩列表

### MassageInterface

定义：实时信息交流

描述：实现了该接口的模型将拥有发送消息、列出连接列表等关于信息交流的方法

该接口包含的具体方法：

- Send				发送私信
- LinkList            列出连接列表
- ChatList           列出聊天列表
- Receive            建立实时接收
- ReceiveLink     建立连接实时接收
- RemoveLink    移除某个连接

### RestInterface

定义：增删查改

描述：实现了该接口的模型将拥有增删查改的方法

该接口包含的具体方法：

- Create	增
- Update   删
- Show      查
- Delete     改
- PageList  查看列表

### SearchInterface

定义：搜索

描述：实现了该接口的模型将拥有搜索的方法

该接口包含的具体方法：

- Search 					实现文本搜索
- SearchLabel            实现标签搜索
- SearchWithLabel    实现带标签搜索

### VisitInterface

定义：游览

描述：实现了该接口的模型将拥有游览、游览人数、游览列表等关于游览的方法

该接口包含的具体方法：

- Visit					游览
- VisitNumber     游览人数
- VisitList              游览列表
- Visits                   用户游览历史记录

## 路由们

### 模型：Article

定义：文章

**基础路由：/article**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create**

    **功能：文章发布**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、category_id（可选），其中title表示文章标题，content表示文章内容，res_long表示长文本备用键值，res_short表示短文本备用键值，category_id表示分类的id(可选)。

    返回值：成功时返回创建成功相关信息和文章信息article，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：文章查看**

    **方法类型：GET**

    请求参数：文章的uuid（在接口地址的id处）。

    返回值：成功找到文章时，将会以json格式给出文章article，article中包含id,user_id,title,content,create_at,updated_at,res_short,res_long，category_id。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新文章**

    **方法类型：PUT**

    请求参数：文章的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、category_id（可选），其中title表示文章标题，content表示文章内容，res_long表示长文本备用键值，res_short表示短文本备用键值，category_id表示分类的id(可选)。

    返回值：成功更新文章时，将会以json格式给出文章article，article中包含id,user_id,title,content,create_at,updated_at,res_short,res_long，category_id。如果失败则返回失败原因。

  - **接口地址：/delete/:id**

    **功能：删除文章**

    **方法类型：DELETE**

    请求参数：文章的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功删除文章时，将会以json格式给出文章article，article中包含id,user_id,title,content,create_at,updated_at,res_short,res_long，category_id。如果失败则返回失败原因。

  - **接口地址：/list**

    **功能：查看文章列表**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：成功时，以json格式返回一个数组articles和total，articles返回了相应列表的文章信息（按照创建时间排序，越新创建排序越前），total表示文章总量，如果失败则返回失败原因。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞或点踩文章**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出文章的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞或点踩文章**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出文章的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出文章的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出文章的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回articleLikes和total，total表示点赞或者点踩的数量，articleLikes为articleLike数组，articleLike包含了user_id表示点赞用户的uid，article_id表示点赞的文章的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回articleLikes和total，total表示点赞或者点踩的数量，articleLikes为articleLike数组，articleLike包含了user_id表示点赞用户的uid，article_id表示点赞的文章的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定文章的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **CollectInterface**

  - **接口地址：/collect/:id**

    **功能：收藏**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定文章的id（即:id部分） 。

    返回值：返回收藏成功信息。

  - **接口地址：/cancel/collect/:id**

    **功能：取消收藏**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定文章的id（即:id部分） 。

    返回值：返回取消收藏成功信息。

  - **接口地址：/collect/show/:id**

    **功能：查看收藏状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定文章的id（即:id部分） 。

    返回值：返回collect，为bool类型，为true表示已经收藏，false表示未收藏。

  - **接口地址：/collect/list/:id**

    **功能：查看收藏列表**

    **方法类型：GET**

    请求参数：在接口地址中给出指定文章的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回articleCollects和total，其为articleCollect数组，包含了user_id表示收藏用户的uid，article_id表示收藏的文章的uid，create_at表示收藏时间。total表示收藏总数。

  - **接口地址：/collect/number/:id**

    **功能：查看文章被收藏数量**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定文章的id（即:id部分） 。

    返回值：返回total表示收藏人次。

  - **接口地址：/collects/:id**

    **功能：查看用户收藏夹**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回articleCollects和total，其为articleCollect数组，包含了user_id表示收藏用户的uid，article_id表示收藏的文章的uid，create_at表示收藏时间。total表示收藏总数。

- **VisitInterface**

  - **接口地址：/visit/:id**

    **功能：游览文章**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定文章的id（即:id部分） 。

    返回值：返回游览成功消息。

  - **接口地址：/visit/number/:id**

    **功能：游览人次**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定文章的id（即:id部分） 。

    返回值：返回total表示游览人次。

  - **接口地址：/visit/list/:id**

    **功能：游览人次列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定文章的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：返回articleVisits和total，total表示游览总量。articleVisits为articleVisit数组，其包含了。包含了user_id表示游览用户的uid，article_id表示游览的文章的uid，create_at表示游览时间。

  - **接口地址：/visits/:id**

    **功能：游览历史**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：返回articleVisits和total，total表示游览总量。articleVisits为articleVisit数组，其包含了。包含了user_id表示游览用户的uid，article_id表示游览的文章的uid，create_at表示游览时间。

- **SearchInterface**

  - **接口地址：/search/:text**

    **功能：按文本搜索文章**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：返回articles和total，total表示搜索到的文章总量。articles为article的数组

  - **接口地址：/search/label**

    **功能：按标签搜索文章**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回articles和total，total表示搜索到的文章总量。articles为article的数组

  - **接口地址：/search/with/label/:text**

    **功能：按文本和标签交集搜索文章**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回articles和total，total表示搜索到的文章总量。articles为article的数组

- **HotInterface**

  - **接口地址：/hot/rank**

    **功能：获取文章热度排行**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：返回articles和total，total表示文章总量。articles的每个元素包含member和score，其中member为article的uid，score为article对应的热度。已按热度排序。

- **LabelInterface**

  - **接口地址：/label/:id/:label**

    **功能：创建文章标签**

    **方法类型：POST**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定文章的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回创建成功消息

  - **接口地址：/label/:id/:label**

    **功能：删除文章标签**

    **方法类型：DELETE**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定文章的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回删除成功消息

  - **接口地址：/label/:id**

    **功能：查看文章标签**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定文章的id（即:id部分）  。

    返回值：返回articleLabels,其为articleLabel数组，每个元素包含了一个 label字符串表示标签，created_at表示创建时间，article_id表示文章uid。

- **其它**

  - **接口地址：/user/list/:id**

    **功能：查看指定用户的文章列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：成功时，以json格式返回一个数组articles和total，articles返回了相应列表的文章信息（按照创建时间排序，越新创建排序越前），total表示文章总量如果失败则返回失败原因。

  - **接口地址：/category/list/:id**

    **功能：查看指定分类的文章列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定分类的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：成功时，以json格式返回两个数组articles，articles返回了相应列表的文章信息（按照创建时间排序，越新创建排序越前），如果失败则返回失败原因。

### 模型：Category

定义：分类

**基础路由：/category**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create**（需要二级权限）

    **功能：创建分类**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含name、content、res_long(可选)、res_short（可选），其中name表示分类名称，content表示分类内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回创建成功相关信息和分类信息category，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看分类**

    **方法类型：GET**

    请求参数：分类的uuid（在接口地址的id处）。

    返回值：成功找到分类时，将会以json格式给出分类category，category中包含id,user_id,name,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

    返回值：成功找到分类时，将会以json格式给出分类category，category中包含id,user_id,name,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/update/:id**（需要二级权限）

    **功能：更新分类**

    **方法类型：PUT**

    请求参数：分类的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含name、content、res_long(可选)、res_short（可选），其中name表示分类名称，content表示分类内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功更新分类时，将会以json格式给出分类category，category中包含id,user_id,name，content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/delete/:id**（需要二级权限）

    **功能：删除分类**

    **方法类型：DELETE**

    请求参数：文章的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功删除文章时，将会以json格式给出分类category，category中包含id,user_id,name，content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/list**

    **功能：查看分类列表**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇分类，默认值为20）。

    返回值：成功时，以json格式返回一个数组categorys和total，categorys返回了相应列表的分类信息（按照创建时间排序，越早创建排序越前），total表示分类总量，如果失败则返回失败原因。

### 模型：Chat

定义：小组聊天

**基础路由：/chat**

实现的接口类型：

- **MassageInterface**

  - **接口地址：/send/:id**

    **功能：创建群聊消息**

    **方法类型：POST**

    请求参数：用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选）。

    返回值：返回创建成功消息。

  - **接口地址：/link/list**

    **功能：查看获取多篇用户组连接**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个数组chats，chats为chat数组返回了相应列表的群聊信息（按照时间排序），每个chat包含created_at表示创建时间，group_id表示用户组，author_id表示作者id，content、res_long(可选)、res_short（可选）表示内容，如果失败则返回失败原因。

  - **接口地址：/list/:id**

    **功能：列出指定用户组的聊天列表**

    **方法类型：GET**

    请求参数：用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个数组chats，chats为chat数组返回了相应列表的群聊信息（按照时间排序），每个chat包含created_at表示创建时间，group_id表示用户组，author_id表示作者id，content、res_long(可选)、res_short（可选）表示内容，如果失败则返回失败原因。

  - **接口地址：/remove/link/:id**

    **功能：移除指定连接**

    **方法类型：DELETE**

    请求参数：用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回移除成功消息。

  - **接口地址：/receive/:id**

    **功能：建立实时接收**

    **方法类型：GET**

    请求参数：用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。注意，该请求为websocket长连接。

    返回值：成功时，将持续实时返回指定组的chat，每个chat包含created_at表示创建时间，group_id表示用户组，author_id表示作者id，content、res_long(可选)、res_short（可选）表示内容，如果失败则返回失败原因。

  - **接口地址：/receivelink**

    **功能：建立连接实时接收**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。注意，该请求为websocket长连接。

    返回值：成功时，将持续实时返回包含该用户的用户组的chat作为用户组的连接请求，每个chat包含created_at表示创建时间，group_id表示用户组，author_id表示作者id，content、res_long(可选)、res_short（可选）表示内容，如果失败则返回失败原因。

### 模型：Comment

定义：题目讨论区

**基础路由：/comment**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create/:id**

    **功能：创建讨论**

    **方法类型：POST**

    请求参数：题目的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示讨论内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回创建成功相关信息和讨论信息comment，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看讨论**

    **方法类型：GET**

    请求参数：讨论的uuid（在接口地址的id处）。

    返回值：成功找到文章时，将会以json格式给出讨论comment，comment中包含id,user_id,content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新讨论**

    **方法类型：PUT**

    请求参数：讨论的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示讨论内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功更新讨论时，将会以json格式给出讨论信息comment。

  - **接口地址：/delete/:id**

    **功能：删除讨论**

    **方法类型：DELETE**

    请求参数：讨论的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功删除讨论时，返回删除成功消息。如果失败则返回失败原因。

  - **接口地址：/list/:id**

    **功能：查看讨论列表**

    **方法类型：GET**

    请求参数：题目的uuid（在接口地址的id处）。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：成功时，以json格式返回一个数组comments和total，comments返回了相应列表的讨论信息（按照创建时间排序，越早创建排序越前），total表示讨论总量，如果失败则返回失败原因。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞、点踩讨论**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出讨论的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞、点踩状态**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出讨论的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出讨论的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出文章的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回commentLikes和total，total表示点赞或者点踩的数量，commentLikes为commentLike数组，commentLike包含了user_id表示点赞用户的uid，comment_id表示点赞的讨论的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回commentLikes和total，total表示点赞或者点踩的数量，commentLikes为commentLike数组，commentLike包含了user_id表示点赞用户的uid，comment_id表示点赞的讨论的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定讨论的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **HotInterface**

  - **接口地址：/hot/rank/:id**

    **功能：获取讨论热度排行**

    **方法类型：GET**

    请求参数： 题目的uuid（在接口地址的id处）。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇讨论，默认值为20）。

    返回值：返回comments和total，total表示讨论总量。comments的每个元素包含member和score，其中member为comment的uid，score为comment对应的热度。已按热度排序。

- **其它**

  - **接口地址：/user/list/:id**

    **功能：获取指定用户的讨论列表**

    **方法类型：GET**

    请求参数：用户的uuid（在接口地址的id处）。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇讨论，默认值为20）。

    返回值：成功时，以json格式返回一个数组comments和total，comments返回了相应列表的讨论信息（按照创建时间排序，越早创建排序越前），total表示讨论总量，如果失败则返回失败原因。

### 模型：Competition

定义：竞赛

**基础路由：/competition**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create/:id**（需要二级权限）

    **功能：创建竞赛**

    **方法类型：POST**

    请求参数：表单的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含 start_time 、 end_time、type，其中type仅可为"Single"或者"Group"，表示为单人赛或者组队赛。

    返回值：成功时返回创建成功相关信息和比赛信息competition，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看竞赛**

    **方法类型：GET**

    请求参数：竞赛的uuid（在接口地址的id处）。

    返回值：成功找到竞赛时，返回成功相关信息和比赛信息competition，否则给出失败原因

  - **接口地址：/update/:id**（需要二级权限）

    **功能：更新竞赛**

    **方法类型：PUT**

    请求参数：竞赛的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含 start_time 、 end_time、type，其中type仅可为"Single"或者"Group"，表示为单人赛或者组队赛。

    返回值：成功更新竞赛时，返回成功相关信息和比赛信息competition，否则给出失败原因

  - **接口地址：/delete/:id**（需要二级权限）

    **功能：删除竞赛**

    **方法类型：DELETE**

    请求参数：竞赛的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功删除竞赛时，返回删除成功消息。如果失败则返回失败原因。

  - **接口地址：/list**

    **功能：查看竞赛列表**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：成功时，以json格式返回一个数组competitions和total，competitions返回了相应列表的竞赛信息（按照创建时间排序，越晚创建排序越前），total表示竞赛总量，如果失败则返回失败原因。

- **其它**

  - **接口地址：/member/rank/:competition/:member**

    **功能：查看指定竞赛指定成员的排名**

    **方法类型：GET**

    请求参数：竞赛的uuid（在接口地址的competition处）以及成员的uuid（在接口地址的member处）。

    返回值：成功时，以json格式返回一个整型rank，表示了当前成员的排名。

  - **接口地址：/member/show/:competition/:member**

    **功能：查看指定竞赛成员罚时情况**

    **方法类型：GET**

    请求参数：竞赛的uuid（在接口地址的competition处）以及成员的uuid（在接口地址的member处）。

    返回值：成功时，以json格式返回一个competitionMembers，其为competitionMember数组，每个元素包含了created_at、updated_at、member_id、competition_id、problem_id、penalties、condition，其中penalties为具体罚时，condition为当前题目状态。

  - **接口地址：/rank/list/:id**

    **功能：查看竞赛排行**

    **方法类型：GET**

    请求参数：竞赛的uuid（在接口地址的id处）。

    返回值：成功时，以json格式返回一个members，其为competitionRank数组，每个元素包含了created_at、updated_at、member_id、competition_id、penalties、 accept_num ，其中penalties为具体罚时， accept_num 为通过题目数量。

  - **接口地址：/rolling/list/:id**

    **功能：进行滚榜**

    **方法类型：GET**

    请求参数：竞赛的uuid（在接口地址的id处）。注意，该请求为websocket长连接。

    返回值：成功时，将持续实时返回指定比赛的rankList，每个rankList包含一个member_id表示排名发生变化的成员id，如果失败则返回失败原因。

  - **接口地址：/match/:id**

    **功能：进行随机匹配**

    **方法类型：POST**

    请求参数：竞赛的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回进入匹配组成功消息。

  - **接口地址：/unmatch/:id**

    **功能：取消随机匹配**

    **方法类型：DELETE**

    请求参数：竞赛的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回退出匹配组成功消息。

### 模型：Friend

定义：好友

**基础路由：/friend**

实现的接口类型：

- **ApplyInterface**

  - **接口地址：/apply/:id**

    **功能：用户申请添加某个好友**

    **方法类型：POST**

    请求参数：指定用户的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示申请内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时，返回申请成功消息。

  - **接口地址：/applying/list**

    **功能：用户查看发出的好友申请**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个friendApplys，其为friendApply数组，每个元素包含了id、create_at、updated_at、user_id、friend_id、condition、content、res_long、res_short。其中condition为bool类型，condition为false时表示申请被拒。friend_id表示发送目标的uid。

  - **接口地址：/applied/list**

    **功能：用户查看接收到的好友申请**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个friendApplys，其为friendApply数组，每个元素包含了id、create_at、updated_at、user_id、friend_id、condition、content、res_long、res_short。其中condition为bool类型，condition为false时表示申请被拒。friend_id表示发送目标的uid。

  - **接口地址：/consent/:id**

    **功能：用户通过好友申请**

    **方法类型：POST**

    请求参数：指定申请的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回通过成功消息。

  - **接口地址：/refuse/:id**

    **功能：用户拒绝申请**

    **方法类型：PUT**

    请求参数：指定申请的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回拒绝成功消息。

  - **接口地址：/quit/:id**

    **功能：用户删除某个好友**

    **方法类型：DELETE**

    请求参数：指定用户的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回删除成功消息。

- **BlockInterface**

  - **接口地址：/block/:id**

    **功能：用户拉黑某用户**

    **方法类型：POST**

    请求参数：指定用户的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回拉黑成功消息。

  - **接口地址：/remove/black/:id**

    **功能：移除某用户的黑名单**

    **方法类型：DELETE**

    请求参数：指定用户的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回移除黑名单成功消息。

  - **接口地址：/black/list**

    **功能：查看黑名单**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个friendBlocks，其为friendBlock数组，每个元素包含了id、create_at、updated_at、user_id、owner_id。其中user_id为被拉黑者的id，owner_id为黑名单持有者id。

### 模型：Group

定义：用户组

**基础路由：/group**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create**

    **功能：创建用户组**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、auto、users(可选)，其中title表示用户组标题，content表示用户组内容，res_long表示长文本备用键值，res_short表示短文本备用键值，auto为bool类型，表示用户组是否自动通过用户申请，users为uuid数组，表示添加这些用户进入用户组，需要二级以上权限。

    返回值：成功时返回创建成功相关信息和用户组信息group，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看用户组**

    **方法类型：GET**

    请求参数：用户组的uuid（在接口地址的id处）。

    返回值：成功找到用户组时，将会以json格式给出用户组group，group中包含title、content、res_long(可选)、res_short（可选）、auto、leader_id，其中leader_id表示用户组创建人的uuid。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新用户组**

    **方法类型：PUT**

    请求参数：用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、auto、users(可选)，其中title表示用户组标题，content表示用户组内容，res_long表示长文本备用键值，res_short表示短文本备用键值，auto为bool类型，表示用户组是否自动通过用户申请，users为uuid数组，表示添加这些用户进入用户组，需要二级以上权限。

    返回值：成功时返回更新成功相关信息和用户组信息group，否则给出失败原因

  - **接口地址：/delete/:id**

    **功能：删除用户组**

    **方法类型：DELETE**

    请求参数：用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时返回删除成功相关信息，否则给出失败原因

  - **接口地址：/list**

    **功能：查看用户组列表**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。

    返回值：成功时，以json格式返回一个数组groups和total，groups返回了相应列表的用户组信息（按照创建时间排序，越新创建排序越前），total表示用户组总量，如果失败则返回失败原因。

- **ApplyInterface**

  - **接口地址：/apply/:id**

    **功能：用户申请加入某个用户组**

    **方法类型：POST**

    请求参数：指定用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示申请内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时，返回申请成功消息。

  - **接口地址：/applying/list**

    **功能：用户查看发出的加入用户组的申请**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个groupApplys，其为groupApply数组，每个元素包含了id、create_at、updated_at、user_id、group_id、condition、content、res_long、res_short。其中condition为bool类型，condition为false时表示申请被拒。group_id表示请求的用户组的uid。

  - **接口地址：/applied/list/:id**

    **功能：用户组组长查看指定组的申请**

    **方法类型：GET**

    请求参数：指定用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个groupApplys，其为groupApply数组，每个元素包含了id、create_at、updated_at、user_id、group_id、condition、content、res_long、res_short。其中condition为bool类型，condition为false时表示申请被拒。group_id表示请求的用户组的uid。

  - **接口地址：/consent/:id**

    **功能：用户组组长通过申请**

    **方法类型：POST**

    请求参数：指定申请的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回通过成功消息。

  - **接口地址：/refuse/:id**

    **功能：用户组组长拒绝申请**

    **方法类型：PUT**

    请求参数：指定申请的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回拒绝成功消息。

  - **接口地址：/quit/:id**

    **功能：用户退出某个用户组**

    **方法类型：DELETE**

    请求参数：指定用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回删除成功消息。

- **BlockInterface**

  - **接口地址：/block/:group/:user**

    **功能：用户组组长拉黑某用户**

    **方法类型：POST**

    请求参数：指定用户的uuid（在接口地址的user处），指定用户组的uuid（在接口地址的group处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回拉黑成功消息。

  - **接口地址：/remove/black/:group/:user**

    **功能：移除某用户的黑名单**

    **方法类型：DELETE**

    请求参数：指定用户的uuid（在接口地址的user处），指定用户组的uuid（在接口地址的group处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回移除黑名单成功消息。

  - **接口地址：/black/list/:id**

    **功能：查看黑名单**

    **方法类型：GET**

    请求参数：指定用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个groupBlocks，其为groupBlock数组，每个元素包含了id、create_at、updated_at、user_id、group_id。其中user_id为被拉黑者的id，group_id为黑名单持有用户组的uuid。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞或点踩**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出用户组的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞或点踩**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出用户组的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出用户组的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出用户组的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回groupLikes和total，total表示点赞或者点踩的数量，groupLikes为groupLike数组，groupLike包含了user_id表示点赞用户的uid，group_id表示点赞的用户组的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回groupLikes和total，total表示点赞或者点踩的数量，groupLikes为groupLike数组，groupLike包含了user_id表示点赞用户的uid，group_id表示点赞的用户组的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定用户组的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **CollectInterface**

  - **接口地址：/collect/:id**

    **功能：收藏**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定用户组的id（即:id部分） 。

    返回值：返回收藏成功信息。

  - **接口地址：/cancel/collect/:id**

    **功能：取消收藏**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定用户组的id（即:id部分） 。

    返回值：返回取消收藏成功信息。

  - **接口地址：/collect/show/:id**

    **功能：查看收藏状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定用户组的id（即:id部分） 。

    返回值：返回collect，为bool类型，为true表示已经收藏，false表示未收藏。

  - **接口地址：/collect/list/:id**

    **功能：查看收藏列表**

    **方法类型：GET**

    请求参数：在接口地址中给出指定用户组的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回groupCollects和total，其为groupCollect数组，包含了user_id表示收藏用户的uid，group_id表示收藏的用户组的uid，create_at表示收藏时间。total表示收藏总数。

  - **接口地址：/collect/number/:id**

    **功能：查看用户组被收藏数量**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户组的id（即:id部分） 。

    返回值：返回total表示收藏人次。

  - **接口地址：/collects/:id**

    **功能：查看用户收藏夹**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回groupCollects和total，其为groupCollect数组，包含了user_id表示收藏用户的uid，group_id表示收藏的用户组的uid，create_at表示收藏时间。total表示收藏总数。

- **LabelInterface**

  - **接口地址：/label/:id/:label**

    **功能：创建用户组标签**

    **方法类型：POST**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定用户组的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回创建成功消息

  - **接口地址：/label/:id/:label**

    **功能：删除用户组标签**

    **方法类型：DELETE**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定用户组的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回删除成功消息

  - **接口地址：/label/:id**

    **功能：查看用户组标签**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户组的id（即:id部分）  。

    返回值：返回groupLabels,其为groupLabel数组，每个元素包含了一个 label字符串表示标签，created_at表示创建时间，group_id表示用户组的uid。

- **SearchInterface**

  - **接口地址：/search/:text**

    **功能：按文本搜索用户组**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。

    返回值：返回groups和total，total表示搜索到的文章总量。groups为group的数组

  - **接口地址：/search/label**

    **功能：按标签搜索用户组**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回groups和total，total表示搜索到的用户组总量。groups为group的数组

  - **接口地址：/search/with/label/:text**

    **功能：按文本和标签交集搜索用户组**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户组，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回groups和total，total表示搜索到的用户组总量。groups为group的数组

- **HotInterface**

  - **接口地址：/hot/rank**

    **功能：获取用户组热度排行**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。

    返回值：返回groups和total，total表示用户组总量。groups的每个元素包含member和score，其中member为group的uid，score为group对应的热度。已按热度排序。

- **其它**

  - **接口地址：/leader/list/:id**

    **功能：查看某一用户创建的用户组列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。

    返回值：成功时，以json格式返回一个数组groups和total，groups返回了相应列表的用户组信息（按照创建时间排序，越新创建排序越前），total表示用户组总量，如果失败则返回失败原因。

  - **接口地址：/member/list/:id**

    **功能：查看某一用户参加的用户组列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。

    返回值：成功时，以json格式返回一个数组groups和total，groups返回了相应列表的用户组信息（按照创建时间排序，越新创建排序越前），total表示用户组总量，如果失败则返回失败原因。

  - **接口地址：/user/list/:id**

    **功能：看某一用户组的用户列表**

    **方法类型：GET**

    请求参数：在接口地址中给出用户组的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。

    返回值：成功时，以json格式返回一个数组groups和total，groups返回了相应列表的用户组信息（按照创建时间排序，越新创建排序越前），total表示用户组总量，如果失败则返回失败原因。

### 模型：Letter

定义：私信

**基础路由：/letter**

实现的接口类型：

- **MassageInterface**

  - **接口地址：/send/:id**

    **功能：创建私信**

    **方法类型：POST**

    请求参数：指定用户的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选）。

    返回值：返回创建成功消息。

  - **接口地址：/link/list**

    **功能：查看获取多篇用户连接**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个数组letters，letters为letter数组返回了相应列表的群聊信息（按照时间排序），每个letter包含created_at表示创建时间，user_id表示接收消息的用户id，author_id表示作者id，content、res_long(可选)、res_short（可选）表示内容，read为bool值表示是否已读，如果失败则返回失败原因。

  - **接口地址：/list/:id**

    **功能：列出指定用户组的聊天列表**

    **方法类型：GET**

    请求参数：用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个数组letters，letters为letter数组返回了相应列表的群聊信息（按照时间排序），每个letter包含created_at表示创建时间，user_id表示接收消息的用户id，author_id表示作者id，content、res_long(可选)、res_short（可选）表示内容，如果失败则返回失败原因。

  - **接口地址：/remove/link/:id**

    **功能：移除指定连接**

    **方法类型：DELETE**

    请求参数：用户的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回移除成功消息。

  - **接口地址：/receive/:id**

    **功能：建立实时接收**

    **方法类型：GET**

    请求参数：用户的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。注意，该请求为websocket长连接。

    返回值：成功时，将持续实时返回指定用户的letter，每个letter包含created_at表示创建时间，user_id表示接收消息的用户id，author_id表示作者id，content、res_long(可选)、res_short（可选）表示内容，如果失败则返回失败原因。

  - **接口地址：/receivelink**

    **功能：建立连接实时接收**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。注意，该请求为websocket长连接。

    返回值：成功时，将持续实时返回包含该用户所有接收到的etter，每个letter包含created_at表示创建时间，user_id表示接收消息的用户id，author_id表示作者id，content、res_long(可选)、res_short（可选）表示内容，如果失败则返回失败原因。

- **BlockInterface**

  - **接口地址：/block/:id**

    **功能：用户私信拉黑某用户**

    **方法类型：POST**

    请求参数：指定用户的uuid（在接口地址的id处），指定用户组的uuid（在接口地址的group处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回拉黑成功消息。

  - **接口地址：/remove/black/:id**

    **功能：移除某用户私信的黑名单**

    **方法类型：DELETE**

    请求参数：指定用户的uuid（在接口地址的id处），指定用户组的uuid（在接口地址的group处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回移除黑名单成功消息。

  - **接口地址：/black/list**

    **功能：查看私信黑名单**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个letterBlocks，其为letterBlock数组，每个元素包含了id、create_at、updated_at、usera_id、userb_id。其中userb_id为被拉黑者的id。

- **其它**

  - **接口地址：/read/:id**

    **功能：已读**

    **方法类型：PUT**

    请求参数：指定私信的uuid（在接口地址的id处），Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回已读成功消息。

### 模型：Message

定义：留言板

**基础路由：/message**

实现的接口类型：

- **其它**

  - **接口地址：/create/:id**

    **功能：创建留言**

    **方法类型：POST**

    请求参数：指定用户的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content表示内容。

    返回值：成功时返回创建成功相关信息，否则给出失败原因

  - **接口地址：/delete/:id**

    **功能：删除留言**

    **方法类型：DELETE**

    请求参数：留言的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。注，删除留言的是留言板的持有者而不是留言的创建者。

    返回值：成功删除留言时，返回删除成功消息。

  - **接口地址：/list/:id**

    **功能：查看留言列表**

    **方法类型：GET**

    请求参数：指定用户的uuid（在接口地址的id处）。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：成功时，以json格式返回一个数组messages和total，messages返回了相应列表的留言信息（按照创建时间排序，越新创建排序越前），total表示留言总量，如果失败则返回失败原因。

### 模型：Post

定义：题解

**基础路由：/post**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create/:id**

    **功能：题解发布**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、其中title表示题解标题，content表示文章内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回创建成功相关信息和题解信息post，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：题解查看**

    **方法类型：GET**

    请求参数：题解的uuid（在接口地址的id处）。

    返回值：成功找到题解时，将会以json格式给出题解post，post中包含id,user_id,title,content,create_at,updated_at,res_short,res_long，problem_id。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新题解**

    **方法类型：PUT**

    请求参数：题解的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、其中title表示题解标题，content表示文章内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功更新题解时，返回更新成功消息。如果失败则返回失败原因。

  - **接口地址：/delete/:id**

    **功能：删除题解**

    **方法类型：DELETE**

    请求参数：题解的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功删除题解时，删除成功消息。如果失败则返回失败原因。

  - **接口地址：/list/:id**

    **功能：查看题解列表**

    **方法类型：GET**

    请求参数：题目的uuid（在接口地址的id处）。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇题解，默认值为20）。

    返回值：成功时，以json格式返回一个数组posts和total，posts返回了相应列表的题解信息（按照创建时间排序，越新创建排序越前），total表示题解总量，如果失败则返回失败原因。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞或点踩题解**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出题解的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞或点踩题解**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出题解的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出题解的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出题解的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回postLikes和total，total表示点赞或者点踩的数量，postLikes为postLike数组，postLike包含了user_id表示点赞用户的uid，post_id表示点赞的题解的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回postLikes和total，total表示点赞或者点踩的数量，postLikes为postLike数组，postLike包含了user_id表示点赞用户的uid，post_id表示点赞的题解的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题解的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **CollectInterface**

  - **接口地址：/collect/:id**

    **功能：收藏**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题解的id（即:id部分） 。

    返回值：返回收藏成功信息。

  - **接口地址：/cancel/collect/:id**

    **功能：取消收藏**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题解的id（即:id部分） 。

    返回值：返回取消收藏成功信息。

  - **接口地址：/collect/show/:id**

    **功能：查看收藏状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题解的id（即:id部分） 。

    返回值：返回collect，为bool类型，为true表示已经收藏，false表示未收藏。

  - **接口地址：/collect/list/:id**

    **功能：查看收藏列表**

    **方法类型：GET**

    请求参数：在接口地址中给出指定题解的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回postCollects和total，其为postCollect数组，包含了user_id表示收藏用户的uid，post_id表示收藏的题解的uid，create_at表示收藏时间。total表示收藏总数。

  - **接口地址：/collect/number/:id**

    **功能：查看题解被收藏数量**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定题解的id（即:id部分） 。

    返回值：返回total表示收藏人次。

  - **接口地址：/collects/:id**

    **功能：查看用户收藏夹**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回postCollects和total，其为postCollect数组，包含了user_id表示收藏用户的uid，post_id表示收藏的题解的uid，create_at表示收藏时间。total表示收藏总数。

- **VisitInterface**

  - **接口地址：/visit/:id**

    **功能：游览题解**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题解的id（即:id部分） 。

    返回值：返回游览成功消息。

  - **接口地址：/visit/number/:id**

    **功能：游览人次**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题解的id（即:id部分） 。

    返回值：返回total表示游览人次。

  - **接口地址：/visit/list/:id**

    **功能：游览人次列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定题解的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：返回postVisits和total，total表示游览总量。postVisits为postVisit数组，其包含了。包含了user_id表示游览用户的uid，post_id表示游览的题解的uid，create_at表示游览时间。

  - **接口地址：/visits/:id**

    **功能：游览历史**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：返回postVisits和total，total表示游览总量。postVisits为postVisit数组，其包含了。包含了user_id表示游览用户的uid，post_id表示游览的题解的uid，create_at表示游览时间。

- **SearchInterface**

  - **接口地址：/search/:id/:text**

    **功能：按文本搜索题解**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇题解，默认值为20）。

    返回值：返回posts和total，total表示搜索到的题解总量。posts为post的数组

  - **接口地址：/search/label/:id**

    **功能：按标签搜索题解**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇题解，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回posts和total，total表示搜索到的题解总量。posts为post的数组

  - **接口地址：/search/with/label/:id/:text**

    **功能：按文本和标签交集搜索题解**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回posts和total，total表示搜索到的题解总量。posts为post的数组

- **HotInterface**

  - **接口地址：/hot/rank/:id**

    **功能：获取题解热度排行**

    **方法类型：GET**

    请求参数：在接口地址中给出题目的id（即:id部分） 。 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇题解，默认值为20）。

    返回值：返回posts和total，total表示题解总量。posts的每个元素包含member和score，其中member为post的uid，score为post对应的热度。已按热度排序。

- **LabelInterface**

  - **接口地址：/label/:id/:label**

    **功能：创建题解标签**

    **方法类型：POST**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定题解的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回创建成功消息

  - **接口地址：/label/:id/:label**

    **功能：删除题解标签**

    **方法类型：DELETE**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定题解的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回删除成功消息

  - **接口地址：/label/:id**

    **功能：查看题解标签**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定题解的id（即:id部分）  。

    返回值：返回postLabels,其为postLabel数组，每个元素包含了一个 label字符串表示标签，created_at表示创建时间，post_id表示题解uid。

- **其它**

  - **接口地址：/user/list/:id**

    **功能：查看指定用户创建的题解列表**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：成功时，以json格式返回一个数组posts和total，posts返回了相应列表的题解信息（按照创建时间排序，越新创建排序越前），total表示题解总量，如果失败则返回失败原因。

### 模型：Problem

定义：题目

**基础路由：/problem**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create**（需要二级权限）

    **功能：题目发布**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、 description 、res_long(可选)、res_short（可选）、 time_limit 、 time_unit 、 memory_limit 、 memory_unit 、 input 、 output 、 sample_input 、 sample_output 、 test_input 、 test_output 、hint、 competition_id 、 source 、 special_judge ,其中title表示题目标题，description表示题目描述，res_long表示长文本备用键值，res_short表示短文本备用键值，time_limit 为uint类型，表示时间限制，time_uint表示时间单位，可为"s"或"ms"，memory_limit为uint类型，表示空间限制， memory_uint表示空间单位，可为"mb"或"kb"或"gb"，input表示输入格式，output表示输出格式、sample_input 表示输入示例、sample_output 表示输出示例，test_input 为string数组，表示测试时的每组输入，test_output 为输出数组，表示测试时的每组输出，hint表示提示，competition_id表示比赛的uid，不为空时表示为某个比赛的题目，source 表示来源，special_judge 表示特判的uid，如果不为空则表示该题为特判题目。

    返回值：成功时返回创建成功相关信息和题目信息problem，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：题目查看**

    **方法类型：GET**

    请求参数：题目的uuid（在接口地址的id处）。

    返回值：成功找到题目时，将会以json格式给出题目problem，problem中包含id,user_id,create_at,updated_at,title、 description 、res_long(可选)、res_short（可选）、 time_limit 、 time_unit 、 memory_limit 、 memory_unit 、 input 、 output 、 sample_input 、 sample_output 、 test_input 、 test_output 、hint、 competition_id 、 source 、 special_judge 。如果失败则返回失败原因。

  - **接口地址：/update/:id**（需要二级权限）

    **功能：更新题目**

    **方法类型：PUT**

    请求参数：题目的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、 description 、res_long(可选)、res_short（可选）、 time_limit 、 time_unit 、 memory_limit 、 memory_unit 、 input 、 output 、 sample_input 、 sample_output 、 test_input 、 test_output 、hint、 competition_id 、 source 、 special_judge ,其中title表示题目标题，description表示题目描述，res_long表示长文本备用键值，res_short表示短文本备用键值，time_limit 为uint类型，表示时间限制，time_uint表示时间单位，可为"s"或"ms"，memory_limit为uint类型，表示空间限制， memory_uint表示空间单位，可为"mb"或"kb"或"gb"，input表示输入格式，output表示输出格式、sample_input 表示输入示例、sample_output 表示输出示例，test_input 为string数组，表示测试时的每组输入，test_output 为输出数组，表示测试时的每组输出，hint表示提示，competition_id表示比赛的uid，不为空时表示为某个比赛的题目，source 表示来源，special_judge 表示特判的uid，如果不为空则表示该题为特判题目。

    返回值：成功更新题目时，返回更新成功消息。如果失败则返回失败原因。

  - **接口地址：/delete/:id**（需要二级权限）

    **功能：删除题目**

    **方法类型：DELETE**

    请求参数：题目的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功删除题目时，删除成功消息。如果失败则返回失败原因。

  - **接口地址：/list**

    **功能：查看题目列表**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇题目，默认值为20）。

    返回值：成功时，以json格式返回一个数组problems和total，problems返回了相应列表的题目信息（按照创建时间排序，越新创建排序越前），total表示题目总量，如果失败则返回失败原因。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞或点踩题目**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出题目的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞或点踩题目**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出题目的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出题目的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出题目的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回problemLikes和total，total表示点赞或者点踩的数量，problemLikes为problemLike数组，problemLike包含了user_id表示点赞用户的uid，problem_id表示点赞的题目的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回problemLikes和total，total表示点赞或者点踩的数量，problemLikes为problemLike数组，problemLike包含了user_id表示点赞用户的uid，problem_id表示点赞的题目的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题目的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **CollectInterface**

  - **接口地址：/collect/:id**

    **功能：收藏**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题目的id（即:id部分） 。

    返回值：返回收藏成功信息。

  - **接口地址：/cancel/collect/:id**

    **功能：取消收藏**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题目的id（即:id部分） 。

    返回值：返回取消收藏成功信息。

  - **接口地址：/collect/show/:id**

    **功能：查看收藏状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题目的id（即:id部分） 。

    返回值：返回collect，为bool类型，为true表示已经收藏，false表示未收藏。

  - **接口地址：/collect/list/:id**

    **功能：查看收藏列表**

    **方法类型：GET**

    请求参数：在接口地址中给出指定题目的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回problemCollects和total，其为problemCollect数组，包含了user_id表示收藏用户的uid，problem_id表示收藏的题目的uid，create_at表示收藏时间。total表示收藏总数。

  - **接口地址：/collect/number/:id**

    **功能：查看题目被收藏数量**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定题目的id（即:id部分） 。

    返回值：返回total表示收藏人次。

  - **接口地址：/collects/:id**

    **功能：查看用户收藏夹**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回problemCollects和total，其为problemCollect数组，包含了user_id表示收藏用户的uid，problem_id表示收藏的题目的uid，create_at表示收藏时间。total表示收藏总数。

- **VisitInterface**

  - **接口地址：/visit/:id**

    **功能：游览题目**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题目的id（即:id部分） 。

    返回值：返回游览成功消息。

  - **接口地址：/visit/number/:id**

    **功能：游览人次**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定题目的id（即:id部分） 。

    返回值：返回total表示游览人次。

  - **接口地址：/visit/list/:id**

    **功能：游览人次列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定题目的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：返回problemVisits和total，total表示游览总量。problemVisits为problemVisit数组，其包含了。包含了user_id表示游览用户的uid，problem_id表示游览的题目的uid，create_at表示游览时间。

  - **接口地址：/visits/:id**

    **功能：游览历史**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：返回problemVisits和total，total表示游览总量。problemVisits为problemVisit数组，其包含了。包含了user_id表示游览用户的uid，problem_id表示游览的题目的uid，create_at表示游览时间。

- **HotInterface**

  - **接口地址：/hot/rank**

    **功能：获取表单热度排行**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个题目，默认值为20）。

    返回值：返回problems和total，total表示题目总量。problems的每个元素包含member和score，其中member为problem的uid，score为problem对应的热度。已按热度排序。

- **SearchInterface**

  - **接口地址：/search/:id/:text**

    **功能：按文本搜索题解**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇题目，默认值为20）。

    返回值：返回problems和total，total表示搜索到的题目总量。problems为problem的数组

  - **接口地址：/search/label/:id**

    **功能：按标签搜索题目**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇题目，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回problems和total，total表示搜索到的题目总量。problems为problem的数组

  - **接口地址：/search/with/label/:id/:text**

    **功能：按文本和标签交集搜索题目**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇题目，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回problems和total，total表示搜索到的题目总量。problems为problem的数组

- **LabelInterface**

  - **接口地址：/label/:id/:label**

    **功能：创建题目标签**

    **方法类型：POST**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定题目的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回创建成功消息

  - **接口地址：/label/:id/:label**

    **功能：删除题目标签**

    **方法类型：DELETE**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定题目的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回删除成功消息

  - **接口地址：/label/:id**

    **功能：查看题目标签**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定题解的id（即:id部分）  。

    返回值：返回problemLabels,其为problemLabel数组，每个元素包含了一个 label字符串表示标签，created_at表示创建时间，problem_id表示题目uid。

- **其它**

  - **接口地址：/user/list/:id**

    **功能：查看指定用户创建的题目列表**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：成功时，以json格式返回一个数组problems和total，problems返回了相应列表的题目信息（按照创建时间排序，越新创建排序越前），total表示题目总量，如果失败则返回失败原因。
    
  - **接口地址：/test/num/:id**
  
    **功能：查看题目测试样例数量**
  
    **方法类型：GET**
  
    请求参数：  在接口地址中给出指定题目的id（即:id部分）  。
  
    返回值：成功时，以json格式返回一个total，total表示测试样例数量，如果失败则返回失败原因。

### 模型：Record

定义：代码提交

**基础路由：/record**

实现的接口类型：

- **其它**

  - **接口地址：/create**

    **功能：创建提交**

    **方法类型：POST**

    请求参数：  Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含language、code、problem_id ,其中language表示语言，code表示提交的代码，problem_id表示题目id。这里的language支持如下："C"、"C#"、"C++"、"C++11"、"Erlang"、"Go"、"Java"、"JavaScript"、"Kotlin"、"Pascal"、"PHP"、"Python"、"Racket"、"Ruby"、"Rust"、"Scala"、 "Swift"

    返回值：成功时，返回成功消息，如果失败则返回失败原因。

  - **接口地址：/show/:id**

    **功能：查看id指定提交状态**

    **方法类型：GET**

    请求参数：  Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定提交的id（即:id部分）  。

    返回值：成功时，以json格式返回一个record，record包含language、code、problem_id、 created_at 、updated_at、 user_id 、condition、competition_id，其中condition表示提交状态，提交状态包含Waiting（等待）、Competition hasn't Started（比赛未开始）、Input Doesn't Exist（输入在数据库中不存在）、Output Doesn't Exist（输入在数据库中不存在）、System Error 1（服务器问题：创建文件失败）、System Error 2（服务器问题：编译指令执行失败）、Compile timeout（编译超时）、Compile Error（编译错误）、System Error 3（服务器问题：消息管道创建失败）、System Error 4（服务器问题：运行指令执行失败）、Time Limit Exceeded（超出时间限制）、Runtime Error（运行时错误）、Memory Limit Exceeded（超出空间限制）、Wrong Answer（错误答案）、System error 5（服务器问题：数据库插入数据失败）、Accepted（提交通过）、Absent from the race（未参加比赛）、Luanguage Error（语言错误）。如果失败则返回失败原因。

  - **接口地址：/list**

    **功能：查看某类特定提交列表**

    **方法类型：GET**

    请求参数：  Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇提交，默认值为20），language表示使用的语言（默认为空），user_id表示提交用户的id（默认为空），problem_id表示题目的id（默认为空），start_time表示在这之后的提交（默认为空），end_time表示在这之前的提交（默认为空），condition表示提交状态（默认为空），pass_low表示提交最少通过多少测试（默认为空），pass_top表示提交至多通过多少测试（默认为空），competition_id表示隶属某个比赛的提交（默认为空）。

    返回值：成功时，以json格式返回复合条件的records和total，如果失败则返回失败原因。

  - **接口地址：/list/case/:id**

    **功能：查看提交的测试通过情况**

    **方法类型：GET**

    请求参数：在接口地址中给出指定提交的id（即:id部分）  。

    返回值：成功时，以json格式返回一个cases，cases为case数组，每个case含有record_id表示为哪一个提交的测试通过情况，id表示为第几个测试，time表示测试使用时间，memory表示测试使用空间，如果失败则返回失败原因。

  - **接口地址：/case/:id**

    **功能：查看某个测试的情况**

    **方法类型：GET**

    请求参数：在接口地址中给出指定提交的id（即:id部分）  

    返回值：成功时，以json格式返回一个case

### 模型：Rejudge

定义：对提交代码的重新判断

**基础路由：/rejudge**

实现的接口类型：

- **其它**

  - **接口地址：/do**（需要二级权限）

    **功能：进行重判**

    **方法类型：PUT**

    请求参数：在Params处提供，language表示使用的语言（默认为空），user_id表示提交用户的id（默认为空），problem_id表示题目的id（默认为空），start_time表示在这之后的提交（默认为空），end_time表示在这之前的提交（默认为空），condition表示提交状态（默认为空），competition_id表示隶属某个比赛的提交（默认为空）。  

    返回值：成功时，重判成功消息

  - **接口地址：/competiton/delete/:id**（需要二级权限）

    **功能：对某场比赛结果进行清空**

    **方法类型：DELETE**

    请求参数：在接口地址中给出指定比赛的id（即:id部分）  

    返回值：成功时，返回成功消息，否则返回失败原因。注意，此操作需要在重判操作之前。

  - **接口地址：/competiton/score/:id**（需要二级权限）

    **功能：对某场比赛结果重新进行分数统计**

    **方法类型：PUT**

    请求参数：在接口地址中给出指定比赛的id（即:id部分）  

    返回值：成功时，返回成功消息，否则返回失败原因。注意，此操作需要在重判操作之后。

### 模型：Remark

定义：文章的回复

**基础路由：/remark**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create/:id**

    **功能：创建回复**

    **方法类型：POST**

    请求参数：文章的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示回复内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回创建成功相关信息，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看回复**

    **方法类型：GET**

    请求参数：回复的uuid（在接口地址的id处）。

    返回值：成功找到回复时，将会以json格式给出回复remark，remark中包含id,user_id,article_id,content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新回复**

    **方法类型：PUT**

    请求参数：回复的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示回复内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回更新成功相关信息，否则给出失败原因

  - **接口地址：/delete/:id**

    **功能：删除回复**

    **方法类型：DELETE**

    请求参数：回复的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时返回删除成功相关信息，否则给出失败原因

  - **接口地址：/list/:id**

    **功能：查看回复列表**

    **方法类型：GET**

    请求参数：文章的uuid（在接口地址的id处）。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇回复，默认值为20）。

    返回值：成功时，以json格式返回一个数组remarks和total，remarks返回了相应列表的回复信息（按照创建时间排序，越新创建排序越前），total表示回复总量，如果失败则返回失败原因。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞或点踩回复**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞或点踩回复**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出回复的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回remarkLikes和total，total表示点赞或者点踩的数量，remarkLikes为remarkLike数组，remarkLike包含了user_id表示点赞用户的uid，remark_id表示点赞的回复的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回remarkLikes和total，total表示点赞或者点踩的数量，remarkLikes为remarkLike数组，remarkLike包含了user_id表示点赞用户的uid，remark_id表示点赞的回复的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定回复的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **HotInterface**

  - **接口地址：/hot/rank/:id**

    **功能：获取回复热度排行**

    **方法类型：GET**

    请求参数：在接口地址中给出文章的id（即:id部分） 。 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇回复，默认值为20）。

    返回值：返回remarks和total，total表示回复总量。remarks的每个元素包含member和score，其中member为remark的uid，score为remark对应的热度。已按热度排序。

- **其它**

  - **接口地址：/user/list/:id**

    **功能：查看指定用户创建的回复列表**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：成功时，以json格式返回一个数组remarks和total，remarks返回了相应列表的回复信息（按照创建时间排序，越新创建排序越前），total表示回复总量，如果失败则返回失败原因。

### 模型：Reply

定义：讨论的回复

**基础路由：/reply**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create/:id**

    **功能：创建回复**

    **方法类型：POST**

    请求参数：讨论的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示回复内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回创建成功相关信息，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看回复**

    **方法类型：GET**

    请求参数：回复的uuid（在接口地址的id处）。

    返回值：成功找到回复时，将会以json格式给出回复remark，remark中包含id,user_id,comment_id,content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新回复**

    **方法类型：PUT**

    请求参数：回复的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示回复内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回更新成功相关信息，否则给出失败原因

  - **接口地址：/delete/:id**

    **功能：删除回复**

    **方法类型：DELETE**

    请求参数：回复的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时返回删除成功相关信息，否则给出失败原因

  - **接口地址：/list/:id**

    **功能：查看回复列表**

    **方法类型：GET**

    请求参数：讨论的uuid（在接口地址的id处）。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇回复，默认值为20）。

    返回值：成功时，以json格式返回一个数组remarks和total，remarks返回了相应列表的回复信息（按照创建时间排序，越新创建排序越前），total表示回复总量，如果失败则返回失败原因。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞或点踩回复**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞或点踩回复**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出回复的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回replyLikes和total，total表示点赞或者点踩的数量，replyLikes为replyLike数组，replyLike包含了user_id表示点赞用户的uid，reply_id表示点赞的回复的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回replyLikes和total，total表示点赞或者点踩的数量，replyLikes为replyLike数组，replyLike包含了user_id表示点赞用户的uid，reply_id表示点赞的回复的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定回复的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **HotInterface**

  - **接口地址：/hot/rank/:id**

    **功能：获取回复热度排行**

    **方法类型：GET**

    请求参数：在接口地址中给出讨论的id（即:id部分） 。 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇回复，默认值为20）。

    返回值：返回remarks和total，total表示回复总量。remarks的每个元素包含member和score，其中member为remark的uid，score为remark对应的热度。已按热度排序。

- **其它**

  - **接口地址：/user/list/:id**

    **功能：查看指定用户创建的回复列表**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：成功时，以json格式返回一个数组replys和total，replys返回了相应列表的回复信息（按照创建时间排序，越新创建排序越前），total表示回复总量，如果失败则返回失败原因。

### 模型：Set

定义：表单

**基础路由：/set**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create**

    **功能：创建表单**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、 topics 、  groups(可选)、 auto_update、 auto_pass 、 pass_num 、 pass_re，其中title表示表单标题，content表示表单内容，res_long表示长文本备用键值，res_short表示短文本备用键值，auto_update为bool类型表示是否每小时更新排名、auto_pass 为bool类型表示是否自动通过用户组申请、pass_num 为int类型表示每组的最大成员数量限制、pass_re为bool类型表示小组成员是否可以重复，topics 表示主题的id数组，表示添加这些主题进入表单，groups表示用户组的id数组，表示添加这些用户组进入表单，需要二级以上权限。

    返回值：成功时返回创建成功相关信息，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看表单**

    **方法类型：GET**

    请求参数：表单的uuid（在接口地址的id处）。

    返回值：成功找到表单时，将会以json格式给出表单set，set中包含id、user_id、ttitle、content、res_long(可选)、res_short（可选）、 auto_update、 auto_pass 、 pass_num 、 pass_re，其中user_id表示表单创建人的uuid。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新表单**

    **方法类型：PUT**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、 topics 、  groups(可选)、 auto_update、 auto_pass 、 pass_num 、 pass_re，其中title表示表单标题，content表示表单内容，res_long表示长文本备用键值，res_short表示短文本备用键值，auto_update为bool类型表示是否每小时更新排名、auto_pass 为bool类型表示是否自动通过用户组申请、pass_num 为int类型表示每组的最大成员数量限制、pass_re为bool类型表示小组成员是否可以重复，topics 表示主题的id数组，表示添加这些主题进入表单，groups表示用户组的id数组，表示添加这些用户组进入表单，需要二级以上权限。

    返回值：成功时返回更新成功相关信息，否则给出失败原因

  - **接口地址：/delete/:id**

    **功能：删除表单**

    **方法类型：DELETE**

    请求参数：表单的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时返回删除成功相关信息，否则给出失败原因

  - **接口地址：/list**

    **功能：查看表单列表**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个表单，默认值为20）。

    返回值：成功时，以json格式返回一个数组sets和total，sets返回了相应列表的用户组信息（按照创建时间排序，越新创建排序越前），total表示表单总量，如果失败则返回失败原因。

- **ApplyInterface**

  - **接口地址：/apply/:id**

    **功能：用户组组长申请加入某个表单**

    **方法类型：POST**

    请求参数：指定表单的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token，在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选）， group_id ，其中content表示申请内容，res_long表示长文本备用键值，res_short表示短文本备用键值， group_id 为需要加入的用户组id，用户必须为该用户组的组长。

    返回值：成功时，返回申请成功消息。

  - **接口地址：/applying/list/:id**

    **功能：用户组组长查看发出的表单的申请**

    **方法类型：GET**

    请求参数：指定用户组的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个setApplys，其为setApply数组，每个元素包含id、created_at、updated_at、condition、content、res_long(可选)、res_short（可选）， group_id ，其中content表示申请内容，res_long表示长文本备用键值，res_short表示短文本备用键值， group_id 为需要加入的用户组id，用户必须为该用户组的组长，condition为bool类型，为false时表示申请被拒。

  - **接口地址：/applied/list/:id**

    **功能：表单创建者查看用户组的申请**

    **方法类型：GET**

    请求参数：指定表单的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个setApplys，其为setApply数组，每个元素包含id、created_at、updated_at、condition、content、res_long(可选)、res_short（可选）， group_id ，其中content表示申请内容，res_long表示长文本备用键值，res_short表示短文本备用键值， group_id 为需要加入的用户组id，用户必须为该用户组的组长，condition为bool类型，为false时表示申请被拒。

  - **接口地址：/consent/:id**

    **功能：表单创建者通过申请**

    **方法类型：POST**

    请求参数：指定申请的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回通过成功消息。

  - **接口地址：/refuse/:id**

    **功能：表单创建者拒绝申请**

    **方法类型：PUT**

    请求参数：指定申请的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回拒绝成功消息。

  - **接口地址：/quit/:set/:group**

    **功能：用户组退出某个表单**

    **方法类型：DELETE**

    请求参数：指定用户组的uuid（在接口地址的group处）、指定表单的uuid（在接口地址的set处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回删除成功消息。

- **BlockInterface**

  - **接口地址：/block/:set/:group**

    **功能：表单拉黑某用户组**

    **方法类型：POST**

    请求参数：指定用户组的uuid（在接口地址的group处）、指定表单的uuid（在接口地址的set处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回拉黑成功消息。

  - **接口地址：/remove/black/:set/:group**

    **功能：移除某用户组的黑名单**

    **方法类型：DELETE**

    请求参数：指定用户组的uuid（在接口地址的group处）、指定表单的uuid（在接口地址的set处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，返回移除黑名单成功消息。

  - **接口地址：/black/list/:id**

    **功能：查看黑名单**

    **方法类型：GET**

    请求参数：指定表单的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时，以json格式返回一个setBlocks，其为setBlock数组，每个元素包含了id、create_at、updated_at、set_id、group_id。其中group_id为被拉黑用户组的id，set_id为黑名单持有表单的uuid。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞或点踩**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出表单的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞或点踩**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出表单的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出表单的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出表单的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回setLikes和total，total表示点赞或者点踩的数量，seteLikes为setLike数组，seetLike包含了user_id表示点赞用户的uid，set_id表示点赞的表单的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回setLikes和total，total表示点赞或者点踩的数量，seteLikes为setLike数组，seetLike包含了user_id表示点赞用户的uid，set_id表示点赞的表单的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定表单的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **CollectInterface**

  - **接口地址：/collect/:id**

    **功能：收藏**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定表单的id（即:id部分） 。

    返回值：返回收藏成功信息。

  - **接口地址：/cancel/collect/:id**

    **功能：取消收藏**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定表单的id（即:id部分） 。

    返回值：返回取消收藏成功信息。

  - **接口地址：/collect/show/:id**

    **功能：查看收藏状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定表单的id（即:id部分） 。

    返回值：返回collect，为bool类型，为true表示已经收藏，false表示未收藏。

  - **接口地址：/collect/list/:id**

    **功能：查看收藏列表**

    **方法类型：GET**

    请求参数：在接口地址中给出指定表单的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回setCollects和total，其为setCollect数组，包含了user_id表示收藏用户的uid，set_id表示收藏的表单的uid，create_at表示收藏时间。total表示收藏总数。

  - **接口地址：/collect/number/:id**

    **功能：查看表单被收藏数量**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定表单的id（即:id部分） 。

    返回值：返回total表示收藏人次。

  - **接口地址：/collects/:id**

    **功能：查看用户收藏夹**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回groupCollects和total，其为groupCollect数组，包含了user_id表示收藏用户的uid，set_id表示收藏的表单的uid，create_at表示收藏时间。total表示收藏总数。

- **LabelInterface**

  - **接口地址：/label/:id/:label**

    **功能：创建表单标签**

    **方法类型：POST**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定表单的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回创建成功消息

  - **接口地址：/label/:id/:label**

    **功能：删除表单标签**

    **方法类型：DELETE**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定表单的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回删除成功消息

  - **接口地址：/label/:id**

    **功能：查看表单标签**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定表单的id（即:id部分）  。

    返回值：返回setLabels,其为setLabel数组，每个元素包含了一个 label字符串表示标签，created_at表示创建时间，set_id表示表单的uid。

- **SearchInterface**

  - **接口地址：/search/:text**

    **功能：按文本搜索表单**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个表单，默认值为20）。

    返回值：返回sets和total，total表示搜索到的表单总量。sets为set的数组

  - **接口地址：/search/label**

    **功能：按标签搜索表单**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个表单，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回sets和total，total表示搜索到的表单总量。sets为set的数组

  - **接口地址：/search/with/label/:text**

    **功能：按文本和标签交集搜索表单**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少表单，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回sets和total，total表示搜索到的表单总量。sets为set的数组

- **HotInterface**

  - **接口地址：/hot/rank**

    **功能：获取表单热度排行**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个表单，默认值为20）。

    返回值：返回sets和total，total表示表单总量。sets的每个元素包含member和score，其中member为set的uid，score为set对应的热度。已按热度排序。

- **其它**

  - **接口地址：/rank/list/:id**

    **功能：查看表单内用户排行**

    **方法类型：GET**

    请求参数： 在接口地址中给出表单的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。

    返回值：成功时，以json格式返回一个数组groups和total，groups返回了相应列表的用户组信息（按照创建时间排序，越新创建排序越前），total表示用户组总量，如果失败则返回失败原因。

  - **接口地址：/rank/update/:id**

    **功能：更新表单排行**

    **方法类型：PUT**

    请求参数： 在接口地址中给出用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个表单，默认值为20）。

    返回值：成功时，以json格式返回一个数组groups和total，groups返回了相应列表的用户组信息（按照创建时间排序，越新创建排序越前），total表示用户组总量，如果失败则返回失败原因。

  - **接口地址：/user/list/:id**

    **功能：查看某一用户的表单列表**

    **方法类型：GET**

    请求参数：在接口地址中给出表单的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个表单，默认值为20）。

    返回值：成功时，以json格式返回一个数组sets和total，sets返回了相应列表的表单信息（按照创建时间排序，越新创建排序越前），total表示表单总量，如果失败则返回失败原因。
    
  - **接口地址：/topic/list/:id**
  
    **功能：查看某一表单的主题列表**
  
    **方法类型：GET**
  
    请求参数：在接口地址中给出表单的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个主题，默认值为20）。
  
    返回值：成功时，以json格式返回一个数组topicLists和total，topicLists为topicList数组，每个topicList含有一个set_id表示表单id，一个topic_id表示表单id，返回了相应列表的主题信息，total表示主题总量，如果失败则返回失败原因。
  
  - **接口地址：/group/list/:id**
  
    **功能：查看某一表单的用户组列表**
  
    **方法类型：GET**
  
    请求参数：在接口地址中给出表单的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。
  
    返回值：成功时，以json格式返回一个数组groupLists和total，groupLists为groupList数组，每个groupList含有一个set_id表示表单id，一个group_id表示用户组id，返回了相应列表的用户组信息，total表示用户组总量，如果失败则返回失败原因。
  
  - **接口地址：/user/group/:user/:set**
  
    **功能：查看某用户在某表单在哪一组中**
  
    **方法类型：GET**
  
    请求参数：在接口地址中给出用户的id（即:user部分） 。在接口地址中给出表单的id（即:set部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少个用户组，默认值为20）。
  
    返回值：成功时，以json格式返回一个数组groups,groups为id数组，表示用户加入的group的id。

### 模型：SpecialJudge

定义：特殊判断

**基础路由：/specialjudge**

实现的接口类型：

- 其它

  - **接口地址：/create**（需要二级权限）

    **功能：创建特判程序**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含 language 、code、input。其要求在以code和language创建程序后，输入input内容可以输出"ok"字符。

    返回值：成功时返回创建成功相关信息并返回specialJudge，包含id、created_at、updated_at、user_id、language 、code、input。否则给出失败原因

  - **接口地址：/update/:id**（需要二级权限）

    **功能：更新特判程序**

    **方法类型：PUT**

    请求参数：在接口地址中给出特判的id（即id部分） 。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含 language 、code、input。其要求在以code和language创建程序后，输入input内容可以输出"ok"字符。

    返回值：成功时返回更新成功相关信息，否则给出失败原因

### 模型：Test

定义：本地测试

**基础路由：/test**

实现的接口类型：

- 其它

  - **接口地址：/create**

    **功能：创建测试**

    **方法类型：POST**

    请求参数：在Body，raw格式给出json类型数据包含 language 、code、input、  memory_limit 表示内存限制（kb）uint类型、 time_limit表示时间限制 （ms）uint类型。这里的language与record提交支持类型相同。

    返回值：成功时返回创建成功相关信息并返回output、condition、memory、time，分别表示输出、状态、运行消耗空间kb、运行消耗时间ms，其中condition与record提交返回格式基本相同，但没有"Accept"以及"Wrong Answer"，取而代之的是"ok"。否则给出失败原因

### 模型：Thread

定义：题解的回复

**基础路由：/thread**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create/:id**

    **功能：创建回复**

    **方法类型：POST**

    请求参数：题解的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示回复内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回创建成功相关信息，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看回复**

    **方法类型：GET**

    请求参数：回复的uuid（在接口地址的id处）。

    返回值：成功找到回复时，将会以json格式给出回复thread，thread中包含id,user_id,post_id,content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新回复**

    **方法类型：PUT**

    请求参数：回复的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示回复内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回更新成功相关信息，否则给出失败原因

  - **接口地址：/delete/:id**

    **功能：删除回复**

    **方法类型：DELETE**

    请求参数：回复的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功时返回删除成功相关信息，否则给出失败原因

  - **接口地址：/list/:id**

    **功能：查看回复列表**

    **方法类型：GET**

    请求参数：题解的uuid（在接口地址的id处）。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇回复，默认值为20）。

    返回值：成功时，以json格式返回一个数组threads和total，threads返回了相应列表的回复信息（按照创建时间排序，越新创建排序越前），total表示回复总量，如果失败则返回失败原因。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞或点踩回复**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞或点踩回复**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出回复的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出回复的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回threadLikes和total，total表示点赞或者点踩的数量，threadLikes为threadLike数组，threadLike包含了user_id表示点赞用户的uid，thread_id表示点赞的回复的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回threadLikes和total，total表示点赞或者点踩的数量，threadLikes为threadLike数组，threadLike包含了user_id表示点赞用户的uid，thread_id表示点赞的回复的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定回复的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **HotInterface**

  - **接口地址：/hot/rank/:id**

    **功能：获取回复热度排行**

    **方法类型：GET**

    请求参数：在接口地址中给出题解的id（即:id部分） 。 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇回复，默认值为20）。

    返回值：返回threads和total，total表示回复总量。threads的每个元素包含member和score，其中member为thread的uid，score为thread对应的热度。已按热度排序。

- **其它**

  - **接口地址：/user/list/:id**

    **功能：查看指定用户创建的回复列表**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：成功时，以json格式返回一个数组threads和total，threads返回了相应列表的回复信息（按照创建时间排序，越新创建排序越前），total表示回复总量，如果失败则返回失败原因。

### 模型：Topic

定义：主题

**基础路由：/topic**

实现的接口类型：

- **RestInterface**

  - **接口地址：/create**

    **功能：创建主题**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、 problems ，其中title表示主题标题，content表示主题内容，res_long表示长文本备用键值，res_short表示短文本备用键值，problems 为string数组，表示题目们的id。

    返回值：成功时返回创建成功相关信息，否则给出失败原因

  - **接口地址：/show/:id**

    **功能：查看主题**

    **方法类型：GET**

    请求参数：主题的uuid（在接口地址的id处）。

    返回值：成功找到主题时，将会以json格式给出主题topic，topic中包含id,user_id,title,content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/update/:id**

    **功能：更新主题**

    **方法类型：PUT**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选）、 problems ，其中title表示主题标题，content表示主题内容，res_long表示长文本备用键值，res_short表示短文本备用键值，problems 为string数组，表示题目们的id。

    返回值：成功更新主题时，返回更新成功信息。如果失败则返回失败原因。

  - **接口地址：/delete/:id**

    **功能：删除主题**

    **方法类型：DELETE**

    请求参数：主题的uuid（在接口地址的id处）。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功删除主题时，返回删除成功信息。如果失败则返回失败原因。

  - **接口地址：/list**

    **功能：查看主题列表**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：成功时，以json格式返回一个数组topics和total，topics返回了相应列表的主题信息（按照创建时间排序，越新创建排序越前），total表示主题总量，如果失败则返回失败原因。

- **LikeInterface**

  - **接口地址：/like/:id**

    **功能：点赞、点踩主题**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出主题的id（即:id部分） 。

    返回值：返回点赞成功消息

  - **接口地址：/cancle/like/:id**

    **功能：取消点赞、点踩状态**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在接口地址中给出主题的id（即:id部分） 。

    返回值：返回取消点赞成功消息

  - **接口地址：/like/number/:id**

    **功能：查看点赞点踩数量**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出主题的id（即:id部分） 。

    返回值：返回total表示点赞或者点踩的数量

  - **接口地址：/like/list/:id**

    **功能：查看点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出文章的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回topicLikes和total，total表示点赞或者点踩的数量，topicLikes为topicLike数组，topicLike包含了user_id表示点赞用户的uid，topic_id表示点赞的主题的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/likes/:id**

    **功能：查看用户点赞、点踩列表**

    **方法类型：GET**

    请求参数：在Params处提供bool类型的like，true为点赞，false为点踩。 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少点赞信息，默认值为20）。

    返回值：返回topicLikes和total，total表示点赞或者点踩的数量，topicLikes为topicLike数组，topicLike包含了user_id表示点赞用户的uid，topic_id表示点赞的主题的uid，create_at表示点赞时间，like表示为点赞true或者点踩false。

  - **接口地址：/like/show/:id**

    **功能：查看用户当前点赞状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定主题的id（即:id部分） 。

    返回值：返回like，like为int类型，0表示无状态，1表示已经点赞，-1表示已经点踩。

- **CollectInterface**

  - **接口地址：/collect/:id**

    **功能：收藏**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定主题的id（即:id部分） 。

    返回值：返回收藏成功信息。

  - **接口地址：/cancel/collect/:id**

    **功能：取消收藏**

    **方法类型：DELETE**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定主题的id（即:id部分） 。

    返回值：返回取消收藏成功信息。

  - **接口地址：/collect/show/:id**

    **功能：查看收藏状态**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定主题的id（即:id部分） 。

    返回值：返回collect，为bool类型，为true表示已经收藏，false表示未收藏。

  - **接口地址：/collect/list/:id**

    **功能：查看收藏列表**

    **方法类型：GET**

    请求参数：在接口地址中给出指定主题的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回topicCollects和total，其为topicCollect数组，包含了user_id表示收藏用户的uid，topic_id表示收藏的主题的uid，create_at表示收藏时间。total表示收藏总数。

  - **接口地址：/collect/number/:id**

    **功能：查看文章被收藏数量**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定文章的id（即:id部分） 。

    返回值：返回total表示收藏人次。

  - **接口地址：/collects/:id**

    **功能：查看用户收藏夹**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少收藏信息，默认值为20）。

    返回值：返回topicCollects和total，其为topicCollect数组，包含了user_id表示收藏用户的uid，topic_id表示收藏的主题的uid，create_at表示收藏时间。total表示收藏总数。

- **VisitInterface**

  - **接口地址：/visit/:id**

    **功能：游览主题**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定主题的id（即:id部分） 。

    返回值：返回游览成功消息。

  - **接口地址：/visit/number/:id**

    **功能：游览人次**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 在接口地址中给出指定主题的id（即:id部分） 。

    返回值：返回total表示游览人次。

  - **接口地址：/visit/list/:id**

    **功能：游览人次列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定主题的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：返回topicVisits和total，total表示游览总量。topicVisits为topicVisit数组，其包含了。包含了user_id表示游览用户的uid，topic_id表示游览的主题的uid，create_at表示游览时间。

  - **接口地址：/visits/:id**

    **功能：游览历史**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：返回topicVisits和total，total表示游览总量。topicVisits为topicVisit数组，其包含了。包含了user_id表示游览用户的uid，topic_id表示游览的主题的uid，create_at表示游览时间。

- **SearchInterface**

  - **接口地址：/search/:text**

    **功能：按文本搜索主题**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇主题，默认值为20）。

    返回值：返回topics和total，total表示搜索到的主题总量。topics为topic的数组

  - **接口地址：/search/label**

    **功能：按标签搜索主题**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇主题，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回topics和total，total表示搜索到的主题总量。topics为topic的数组

  - **接口地址：/search/with/label/:text**

    **功能：按文本和标签交集搜索主题**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇主题，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回topics和total，total表示搜索到的主题总量。topics为topic的数组

- **HotInterface**

  - **接口地址：/hot/rank**

    **功能：获取文章热度排行**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：返回articles和total，total表示文章总量。articles的每个元素包含member和score，其中member为article的uid，score为article对应的热度。已按热度排序。

- **LabelInterface**

  - **接口地址：/label/:id/:label**

    **功能：创建主题标签**

    **方法类型：POST**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定主题的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回创建成功消息

  - **接口地址：/label/:id/:label**

    **功能：删除主题标签**

    **方法类型：DELETE**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。  在接口地址中给出指定主题的id（即:id部分） 。 在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回删除成功消息

  - **接口地址：/label/:id**

    **功能：查看主题标签**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定文章的id（即:id部分）  。

    返回值：返回topicLabels,其为topicLabel数组，每个元素包含了一个 label字符串表示标签，created_at表示创建时间，topic_id表示主题uid。

- **其它**

  - **接口地址：/user/list/:id**

    **功能：查看指定用户的主题列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定用户的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：成功时，以json格式返回一个数组topics和total，topics返回了相应列表的主题信息（按照创建时间排序，越新创建排序越前），total表示主题总量，如果失败则返回失败原因。

  - **接口地址：/problem/list/:id**

    **功能：查看某一主题的题目列表**

    **方法类型：GET**

    请求参数： 在接口地址中给出指定主题的id（即:id部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少游览信息，默认值为20）。

    返回值：成功时，以json格式返回problemLists和total，problemLists为problemList数组，含有topic_id表示主题的id和problem_id表示题目的id，total表示总量，如果失败则返回失败原因。

### 模型：User

定义：用户

**基础路由：/user**

实现的接口类型：

- **SearchInterface**

  - **接口地址：/search/:text**

    **功能：按文本搜索用户**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户，默认值为20）。

    返回值：返回users和total，total表示搜索到的用户总量。users为user的数组，每个user含有name（名称）、email（邮箱地址）、blog（博客地址）、sex（bool类型，性别）、icon（头像的url）、level（权限等级）、score（竞赛分数）、like_num（不算今日的点赞数量）、unlike_num（不算今日的点踩数量）、collect_num（不算今日的收藏数量）、visit_num（不算今日的游览数量）

  - **接口地址：/search/label**

    **功能：按标签搜索用户**

    **方法类型：GET**

    请求参数： 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回users和total，total表示搜索到的用户总量。users为user的数组，每个user含有name（名称）、email（邮箱地址）、blog（博客地址）、sex（bool类型，性别）、icon（头像的url）、level（权限等级）、score（竞赛分数）、like_num（不算今日的点赞数量）、unlike_num（不算今日的点踩数量）、collect_num（不算今日的收藏数量）、visit_num（不算今日的游览数量）

  - **接口地址：/search/with/label/:text**

    **功能：按文本和标签交集搜索文章**

    **方法类型：GET**

    请求参数： 在接口地址中给出需要搜索的字符串（即:text部分） 。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户，默认值为20）。在Body，raw格式给出json类型数据包含一个labels，其为字符串数组表示要搜索的标签。

    返回值：返回users和total，total表示搜索到的用户总量。users为user的数组，每个user含有name（名称）、email（邮箱地址）、blog（博客地址）、sex（bool类型，性别）、icon（头像的url）、level（权限等级）、score（竞赛分数）、like_num（不算今日的点赞数量）、unlike_num（不算今日的点踩数量）、collect_num（不算今日的收藏数量）、visit_num（不算今日的游览数量）

- **LabelInterface**

  - **接口地址：/label/:label**

    **功能：创建用户标签**

    **方法类型：POST**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。   在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回创建成功消息

  - **接口地址：/label/:label**

    **功能：删除用户标签**

    **方法类型：DELETE**

    请求参数： Authorization中的Bearer Token中提供注册、登录时给出的token。   在接口地址中给出指定标签内容（即:label部分） 。

    返回值：返回删除成功消息

  - **接口地址：/label/:id**

    **功能：查看用户标签**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：返回userLabels,其为userLabel数组，每个元素包含了一个 label字符串表示标签，created_at表示创建时间，user_id表示用户uid。

- **其它**

  - **接口地址：/verify/:id**

    **功能：验证码获取**

    **方法类型：GET**

    请求参数：需要在接口地址部分（:id）给出用户邮箱

    返回值：成功则返回200并向相应邮箱发送验证邮件，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/regist**

    **功能：用户注册**

    **方法类型：POST**

    请求参数：Body部分，form-data类型，接收四个字符串分别为Email，Password，Name，Verify。其中Email需要符合邮箱格式，Password最少需要六位，Name最多为20位长度，Verify必须与邮箱验证码相同，注意，用户的邮箱和名称均不能重复。

    返回值：成功则返回200与token，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/login**

    **功能：用户登录**

    **方法类型：POST**

    请求参数：Body部分，form-data类型，接收两个字符串分别为Email，Password。其中Email需要符合邮箱格式，Password最少需要六位。

    返回值：成功则返回200与token，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/security**

    **功能：找回密码**

    **方法类型：PUT**

    请求参数：Body部分，form-data类型，接收两个字符串分别为Email，Verify。其中Verify必须与邮箱验证码相同。

    返回值：成功则返回200并向相应邮箱发送重置后的密码，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/update/password**

    **功能：更新密码**

    **方法类型：PUT**

    请求参数：Body部分，form-data类型，接收两个字符串分别为first，second。其中first为旧密码，second为新密码。

    返回值：成功则返回200并修改数据库中的密码，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/info**

    **功能：返回当前登录的用户**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。 

    返回值：返回用户的一些个人信息，格式为json包含name（用户名称）、email（邮箱）、blog（博客地址）、sex（bool类型，性别）、address（地址）、icon（头像）、level(权限等级)、score(竞赛分数)、like_num(收到点赞数量)、unlike_num（收到点踩数量）、collect_num（收到收藏数量）、visit_num（被游览数量）

  - **接口地址：/show/:id**

    **功能：获取某个用户的所有信息**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：返回用户的一些个人信息，格式为json包含name（用户名称）、email（邮箱）、blog（博客地址）、sex（bool类型，性别）、address（地址）、icon（头像）、level(权限等级)、score(竞赛分数)、like_num(收到点赞数量)、unlike_num（收到点踩数量）、collect_num（收到收藏数量）、visit_num（被游览数量）

  - **接口地址：/update**

    **功能：修改用户信息**

    **方法类型：PUT**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body 中，raw格式提供json包含name（用户名称）、email（邮箱）、blog（博客地址）、sex（bool类型，性别）、address（地址）、icon（头像）

    返回值：更新成功后返回用户更新后的个人信息，否则返回错误原因。

  - **接口地址：/update/level/:id/:level**

    **功能：修改用户等级**

    **方法类型：PUT**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）， 在接口地址中给出需要修改的level（即:level部分） 。

    返回值：更新成功后返回更新成功信息，否则返回错误原因。

  - **接口地址：/accept/rank/show/:id**

    **功能：获取某个用户的ac题目数量排行情况**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：返回一个rank，为int类型，表示指定用户当前的排行。

  - **接口地址：/accept/rank/list**

    **功能：获取ac题目数量排行榜**

    **方法类型：GET**

    请求参数：  在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户，默认值为20）。

    返回值：返回acceptRanks和total，total表示排行总量。acceptRanks的每个元素包含accept_num和user_id，其中user_id为user的uid，accept_num为user对应的ac数量。已按ac数量排序。

  - **接口地址：/accept/num/:id**

    **功能：获取某个用户的ac题目数量**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：返回一个num，为int类型，表示指定用户的ac题目数量。

  - **接口地址：/score/change/:id**

    **功能：查看某一用户分数变化情况**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。 在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少信息，默认值为20）。

    返回值：返回userScoreChanges和total，total表示信息总量。userScoreChanges的每个元素包含id、 created_at 、 updated_at 、 score_change(浮点类型，表示分数变化) 、competition_id（竞赛id）和user_id，按照创建时间降序排序。

  - **接口地址：/hot/:id**

    **功能：查看某一用户今日热度数据**

    **方法类型：GET**

    请求参数：  在接口地址中给出指定用户的id（即:id部分）  。

    返回值：返回VisitNum、LikeNum、UnLikeNum、CollectNum分别表示指定用户的今日被游览人次、收到点赞数量、收到点踩数量、被收藏数量。

  - **接口地址：/like/rank**

    **功能：查看今日点赞榜单**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户，默认值为20）。

    返回值：返回users和total,total表示总量，users为数组，每个元素还有member和score，member为用户的uid，score为点赞数量。

  - **接口地址：/unlike/rank**

    **功能：查看今日点踩榜单**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户，默认值为20）。

    返回值：返回users和total,total表示总量，users为数组，每个元素还有member和score，member为用户的uid，score为点踩数量。

  - **接口地址：/collect/rank**

    **功能：查看今日收藏榜单**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户，默认值为20）。

    返回值：返回users和total,total表示总量，users为数组，每个元素还有member和score，member为用户的uid，score为收藏数量。

  - **接口地址：/visit/rank**

    **功能：查看今日游览榜单**

    **方法类型：GET**

    请求参数：在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少用户，默认值为20）。。

    返回值：返回users和total,total表示总量，users为数组，每个元素还有member和score，member为用户的uid，score为游览数量。

