# coreDog
A smart CD tools

## API List

* GET /v1/get/all/services 获取指定驱动所支持的所有服务

    当调用systemd驱动时:
    |query param|value|desc|
    |-----------|-----|----|
    |driver|systemd|调用coreos的systemd服务|
    当调用rancher驱动时:
    |query param|value|desc|
    |-----------|-----|----|
    |driver|rancher|调用第三方Rancher服务|
    |accesskey|<string>|rancher生成的访问key|
    |secretkey|<string>|rancher生成的访问密钥|
    |env|<string>|rancher环境名称,大小写敏感|
    |domain|<string>|rancher访问地址,若为空，则为localhost:8080|


* PUT /v1/put/docker/img 通知指定驱动下载指定镜像
    |query param|value|desc|
    |-----------|-----|----|
    |driver|systemd|调用coreos的systemd服务|
    |driver|rancher|调用第三方Rancher服务|
    
    query param: 
    driver:
        * systemd
    img: 镜像名称

* POST /v1/post/start/service 通知指定驱动启动指定服务
    query param: 
    driver:
        * systemd
    srv: 服务名称

* PUT /v1/put/restart/service 通知指定驱动重启指定服务
    query param: 
    driver:
        * systemd
    srv: 服务名称