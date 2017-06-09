package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// RancherProjectInfo 每个Project/Service详细信息
// 当为Project时保存ID,Name两项信息
// 当为Service时保存ID,Name,State三项信息
type RancherProjectInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state,omitempty"`
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

	log.Printf("Env [%s] ID [%s] \n", r.Env, envID)

	return r.getService(envID)
}

// Start 启动Rancher中指定服务
func (r Rancher) Start(srv model.Service) error {
	eid, err := r.getEnvID(r.Env)
	if err != nil {
		return err
	}

	sid, err := r.getSrvID(eid, r.Service)
	if err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s%s%s/services/%s?action=activate", r.Domain, ENVAPI, eid, sid)
	client := &http.Client{}
	// log.Println(reqURL)
	req, err := http.NewRequest("POST", reqURL, nil)
	req.SetBasicAuth(r.AccessKey, r.SecretKey)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode > 300 {
		return errors.New("Service Activate Error. " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

// Restart 重启Rancher中的指定服务
func (r Rancher) Restart(srv model.Service) error {
	eid, err := r.getEnvID(r.Env)
	if err != nil {
		return err
	}

	sid, err := r.getSrvID(eid, r.Service)
	if err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s%s%s/services/%s?action=restart", r.Domain, ENVAPI, eid, sid)
	client := &http.Client{}

	// 固定重启策略
	req, err := http.NewRequest("POST", reqURL, strings.NewReader("{ \"rollingRestartStrategy\": { \"batchSize\": 1, \"intervalMillis\": 2000 }}"))
	req.SetBasicAuth(r.AccessKey, r.SecretKey)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode > 300 {
		return errors.New("Service Restart Error. " + strconv.Itoa(resp.StatusCode))
	}

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

	return "", errors.New("Can not find this env ID: " + env)
}

// getService 获取指定服务列表
// envid 服务所在的Env ID
func (r Rancher) getService(envid string) ([]model.Service, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", r.Domain+ENVAPI+envid+"/service", nil)
	req.SetBasicAuth(r.AccessKey, r.SecretKey)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var srvs RancherProject

	err = json.Unmarshal(content, &srvs)
	if err != nil {
		return nil, err
	}

	var msrv []model.Service

	for _, s := range srvs.Data {
		var ms = model.Service{
			Name:   s.Name,
			Status: s.State,
			Type:   "rancher service",
			SrvID:  s.ID,
		}

		msrv = append(msrv, ms)
	}

	return msrv, nil
}

// getSrvID 获取指定服务ID
// envid 服务所在的Env ID
// srvname 服务名称
func (r Rancher) getSrvID(envid, srvname string) (string, error) {
	srvs, err := r.getService(envid)
	if err != nil {
		return "", err
	}

	for _, s := range srvs {
		if strings.Compare(s.Name, srvname) == 0 {
			return s.SrvID, nil
		}
	}

	return "", errors.New("Can not find this service ID")
}
