package driver

import (
	"errors"
	"fmt"

	"github.com/andy-zhangtao/coreDog/model"
	"github.com/coreos/go-systemd/dbus"
)

// Systemd Coreos系统
type Systemd struct{}

// List 获取coreos系统中所有服务
func (s Systemd) List() ([]model.Service, error) {
	conn, err := dbus.New()
	if err != nil {
		msg := fmt.Sprintf("DBUS链接失败[%s]", err.Error())
		return nil, errors.New(msg)
	}

	units, err := conn.ListUnits()
	if err != nil {
		msg := fmt.Sprintf("查询服务文件失败[%s]", err.Error())
		return nil, errors.New(msg)
	}

	var src []model.Service

	for _, u := range units {
		s := model.Service{
			Name:        u.Name,
			Type:        u.JobType,
			Status:      u.ActiveState,
			Description: u.Description,
		}

		src = append(src, s)
	}

	return src, nil
}

func (s Systemd) Start(src model.Service) error {
	return nil
}

func (s Systemd) Restart(src model.Service) error {
	return nil
}

func (s Systemd) PullImg(img string) error {
	return nil
}
