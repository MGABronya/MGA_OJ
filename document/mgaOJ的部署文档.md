# mgaOJ的部署文档

一共有四个服务需要部署如下

img，简易图床

test，本地测试

judge，判题机

mgaoj，主要服务

这里将简易的说明这四个服务的部署流程

## Img

#### 服务：简易图床

#### 部署文件包：https://github.com/MGABronya/MGA_OJ/blob/main/deploy/Img/Img.zip

你可以使用以下命令下载它

````
wget https://github.com/MGABronya/MGA_OJ/blob/main/deploy/Img/Img.zip
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
                        listen       80;
                        server_name  localhost;

                        location / {
                        root   /home/ubuntu;
                        }
        }
···
}
````

这里的80可以修改为任意一个未被占用的端口，如需访问，这里需要该端口的防火墙关闭。

/home/ubuntu需要修改为你部署Img服务的目录。

修改完成后使用指令

````
nginx -s reload
````

重新读入配置文件。

在这之后，你可以通过http://localhost:80/后加文件名的方式访问到上传的图片。

至于上传文件的方式，可以参考接口文档中Img部分https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3.md

#### 实例

上传文件：

api_img.mgaronya.com/img/upload

查看文件MGA1.jpg：

img.mgaronya.com/MGA1.jpg

## Test

#### 服务：本地测试

#### 部署文件包：https://github.com/MGABronya/MGA_OJ/blob/main/deploy/Test/Test.zip

你可以使用以下命令下载它

````
wget https://github.com/MGABronya/MGA_OJ/blob/main/deploy/Test/Test.zip
````

你可以将该包解压在linux机器上。

````
unzip -o Test.zip -d /home/ubuntu/Test0
````

这里的/home/ubuntu/Test0可以替换为你想要部署的目录。

运行该程序需要有C,C#,C++,C++11,Erlang,Go,Java,JavaScript,Kotlin,Pascal,PHP,Python,Racket,Ruby,Rust,Scala,Swift的语言运行环境。你可以通过以下命令获得实现了该环境的镜像。

````
docker pull mgaronya/oj-test
````

在这之后，将服务目录挂载到docker的home目录下，并绑定单核

````
docker run -itd --cpuset-cpus 0 -v /home/lighthouse/TestMachine/Test0:/home -p 1000:1000 --name ubuntu-test0 mgaronya/oj-test /bin/bash
````

你可以将0换为其它的核，/home/lighthouse/TestMachine/Test0挂载目录应该为你部署的目录，端口映射1000:1000应该为你部署服务时准备使用的端口，ubuntu-test0可以换为其它你想使用的名字。

那之后，通过以下指令找到运行时的容器

````
docker ps
````

然后进入它

````
docker exec -it id /bin/bash
````

这里的id为docker的id。那之后，你可以在挂载的home目录下找到部署包。

您可以修改该包config目录下的application.yml文件来修改监听的端口。原application.yml内容如下

````
server:
  port: 1000
````

你可以修改1000的值以此来修改监听端口，但是要与运行docker时映射的docker端口相同。

那之后，在部署目录下的Test目录下使用指令

````
nohup ./Test > Test.log 2>&1 &
````

将Test服务挂在后台，你可以查看文件Test.log了解日志信息。

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

如果你部署了多个Test服务，你或许需要使用到nginx来完成负载均衡。

在下载并启动nginx服务后，你可以使用

````
whereis nginx
````

寻找nginx的目录。在该示例中，其目录为/etc/nginx

再这之后在nginx.conf文件中的http中添加以下内容

````
http{
····
	upstream  test-server {
                        least_conn;
                        server    100.10.10.100:1000;
                        server    100.10.10.100:1001;
                        server    40.100.200.200:1001;
                        server    40.100.200.200:1002;
                        server    40.101.100.160:1003;
        }
        server {
                        listen       10000;
                        server_name test_oj.mgaronya.com;

                        location / {
                        proxy_pass http://test-server;
                        proxy_redirect default;
                        }
        }
···
}
````

你需要将test-server中的server们修改为你部署服务docker所在的服务器ip以及其映射的端口。

那之后，在server中你可以修改10000为任何一个你想要使用的端口，并使用你自己的域名替换test_oj.mgaronya.com

至于测试的方式，可以参考接口文档中Test部分https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3.md

实例

test_oj.mgaronya.com

## Judge

#### 服务：判题机

#### 部署文件包：https://github.com/MGABronya/MGA_OJ/blob/main/deploy/Judge/Judge.zip

你可以使用以下命令下载它

````
wget https://github.com/MGABronya/MGA_OJ/blob/main/deploy/Judge/Judge.zip
````

你可以将该包解压在linux机器上。

````
unzip -o Judge.zip -d /home/ubuntu/Judge0
````

这里的/home/ubuntu/Judge0可以替换为你想要部署的目录。

