package config

type WeChat struct {
	WeChatAppID     string `yaml:"WeChatAppID"`
	WeChatMchID     string `yaml:"WeChatMchID"`
	WeChatAPIKey    string `yaml:"WeChatAPIKey"`
	WeChatNotifyURL string `yaml:"WeChatNotifyURL"`
}
