// 1. 权限多叉树解析获取功能ID
// 2. 修改功能
// 3. 删除功能
// 4. 功能列表与搜索
// 5. 新增功能
// 6. 解析URL，获取功能ID
package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/1046102779/grbac/common/consts"
	"github.com/1046102779/grbac/common/utils"
	"github.com/1046102779/grbac/conf"
	. "github.com/1046102779/grbac/logger"
	"github.com/1046102779/grbac/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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
			"err_code": consts.ERROR_CODE__PARAM__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	// 获取user_id和company_id
	companyId, retcode, err = utils.GetCompanyIdFromHeader(t.Ctx.Request)
	if err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	userId, retcode, err = utils.GetUserIdFromHeader(t.Ctx.Request)
	if err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
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
			result := conf.Redis__Client.SIsMember(fmt.Sprintf("YCFM_%d_%d", userId, funcId), companyId)
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
		result := conf.Redis__Client.SIsMember(fmt.Sprintf("YCFM_%d_%d", userId, funcId), companyId)
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

type FunctionInfo struct {
	Uri    string `json:"uri"`
	Method string `json:"method"`
	Name   string `json:"name"`
}

// 2. 修改功能
// @router /:id [PUT]
func (t *FunctionsController) ModifyFunction() {
	var (
		funcInfo *FunctionInfo = new(FunctionInfo)
	)
	// 获取功能ID
	funcId, _ := t.GetInt(":id")
	if funcId <= 0 {
		err := errors.New("param `:id` empty")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__PARAM__ILLEGAL,
			"err_msg":  err.Error(),
		}
		t.ServeJSON()
		return
	}
	// 解析json
	if err := jsoniter.Unmarshal(t.Ctx.Input.RequestBody, funcInfo); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__JSON__PARSE_FAILED,
			"err_msg":  err.Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	// 读取相关数据
	function := &models.Functions{
		Id: funcId,
	}
	if retcode, err := function.ReadFunctionNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	// 更新相关数据
	function.Uri = funcInfo.Uri
	function.MethodType = int16(models.GetMethodTypeByName(strings.ToUpper(funcInfo.Method)))
	function.Name = funcInfo.Name
	function.UpdatedAt = now
	if retcode, err := function.UpdateFunctionNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	t.Data["json"] = map[string]interface{}{
		"err_code": 0,
		"err_msg":  "",
	}
	t.ServeJSON()
	return
}

// 3. 删除功能
// @router /:id/invalid [PUT]
func (t *FunctionsController) DeleteFuntion() {
	funcId, _ := t.GetInt(":id")
	if funcId <= 0 {
		err := errors.New("param `:id` empty")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__PARAM__ILLEGAL,
			"err_msg":  err.Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	function := &models.Functions{
		Id: funcId,
	}
	if retcode, err := function.ReadFunctionNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	function.UpdatedAt = now
	function.Status = consts.STATUS_DELETED
	if retcode, err := function.UpdateFunctionNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	t.Data["json"] = map[string]interface{}{
		"err_code": 0,
		"err_msg":  "",
	}
	t.ServeJSON()
	return
}

// 4. 功能列表与搜索
// @router / [GET]
func (t *FunctionsController) GetFunctions() {
	pageIndex, _ := t.GetInt64("page_index", 1)
	pageSize, _ := t.GetInt64("page_size", 100)
	searchKey := t.GetString("search_key")
	if functions, count, realCount, retcode, err := models.GetFunctions(pageIndex-1, pageSize, searchKey); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	} else {
		t.Data["json"] = map[string]interface{}{
			"err_code":   0,
			"err_msg":    "",
			"count":      count,
			"real_count": realCount,
			"functions":  functions,
		}
	}
	t.ServeJSON()
	return
}

// 5. 新增功能
// @router / [POST]
func (t *FunctionsController) AddFunction() {
	var (
		functionInfo *FunctionInfo = new(FunctionInfo)
	)
	if err := jsoniter.Unmarshal(t.Ctx.Input.RequestBody, functionInfo); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__JSON__PARSE_FAILED,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	// 解析获取thid_region_mark_key
	// eg. /v1/storages/warehouses/:id/invalid -> :id
	markKey, retcode, err := models.GetThirdReginMarkKey(functionInfo.Uri)
	if err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	function := &models.Functions{
		Uri:                functionInfo.Uri,
		MethodType:         int16(models.GetMethodTypeByName(functionInfo.Method)),
		Name:               functionInfo.Name,
		ThirdRegionMarkKey: markKey,
		Status:             consts.STATUS_VALID,
		UpdatedAt:          now,
		CreatedAt:          now,
	}
	// 判断功能是否已经存在
	count, err := o.QueryTable(function.TableName()).Filter("uri", function.Uri).Filter("method_type", function.MethodType).Count()
	if err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__DB__READ,
			"err_msg":  err.Error(),
		}
		t.ServeJSON()
		return
	}
	if count > 0 {
		err = errors.New("`" + functionInfo.Method + "-" + functionInfo.Uri + "` already exist")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__PARAM__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	if retcode, err = function.InsertFunctionNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	t.Data["json"] = map[string]interface{}{
		"err_code": 0,
		"err_msg":  "",
	}
	t.ServeJSON()
	return
}
