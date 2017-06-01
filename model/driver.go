package model

// Driver 部署服务接口
type Driver interface {
	// List 获取当前Driver所支持的所有服务
	List() ([]Service, error)
	// Start 启动指定服务
	Start(src Service) error
	// Restart 重启指定服务
	Restart(src Service) error
	// PullImg 拉取指定镜像
	PullImg(img string) error
}

// Service 服务信息
type Service struct {
	// Name 服务名称
	Name string `json:"name"`
	// Type 服务类型
	Type string `json:"type"`
	// Status 服务当前状态
	// 0: Running
	// 1: Stopped
	// 2: Hangup
	Status string `json:"status"`
	// Description 服务描述信息
	Description string `json:"description"`
}

// MetaData 部署原信息
type MetaData struct {
	// Auth 认证信息
	Auth string `json:"auth"`
}
