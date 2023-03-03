package api

import (
	"context"
	"fmt"
	"math/rand"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GenerateSmsCode(width int) string {
	// 生成width长度的短信验证码
	numeric := [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	return sb.String()
}

func SendSms(ctx *gin.Context) {
	// 表单验证
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(err, ctx)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey(
		"cn-chengdu",
		global.ServerConfig.AliSmsInfo.ApiKey,
		global.ServerConfig.AliSmsInfo.ApiSecret,
	)
	if err != nil {
		panic(err)
	}

	code := global.ServerConfig.AliSmsInfo.Code
	smsCode := GenerateSmsCode(6)

	if len(smsCode) != 6 {
		smsCode = GenerateSmsCode(6)
	}

	request := requests.NewCommonRequest()                              // 构造一个公共请求
	request.Method = "POST"                                             // 设置请求方式
	request.Scheme = "https"                                            // https | http
	request.Domain = "dysmsapi.aliyuncs.com"                            // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2017-05-25"                                      // 指定产品版本
	request.ApiName = "SendSms"                                         // 指定接口名
	request.QueryParams["RegionId"] = "cn-chengdu"                      // 地区
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile            //手机号
	request.QueryParams["SignName"] = "Waka生鲜"                          //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = code                          //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成

	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			global.ServerConfig.RedisInfo.Host,
			global.ServerConfig.RedisInfo.Port,
		),
	})

	rdb.Set(
		context.Background(),
		sendSmsForm.Mobile,
		smsCode,
		time.Second*time.Duration(global.ServerConfig.RedisInfo.Expire),
	)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
