package controllers

import (
	"fmt"
	"strconv"

	utils "github.com/1046102779/common"
	. "github.com/1046102779/common/utils"
	. "github.com/1046102779/grbac/logger"
	"github.com/1046102779/grbac/models"
	"github.com/astaxie/beego"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// FunctionsController operations for Functions
type FunctionsController struct {
	beego.Controller
}

// 权限分两步解析:
// >> 1. 获取功能ID
// >> 2. 判断功能和实体ID，是否有对应的用户权限

// >> 1. 获取功能ID
// @router / [POST]
func (t *FunctionsController) GetFuncId() {
	var (
		info                       *models.HttpRequestInfo = new(models.HttpRequestInfo)
		funcId                     int
		entityStr                  string
		entityId                   int64
		err                        error
		userId, companyId, retcode int
	)
	if err = jsoniter.Unmarshal(t.Ctx.Input.RequestBody, info); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": utils.SOURCE_DATA_ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	// 获取user_id和company_id
	if info, retcode, err := GetHeaderParams(t.Ctx.Request); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	} else if info != nil && info.CompanyId > 0 {
		companyId = info.CompanyId
		userId = info.UserId
	} else {
		err := errors.New("please login homepage")
		t.Data["json"] = map[string]interface{}{
			"err_code": utils.USER_LOGGED_IN,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	fmt.Println("info: ", *info)
	// 1. 先判断URL是否在白名单列表
	// 2. 获取URL在构建树中是否存在
	if funcId, entityStr, retcode, err = models.GetFuncId(info); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	// 根据解析树的结果，分析判断
	//   funcId, entityStr
	// 1.  -1 ,    ""    // 没有匹配到
	// 2.  >0 ,    ""    // 找到，但是URI.PATH中没有实体ID
	// 3.  >0 ,    !=""  // 找到，且有实体ID，但是还需要对entityStr判断，看是否可以转化为int型，否则归为第二种情况
	if funcId == -1 && entityStr == "" {
		t.Ctx.Output.SetStatus(403)
		return
	}
	if funcId > 0 && entityStr != "" {
		if entityId, err = strconv.ParseInt(entityStr, 10, 64); err != nil {
			entityStr = "" // 无法转化为int型，归为第二种情况
		} else {
			// 判断用户ID持有的<funcId, entityId>，在<funcId, entityIds> 列表中是否有对应的entityId存在
			// 如果有，则表示该用户有访问entityId资源的权限
			result := RedisClient.SIsMember(fmt.Sprintf("YCFM_%d_%d", userId, funcId), companyId)
			if result != nil || result.Err() != nil || !result.Val() {
				if result.Err() != nil {
					Logger.Error(err.Error())
				}
				t.Ctx.Output.SetStatus(403) // 没有访问权限
				return
			}
			t.Ctx.Output.SetStatus(200) // 有访问权限
			return
		}
		fmt.Println("hello,world, entityId=", entityId)
		return
	}
	if funcId > 0 && entityStr == "" {
		// 判断用户ID持有的funcId, 在<funcId, entityIds> 列表中对应的entityIds为空
		// 如果为空，则表示该用户只要访问这个URI，一定会有访问权限
		// 否则，没有访问权限
		result := RedisClient.SIsMember(fmt.Sprintf("YCFM_%d_%d", userId, funcId), companyId)
		if result != nil || result.Err() != nil || !result.Val() {
			if result.Err != nil {
				Logger.Error(err.Error())
			}
			t.Ctx.Output.SetStatus(403) // 没有访问权限
			return
		}
		t.Ctx.Output.SetStatus(200) // 有访问权限
		return
	}
	return
}
