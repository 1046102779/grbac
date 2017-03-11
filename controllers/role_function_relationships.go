// 角色与功能关系管理
// 1. 删除角色与功能关系
// 2. 新增角色与功能关系
// 3. 角色与功能关系列表和搜索

package controllers

import (
	"time"

	"github.com/1046102779/grbac/common/consts"
	. "github.com/1046102779/grbac/logger"
	"github.com/1046102779/grbac/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// RoleFunctionRelationshipsController operations for RoleFunctionRelationships
type RoleFunctionRelationshipsController struct {
	beego.Controller
}

// 1. 删除角色与功能关系
// @router /:id/invalid [PUT]
func (t *RoleFunctionRelationshipsController) DeleteRoleFunction() {
	id, _ := t.GetInt(":id")
	if id <= 0 {
		err := errors.New("param `:id` empty")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__SOURCE_DATA__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	roleFunction := &models.RoleFunctionRelationships{
		Id: id,
	}
	if retcode, err := roleFunction.ReadRoleFunctionNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	roleFunction.Status = consts.STATUS_DELETED
	roleFunction.UpdatedAt = now
	if retcode, err := roleFunction.UpdateRoleFunctionNoLock(&o); err != nil {
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

// 2. 新增角色与功能关系
// @router / [POST]
func (t *RoleFunctionRelationshipsController) AddRoleFunction() {
	// Id: function id, Name: role name
	type RoleFunctionInfo struct {
		Id     int `json:"function_id"`
		RoleId int `json:"role_id"`
	}
	var (
		roleFunctionInfo *RoleFunctionInfo = new(RoleFunctionInfo)
	)
	if err := jsoniter.Unmarshal(t.Ctx.Input.RequestBody, roleFunctionInfo); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__JSON__PARSE_FAILED,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	if roleFunctionInfo.Id <= 0 || roleFunctionInfo.RoleId <= 0 {
		err := errors.New("param `function_id | role_id` empty")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__SOURCE_DATA__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	roleFunction := &models.RoleFunctionRelationships{
		RoleId:     roleFunctionInfo.RoleId,
		FunctionId: roleFunctionInfo.Id,
		Status:     consts.STATUS_VALID,
		UpdatedAt:  now,
		CreatedAt:  now,
	}
	if retcode, err := roleFunction.InsertRoleFunctionNoLock(&o); err != nil {
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

// 3. 角色与功能关系列表和搜索
// @router / [GET]
func (t *RoleFunctionRelationshipsController) GetRoleFunctions() {
	type RoleFunctionInfo struct {
		FunctionName               string `json:"function_name"`
		Method                     string `json:"method"`
		RoleFunctionRelationshipId int    `json:"role_function_relationship_id"`
		RoleName                   string `json:"role_name"`
		Uri                        string `json:"uri"`
	}
	var (
		roleFunctionInfos []*RoleFunctionInfo = []*RoleFunctionInfo{}
		roles             []*models.Role      = []*models.Role{}
	)
	pageIndex, _ := t.GetInt64("page_index", 1)
	pageSize, _ := t.GetInt64("page_size", 100)
	roleId, _ := t.GetInt("role_id")
	searchKey := t.GetString("search_key")
	roleFunctions, count, realCount, retcode, err := models.GetRoleFunctions(pageIndex-1, pageSize, roleId, searchKey)
	if err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	if roles, retcode, err = models.GetRoles(); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	for index := 0; roleFunctions != nil && index < len(roleFunctions); index++ {
		function := &models.Functions{
			Id: roleFunctions[index].FunctionId,
		}
		if retcode, err = function.ReadFunctionNoLock(&o); err != nil {
			Logger.Error(err.Error())
			t.Data["json"] = map[string]interface{}{
				"err_code": retcode,
				"err_msg":  errors.Cause(err).Error(),
			}
			t.ServeJSON()
			return
		}
		roleName, _, _ := models.GetRoleName(roles, int64(roleFunctions[index].RoleId))
		roleFunctionInfos = append(roleFunctionInfos, &RoleFunctionInfo{
			RoleFunctionRelationshipId: roleFunctions[index].Id,
			FunctionName:               function.Name,
			Method:                     models.GetMethodNameByType(int(function.MethodType)),
			RoleName:                   roleName,
			Uri:                        function.Uri,
		})
	}
	t.Data["json"] = map[string]interface{}{
		"err_code":       0,
		"err_msg":        "",
		"count":          count,
		"real_count":     realCount,
		"role_functions": roleFunctionInfos,
	}
	t.ServeJSON()
	return
}
