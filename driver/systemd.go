package driver

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

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

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

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

// Start 定向启动指定服务
func (s Systemd) Start(srv model.Service) error {
	if srv.Name == "" {
		return errors.New("Service name is empty!")
	}

	conn, err := dbus.New()
	if err != nil {
		msg := fmt.Sprintf("DBUS链接失败[%s]", err.Error())
		return errors.New(msg)
	}

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	units, err := conn.ListUnits()
	if err != nil {
		msg := fmt.Sprintf("查询服务文件失败[%s]", err.Error())
		return errors.New(msg)
	}

	hasThisSrv := false
	var us dbus.UnitStatus

	for _, u := range units {
		if strings.Compare(u.Name, srv.Name) == 0 {
			hasThisSrv = true
			us = u
			break
		}
	}

	if !hasThisSrv {
		return errors.New("I can not find this service! Please confirem service name")
	}

	if us.ActiveState != "active" {
		msg := fmt.Sprintf("服务文件[%s]未处于active状态,当前状态[%s]", us.Name, us.ActiveState)
		return errors.New(msg)
	}

	reschan := make(chan string)
	_, err = conn.StartUnit(us.Name, "replace", reschan)
	if err != nil {
		msg := fmt.Sprintf("服务[%s]重启失败[%s]", us.Name, err.Error())
		return errors.New(msg)
	}

	job := <-reschan
	if job != "done" {
		msg := fmt.Sprintf("服务重启失败,重启返回信息[%s]", job)
		return errors.New(msg)
	}

	return nil
}

func (s Systemd) Restart(srv model.Service) error {
	if srv.Name == "" {
		return errors.New("Service name is empty!")
	}

	conn, err := dbus.New()
	if err != nil {
		msg := fmt.Sprintf("DBUS链接失败[%s]", err.Error())
		return errors.New(msg)
	}

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	units, err := conn.ListUnits()
	if err != nil {
		msg := fmt.Sprintf("查询服务文件失败[%s]", err.Error())
		return errors.New(msg)
	}

	hasThisSrv := false
	var us dbus.UnitStatus

	for _, u := range units {
		if strings.Compare(u.Name, srv.Name) == 0 {
			hasThisSrv = true
			us = u
			break
		}
	}

	if !hasThisSrv {
		return errors.New("I can not find this service! Please confirem service name")
	}

	if us.ActiveState != "active" {
		msg := fmt.Sprintf("服务文件[%s]未处于active状态,当前状态[%s]", us.Name, us.ActiveState)
		return errors.New(msg)
	}

	reschan := make(chan string)
	_, err = conn.RestartUnit(us.Name, "replace", reschan)
	if err != nil {
		msg := fmt.Sprintf("服务[%s]重启失败[%s]", us.Name, err.Error())
		return errors.New(msg)
	}

	job := <-reschan
	if job != "done" {
		msg := fmt.Sprintf("服务重启失败,重启返回信息[%s]", job)
		return errors.New(msg)
	}

	return nil
}

// PullImg 拉取指定镜像
func (s Systemd) PullImg(img string) error {
	sock, err := s.getDockerSocket()
	if err != nil {
		log.Println(err)
	}
	if sock == "" {
		return errors.New("Get [docker.sock] error!")
	}

	return nil
}

func (s Systemd) getDockerSocket() (string, error) {
	conn, err := dbus.New()
	if err != nil {
		msg := fmt.Sprintf("DBUS链接失败[%s]", err.Error())
		return "", errors.New(msg)
	}

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	// dockerInfo, err := conn.GetUnitProperties("docker.socket")
	// if err != nil {
	// 	return "", err
	// }

	// log.Println(dockerInfo)

	sc, err := conn.GetUnitTypeProperties("docker.socket", "Socket")
	if err != nil {
		return "", err
	}

	listen := sc["Listen"]

	st := reflect.ValueOf(listen)
	if st.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	ret := make([][]interface{}, st.Len())

	for i := 0; i < st.Len(); i++ {
		rvalue := interfaceSlice(st.Index(i).Interface())
		ret[i] = rvalue
	}

	var sock string
	for _, r := range ret {
		for _, rn := range r {
			s := reflect.ValueOf(rn)
			if s.Kind() == reflect.String {
				skey := s.Interface().(string)
				if strings.Contains(skey, "docker.sock") {
					sock = skey
				}
			}
		}

	}
	return sock, nil
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
