// 角色管理
// 1. 修改角色
// 2. 新增角色
// 3. 获取角色列表
package controllers

import (
	"strings"
	"time"

	"github.com/1046102779/grbac/common/consts"
	. "github.com/1046102779/grbac/logger"
	"github.com/1046102779/grbac/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// RolesController operations for Roles
type RolesController struct {
	beego.Controller
}

type RoleInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// 1. 修改角色
// @router /:id [PUT]
func (t *RolesController) ModifyRoles() {
	var (
		info *RoleInfo = new(RoleInfo)
	)
	// 获取角色id
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
	// json解析
	if err := jsoniter.Unmarshal(t.Ctx.Input.RequestBody, info); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__JSON__PARSE_FAILED,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	if strings.TrimSpace(info.Code) == "" || strings.TrimSpace(info.Name) == "" {
		err := errors.New("param `code | name` empty")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__SOURCE_DATA__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	role := &models.Role{
		RoleId: id,
	}
	if retcode, err := role.ReadRoleNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	role.Code = info.Code
	role.Name = info.Name
	role.UpdatedAt = now
	if retcode, err := role.UpdateRoleNoLock(&o); err != nil {
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

// 2. 新增角色
// @router / [POST]
func (t *RolesController) AddRole() {
	var (
		info *RoleInfo = new(RoleInfo)
	)
	if err := jsoniter.Unmarshal(t.Ctx.Input.RequestBody, info); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__JSON__PARSE_FAILED,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	if strings.TrimSpace(info.Name) == "" {
		err := errors.New("param `name` empty")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__JSON__PARSE_FAILED,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	role := &models.Role{
		Code:      info.Code,
		Name:      info.Name,
		Status:    consts.STATUS_VALID,
		UpdatedAt: now,
		CreatedAt: now,
	}
	if retcode, err := role.InsertRoleNoLock(&o); err != nil {
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

// 3. 获取角色列表
// @router / [GET]
func (t *RolesController) GetRoles() {
	roles, retcode, err := models.GetRoles()
	if err != nil {
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
		"roles":    roles,
	}
	t.ServeJSON()
	return
}
