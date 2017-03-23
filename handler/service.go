package handler

import (
	"fmt"
	"net/http"

	"github.com/andy-zhangtao/Sandstorm"
	"github.com/coreos/go-systemd/dbus"
	"github.com/gorilla/mux"
)

var msg string

// RestartService 重启服务
// @name 服务名称
func RestartService(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		msg = fmt.Sprintf("重启服务名不得为空")
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	target := name + ".service"
	conn, err := dbus.New()
	if err != nil {
		msg = fmt.Sprintf("DBUS链接失败[%s]", err.Error())
	}

	defer func() {
		fmt.Println(msg)
		if conn != nil {
			conn.Close()
		}
	}()

	units, err := conn.ListUnits()
	if err != nil {
		msg = fmt.Sprintf("查询服务文件失败[%s]", err.Error())
	}

	unit := getUnitStatus(units, target)
	if unit == nil {
		msg = fmt.Sprintf("没有查询到服务文件[%s]", target)
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	} else if unit.ActiveState != "active" {
		msg = fmt.Sprintf("服务文件[%s]未处于active状态,当前状态[%s]", target, unit.ActiveState)
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	// Restart the unit
	reschan := make(chan string)
	_, err = conn.RestartUnit(target, "replace", reschan)
	if err != nil {
		msg = fmt.Sprintf("服务[%s]重启失败[%s]", target, err.Error())
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	job := <-reschan
	if job != "done" {
		msg = fmt.Sprintf("服务重启失败,重启返回信息[%s]", job)
		Sandstorm.HTTPError(w, msg, http.StatusInternalServerError)
		return
	}

	
}

func getUnitStatus(units []dbus.UnitStatus, name string) *dbus.UnitStatus {
	for _, u := range units {
		if u.Name == name {
			return &u
		}
	}
	return nil
}
