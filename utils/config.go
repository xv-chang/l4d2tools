package utils

import (
	"encoding/json"
	"log"
)

type WebConfig struct {
	Auth      string `json:"auth"`
	AddonsDir string `json:"addons_dir"`
	WebUrl    string `json:"web_url"`
}

func DefaultConfig() *WebConfig {
	return &WebConfig{}
}
func LoadConfig(p string) *WebConfig {
	if !PathExists(p) {
		log.Fatalf("尝试加载配置文件 %v 失败: 文件不存在", p)
		return nil
	}
	c := WebConfig{}
	err := json.Unmarshal([]byte(ReadAllText(p)), &c)
	if err != nil {
		log.Fatalf("尝试加载配置文件 %v 时出现错误: %v", p, err)
		return nil
	}
	return &c
}
