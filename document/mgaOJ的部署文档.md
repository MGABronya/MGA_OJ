# mgaOJ的部署文档

一共有四个服务需要部署如下

img，简易图床

test，本地测试

judge，判题机

mgaoj，主要服务

这里将简易的说明这四个服务的部署流程

## Img

#### 服务：简易图床

#### 部署文件包：https://github.com/MGABronya/MGA_OJ/blob/main/Img/Img.zip

你可以使用以下命令下载它

````
wget https://github.com/MGABronya/MGA_OJ/blob/main/Img/Img.zip
````

你可以将该包解压在linux机器上。

````
unzip -o Img.zip -d /home/ubuntu
````

这里的/home/ubuntu可以替换为你想要部署的目录。

您可以修改该包config目录下的application.yml文件来修改监听的端口。原application.yml内容如下

````
server:
  port: 1000
````

你可以修改1000的值以此来修改监听端口，如需访问，请关闭该端口的防火墙。

那之后，在部署目录下使用指令

````
nohup ./Img > Img.log 2>&1 &
````

将Img服务挂在后台，你可以查看文件Img.log了解日志信息。

如果没有运行权限，可以运行以下指令获取权限

````
chmod 777 Img
````

如果需要停止它，可以使用

````
ps -aux | grep ./Img
````

找到进程id，然后

````
kill -9 id
````

杀死它。

如果该服务是为了支持MGA_OJ而部署的，请保证同目录下存在以下文件。

````
MGA1.jpg	MGA2.jpg	MGA3.jpg	MGA4.jpg	MGA5.jpg	MGA6.jpg	MGA7.jpg	MGA8.jpg	MGA9.jpg
````

如果不存在，你可以在https://github.com/MGABronya/MGA_OJ/tree/main/Icon找到他们

当然，只要确保文件名相同，你也可以将这些图片替换为其它图片。

在这之后，你可能需要用到nginx来提供反向代理服务。在你下载并启动nginx服务后，可以通过修改nginx的配置文件达到反向代理目的。这里给出一个修改示例。

你可以使用

````
whereis nginx
````

寻找nginx的目录。在该示例中，其目录为/etc/nginx

再这之后在nginx.conf文件中的http中添加以下内容

````
http{
····
	server {
                        listen       84;
                        server_name  localhost;

                        location / {
                        root   /home/ubuntu;
                        }
        }
···
}
````

这里的84可以修改为任意一个未被占用的端口，如需访问，这里需要该端口的防火墙关闭。

/home/ubuntu需要修改为你部署Img服务的目录。

修改完成后使用指令

````
nginx -s reload
````

重新读入配置文件。

在这之后，你可以通过http://localhost:84/后加文件名的方式访问到上传的图片。

至于上传文件的方式，可以参考接口文档中Img部分https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3.md

#### 实例

上传文件：

api_img.mgaronya.com/img/upload

查看文件MGA1.jpg：

img.mgaronya.com/MGA1.jpg

## Test

#### 服务：本地测试

#### 部署文件包：https://github.com/MGABronya/MGA_OJ/blob/main/Test/TestMachine.zip

你可以使用以下命令下载它

````
wget https://github.com/MGABronya/MGA_OJ/blob/main/Test/TestMachine.zip
````

你可以将该包解压在linux机器上。

````
unzip -o TestMachine.zip -d /home/ubuntu/Test0
````

这里的/home/ubuntu/Test0可以替换为你想要部署的目录。

您可以修改该包config目录下的application.yml文件来修改监听的端口。原application.yml内容如下

````
server:
  port: 1000
````

你可以修改1000的值以此来修改监听端口，如需访问，请关闭该端口的防火墙。

那之后，在部署目录下的Test目录下使用指令

````
nohup ./Test > Test.log 2>&1 &
````

将Test服务挂在后台，你可以查看文件Test.log了解日志信息。

如果没有运行权限，可以运行以下指令获取权限

如果没有运行权限，可以运行以下指令获取权限

````
chmod 777 Test
````

如果需要停止它，可以使用

````
ps -aux | grep ./Test
````

找到进程id，然后

````
kill -9 id
````

杀死它。