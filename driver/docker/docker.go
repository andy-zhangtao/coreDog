package docker

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

type DockerInfo struct {
	ServerVersion string `json:"ServerVersion"`
}

// GetDockerVersion 获取当前Docker Server版本
func GetDockerVersion(sock string) (string, error) {
	addr := net.UnixAddr{sock, "unix"}
	conn, err := net.DialUnix("unix", nil, &addr)
	if err != nil {
		return "", err
	}

	_, err = conn.Write([]byte("GET /info HTTP/1.0\r\n\r\n"))
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return "", err
	}

	body := getBody(result[:])

	var di DockerInfo

	err = json.Unmarshal(body, &di)
	if err != nil {
		return "", err
	}

	return di.ServerVersion, nil
}

func getBody(result []byte) (body []byte) {
	for i := 0; i <= len(result)-4; i++ {
		if result[i] == 13 && result[i+1] == 10 && result[i+2] == 13 && result[i+3] == 10 {
			body = result[i+4:]
			break
		}
	}
	return
}

// PullImg 拉取指定镜像
// sock docker.sock文件位置
// api docker api endpoint
func PullImg(sock, api string) error {
	addr := net.UnixAddr{sock, "unix"}
	conn, err := net.DialUnix("unix", nil, &addr)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte("POST " + api + " HTTP/1.0\r\n\r\n"))
	if err != nil {
		return err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return err
	}

	body := getBody(result[:])

	log.Println(string(body))
	if strings.Contains(string(body), "errorDetail") && strings.Contains(string(body), "error") {
		return errors.New(string(body))
	}

	return nil
}
