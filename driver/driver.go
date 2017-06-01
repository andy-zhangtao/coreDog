package driver

import "github.com/andy-zhangtao/coreDog/model"

// ListService 获取所有服务
func ListService(d model.Driver) ([]model.Service, error) {
	return d.List()
}

// StartService 启动指定服务
// srv 服务名称
func StartService(d model.Driver, srv string) error {
	s := model.Service{
		Name: srv,
	}
	return d.Start(s)
}

// RestartService 重启指定服务
// srv 服务名称
func RestartService(d model.Driver, srv string) error {
	s := model.Service{
		Name: srv,
	}
	return d.Restart(s)
}

// PullImg 下载指定镜像
// img 镜像名称
func PullImg(d model.Driver, img string) error {
	return d.PullImg(img)
}
