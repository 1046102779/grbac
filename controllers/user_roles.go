// 用户与角色管理
// 1. 删除用户与角色关系
// 2. 新增用户与角色关系
// 3. 用户与角色关系列表和搜索

package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/1046102779/grbac/common/consts"
	"github.com/1046102779/grbac/conf"
	pb "github.com/1046102779/grbac/igrpc"
	. "github.com/1046102779/grbac/logger"
	"github.com/1046102779/grbac/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// UserRolesController operations for UserRoles
type UserRolesController struct {
	beego.Controller
}

// 1. 删除用户与角色关系
// @router /:id/invalid [PUT]
func (t *UserRolesController) DeleteUserRoles() {
	id, _ := t.GetInt(":id")
	if id <= 0 {
		err := errors.New("param `:id` empty")
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__SOURCE_DATA__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	userRole := &models.UserRole{
		UserRoleId: id,
	}
	if retcode, err := userRole.ReadUserRoleNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	userRole.Status = consts.STATUS_DELETED
	userRole.UpdatedAt = now
	if retcode, err := userRole.UpdateUserRoleNoLock(&o); err != nil {
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

// 2. 新增用户与角色关系
// @router / [POST]
func (t *UserRolesController) InsertUserRole() {
	type UserInfo struct {
		Mobile string `json:"mobile"`
		RoleId int    `json:"role_id"`
	}
	var (
		info *UserInfo = new(UserInfo)
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
	if strings.TrimSpace(info.Mobile) == "" || info.RoleId <= 0 {
		err := errors.New("param `mobile | role_id` empty")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__SOURCE_DATA__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	// rpc服务，通过手机号，获取用户ID
	user := &pb.User{
		Mobile: info.Mobile,
	}
	conf.AccountClient.Call(fmt.Sprintf("%s.%s", "accounts", "GetUserByMobile"), user, user)
	// 通过角色名称， 获取角色ID
	// 新增角色与用户的关系
	o := orm.NewOrm()
	now := time.Now()
	userRole := &models.UserRole{
		UserId:    int(user.UserId),
		RoleId:    info.RoleId,
		Status:    consts.STATUS_VALID,
		UpdatedAt: now,
		CreatedAt: now,
	}
	if retcode, err := userRole.InsertUserRoleNoLock(&o); err != nil {
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

// 3. 用户与角色关系列表和搜索
// @router / [GET]
func (t *UserRolesController) GetUserRoles() {
	pageIndex, _ := t.GetInt64("page_index", 1)
	pageSize, _ := t.GetInt64("page_size", 100)
	roleId, _ := t.GetInt("role_id")
	searchKey := t.GetString("search_key")
	userRoles, count, realCount, retcode, err := models.GetUserRoles(pageIndex-1, pageSize, roleId, searchKey)
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
		"err_code":   0,
		"err_msg":    "",
		"count":      count,
		"real_count": realCount,
		"user_roles": userRoles,
	}
	t.ServeJSON()
	return
}
