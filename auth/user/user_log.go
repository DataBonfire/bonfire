package user

import "github.com/databonfire/bonfire/resource"

type UserLog struct {
	resource.Model

	// 用户
	UserId uint `json:"user_id" gorm:"index" react_type:"reference" react_reference:"users.nick_name"`
	// .eg: GET./offers
	MethodPath string `json:"method_path" gorm:"type:varchar(255);index"`
	// http code
	ResponseCode int `json:"response_code"`
	// int: timestamp
	AccessAt int64 `json:"access_at"`
	// 请求消耗时间，单位为 毫秒
	CostTime int64 `json:"cost_time"`
	// 客户端的UserAgent
	UA string `json:"ua"`
	// 发起请求的客户端的IP
	IP string `json:"ip"`
	// 请求协议
	Protocol string `json:"protocol" gorm:"type:varchar(40)" `
	// 请求URI
	RequestUri string `json:"request_uri"`
	// string: multiline
	RequestBody string `json:"request_body"`
}
