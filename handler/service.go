package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/andy-zhangtao/Sandstorm"
	"github.com/andy-zhangtao/coreDog/driver"
	"github.com/andy-zhangtao/coreDog/model"
	"github.com/andy-zhangtao/coreDog/util"
	"github.com/coreos/go-systemd/dbus"
)

var msg string

// RestartService 重启指定服务
// driver 驱动名称
// srv 服务名称
func RestartService(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	srv, ok := vals["srv"]
	if !ok {
		msg = fmt.Sprintf("服务名称不得为空")
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	oriDriver, ok := vals["driver"]
	if !ok {
		msg = fmt.Sprintf("驱动名称不得为空")
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	var err error

	switch strings.ToUpper(oriDriver[0]) {
	case util.SYSTEMD:
		dri := driver.Systemd{}
		err = driver.RestartService(dri, srv[0])
	}

	if err != nil {
		Sandstorm.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Sandstorm.HTTPSuccess(w, "")
	return
}

// PullImgService 拉取指定镜像
// img 镜像名称
// driver 驱动名称
func PullImgService(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	img, ok := vals["img"]
	if !ok {
		msg = fmt.Sprintf("镜像名称不得为空")
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	oriDriver, ok := vals["driver"]
	if !ok {
		msg = fmt.Sprintf("驱动名称不得为空")
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	var err error

	switch strings.ToUpper(oriDriver[0]) {
	case util.SYSTEMD:
		dri := driver.Systemd{}
		err = driver.PullImg(dri, img[0])
	}

	if err != nil {
		Sandstorm.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Sandstorm.HTTPSuccess(w, "")
	return
}

// ListService 获取指定驱动所支持的所有服务
// driver: 指定驱动类型。 systemd/
func ListService(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	oriDriver, ok := vals["driver"]
	if !ok {
		msg = fmt.Sprintf("驱动类型不得为空")
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	var srv []model.Service
	var err error

	switch strings.ToUpper(oriDriver[0]) {
	case util.SYSTEMD:
		dri := driver.Systemd{}
		srv, err = driver.ListService(dri)
	case util.RANCHER:
		accessKey, ok := vals["accesskey"]
		if !ok {
			msg = fmt.Sprintf("驱动为Rancher时,AccessKey 不得为空")
			Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
			return
		}
		secretKey, ok := vals["secretkey"]
		if !ok {
			msg = fmt.Sprintf("驱动为Rancher时,secretKey 不得为空")
			Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
			return
		}

		domain, ok := vals["domain"]
		if !ok {
			domain = append(domain, "localhost:8080")
		}
		env, ok := vals["env"]
		if !ok {
			msg = fmt.Sprintf("驱动为Rancher时,env 不得为空")
			Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
			return
		}

		dri := driver.Rancher{
			AccessKey: accessKey[0],
			SecretKey: secretKey[0],
			Domain:    domain[0],
			Env:       env[0],
		}

		srv, err = driver.ListService(dri)
	}

	if err != nil {
		Sandstorm.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := json.Marshal(srv)
	if err != nil {
		Sandstorm.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Sandstorm.HTTPSuccess(w, string(content))
	return
}

// StartService 启动指定服务
// driver 驱动类型
// service 服务名称
func StartService(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	srv, ok := vals["srv"]
	if !ok {
		msg = fmt.Sprintf("服务名称不得为空")
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	oriDriver, ok := vals["driver"]
	if !ok {
		msg = fmt.Sprintf("驱动名称不得为空")
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	var err error

	switch strings.ToUpper(oriDriver[0]) {
	case util.SYSTEMD:
		dri := driver.Systemd{}
		err = driver.StartService(dri, srv[0])
	}

	if err != nil {
		Sandstorm.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Sandstorm.HTTPSuccess(w, "")
	return

}

// RestartService 重启服务
// @name 服务名称
// func RestartService(w http.ResponseWriter, r *http.Request) {
// 	name := mux.Vars(r)["name"]
// 	if name == "" {
// 		msg = fmt.Sprintf("重启服务名不得为空")
// 		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
// 		return
// 	}

// 	target := name + ".service"
// 	conn, err := dbus.New()
// 	if err != nil {
// 		msg = fmt.Sprintf("DBUS链接失败[%s]", err.Error())
// 	}

// 	defer func() {
// 		fmt.Println(msg)
// 		if conn != nil {
// 			conn.Close()
// 		}
// 	}()

// 	units, err := conn.ListUnits()
// 	if err != nil {
// 		msg = fmt.Sprintf("查询服务文件失败[%s]", err.Error())
// 	}

// 	unit := getUnitStatus(units, target)
// 	if unit == nil {
// 		msg = fmt.Sprintf("没有查询到服务文件[%s]", target)
// 		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
// 		return
// 	} else if unit.ActiveState != "active" {
// 		msg = fmt.Sprintf("服务文件[%s]未处于active状态,当前状态[%s]", target, unit.ActiveState)
// 		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
// 		return
// 	}

// 	// Restart the unit
// 	reschan := make(chan string)
// 	_, err = conn.RestartUnit(target, "replace", reschan)
// 	if err != nil {
// 		msg = fmt.Sprintf("服务[%s]重启失败[%s]", target, err.Error())
// 		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
// 		return
// 	}

// 	job := <-reschan
// 	if job != "done" {
// 		msg = fmt.Sprintf("服务重启失败,重启返回信息[%s]", job)
// 		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
// 		return
// 	}

// }

func getUnitStatus(units []dbus.UnitStatus, name string) *dbus.UnitStatus {
	for _, u := range units {
		if u.Name == name {
			return &u
		}
	}
	return nil
}
