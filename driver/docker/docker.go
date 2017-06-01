package docker

import (
	"encoding/json"
	"io/ioutil"
	"net"
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
