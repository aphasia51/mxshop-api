package main

import (
	"fmt"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
	myvalidator "mxshop-api/user-web/validator"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 初始化logger
	initialize.InitLogger()

	// 初始化配置文件
	initialize.InitConfig()

	// 初始化routers
	Router := initialize.Routers()

	// 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	// 初始化srv的连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	// 如果是本地开发环境 端口号固定  线上环境启动获取端口好
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation(
			"mobile",
			global.Trans,
			func(ut ut.Translator) error { return ut.Add("mobile", "{0} 非法的手机号码!", true) },
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("mobile", fe.Field())
				return t
			},
		)
	}

	/*
		1 S()可以获取一个全局的sugar， 可以让我们自己设置一个全局的logger
		2 日志是分级别的 debug, info, warn, error, fatal
		3 S函数和L()函数很有用，提供了一个全局的安全访问logger的途径
	*/

	zap.S().Debugf("启动服务, 端口: %d", global.ServerConfig.Port)

	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败: ", err.Error())
	}
}
