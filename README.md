# coreDog
A smart CD tools

## API List

* GET /v1/get/all/services 获取指定驱动所支持的所有服务
    query param: 
    driver:
        * systemd

* PUT /v1/put/docker/img 通知指定驱动下载指定镜像
    query param: 
    driver:
        * systemd
    img: 镜像名称

* POST /v1/post/start/service 通知指定驱动启动指定服务
    query param: 
    driver:
        * systemd
    srv: 服务名称
    