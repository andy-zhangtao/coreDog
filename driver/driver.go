package driver

import "github.com/andy-zhangtao/coreDog/model"

// ListService 获取所有服务
func ListService(d model.Driver) ([]model.Service, error) {
	return d.List()
}
