package driver

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/andy-zhangtao/coreDog/model"
)

const (
	// ENVAPI 获取Rancher env列表
	ENVAPI = "/v2-beta/projects/"
)

// Rancher rancher对象
type Rancher struct {
	// AccessKey Rancher用户生成的访问Key
	AccessKey string `json:"accesskey"`
	// SecretKey 与AccessKey相对应的密钥Key
	SecretKey string `json:"secretkey"`
	// Domain Rancher访问域名, 默认为localhost:8080
	Domain string `json:"domain"`
	// Env Rancher环境名称, 此处Env与Project为同一事物
	Env string `json:"env"`
	// Service Rancher服务名称, 此服务必须处于指定的env中
	Service string `json:"service"`
}

// RancherProject Rancher工程信息及Metadata信息
// 目前仅保存ProjectInfo信息
type RancherProject struct {
	Data []RancherProjectInfo `json:"data"`
}

// RancherProjectInfo 每个工程详细信息
// 目前仅保存ID和Name两项信息
type RancherProjectInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// List 获取Rancher中所有的服务
func (r Rancher) List() ([]model.Service, error) {
	if strings.TrimSpace(r.Env) == "" {
		return nil, errors.New("Rancher Env Can not be empty!")
	}

	envID, err := r.getEnvID(r.Env)
	if err != nil {
		return nil, err
	}

	log.Printf("Env ID [%s] \n", envID)

	return nil, nil
}

// Start 启动Rancher中指定服务
func (r Rancher) Start(srv model.Service) error {
	return nil
}

// Restart 重启Rancher中的指定服务
func (r Rancher) Restart(srv model.Service) error {
	return nil
}

// PullImg 拉取指定镜像
func (r Rancher) PullImg(img string) error {
	return nil
}

// getEnvID 获取指定env id
// env 指定的Env名称, 大小写敏感
func (r Rancher) getEnvID(env string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", r.Domain+ENVAPI, nil)
	req.SetBasicAuth(r.AccessKey, r.SecretKey)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var rancher RancherProject

	err = json.Unmarshal(content, &rancher)
	if err != nil {
		return "", err
	}

	for _, r := range rancher.Data {
		if strings.Compare(r.Name, env) == 0 {
			return r.ID, nil
		}
	}

	return "", errors.New("Can not find this env ID:" + env)
}

func (r Rancher) getService(envid string) ([]model.Service, error) {

}
