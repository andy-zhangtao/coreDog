package util

const (
	// RESTART_SRV = "/post/restart/{name}"
	// LISTALLSRV 获取所有service API
	LISTALLSRV = "/get/all/services"
	// PULLIMG 拉取镜像API
	PULLIMG = "/put/docker/img"
	// STARTSRV 启动service API
	STARTSRV = "/post/start/service"
	// RESTAETSRV 重启service API
	RESTAETSRV = "/put/restart/service"
	// SYSTEMD systemd驱动名称
	SYSTEMD = "SYSTEMD"
	// RANCHER rancher 驱动名称
	RANCHER = "RANCHER"
)

var (
	// VERSION docker版本列表
	VERSION = map[string]string{
		"17.04":   "1.28",
		"17.03.1": "1.27",
		"1.13.1":  "1.26",
		"17.03.0": "1.26",
		"1.13.0":  "1.25",
		"1.12":    "1.24",
		"1.11":    "1.23",
		"1.10":    "1.22",
		"1.9":     "1.21",
		"1.8":     "1.20",
		"1.7":     "1.19",
		"1.6":     "1.18",
	}
)
