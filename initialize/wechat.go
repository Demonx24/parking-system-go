package initialize

import (
	"go.uber.org/zap"
	"log"
	"parking-system-go/config"
	"parking-system-go/global"
)

// InitWeChat 从 global.Config 读取微信配置，返回指针
func InitWeChat() *config.WeChat {
	if global.Config == nil {
		log.Fatalf("全局配置为空，无法初始化微信配置")
	}

	// 直接取 config.WeChat
	wechatCfg := global.Config.WeChat

	if wechatCfg.WeChatAppID == "" || wechatCfg.WeChatMchID == "" {
		log.Fatalf("微信配置不完整，请检查 config 文件")
	}

	global.Log.Info("微信配置初始化成功",
		zap.String("AppID", wechatCfg.WeChatAppID),
		zap.String("MchID", wechatCfg.WeChatMchID),
	)

	return &wechatCfg
}
