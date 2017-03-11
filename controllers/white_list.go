// 白名单管理:
// 1. 修改白名单
// 2. 删除白名单
// 3. 新增白名单
// 4. 白名单列表
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

// WhiteListController operations for WhiteList
type WhiteListController struct {
	beego.Controller
}

type WhiteListInfo struct {
	Method string `json:"method"`
	Name   string `json:"name"`
	Uri    string `json:"uri"`
}

// 1. 修改白名单
// @router /:id [PUT]
func (t *WhiteListController) ModifyWhitList() {
	var (
		whiteListInfo *WhiteListInfo = new(WhiteListInfo)
	)
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
	if err := jsoniter.Unmarshal(t.Ctx.Input.RequestBody, whiteListInfo); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__JSON__PARSE_FAILED,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	whiteList := &models.WhiteList{
		Id: id,
	}
	if retcode, err := whiteList.ReadWhiteListNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	whiteList.UpdatedAt = now
	whiteList.Url = whiteListInfo.Uri
	whiteList.Name = whiteListInfo.Name
	whiteList.MethodType = int16(models.GetMethodTypeByName(whiteListInfo.Method))
	if retcode, err := whiteList.UpdateWhiteListNoLock(&o); err != nil {
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

// 2. 删除白名单
// @router /:id/invalid [PUT]
func (t *WhiteListController) DeleteWhiteList() {
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
	whiteList := &models.WhiteList{
		Id: id,
	}
	if retcode, err := whiteList.ReadWhiteListNoLock(&o); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	whiteList.UpdatedAt = now
	whiteList.Status = consts.STATUS_DELETED
	if retcode, err := whiteList.UpdateWhiteListNoLock(&o); err != nil {
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

// 3. 新增白名单
// @router / [POST]
func (t *WhiteListController) AddWhiteList() {
	var (
		whiteListInfo *WhiteListInfo = new(WhiteListInfo)
	)
	if err := jsoniter.Unmarshal(t.Ctx.Input.RequestBody, whiteListInfo); err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__SOURCE_DATA__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	if strings.TrimSpace(whiteListInfo.Method) == "" || strings.TrimSpace(whiteListInfo.Name) == "" ||
		strings.TrimSpace(whiteListInfo.Uri) == "" {
		err := errors.New("param `method | name | uri` empty")
		t.Data["json"] = map[string]interface{}{
			"err_code": consts.ERROR_CODE__SOURCE_DATA__ILLEGAL,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	o := orm.NewOrm()
	now := time.Now()
	whiteList := &models.WhiteList{
		Name:       whiteListInfo.Name,
		Url:        whiteListInfo.Uri,
		MethodType: int16(models.GetMethodTypeByName(whiteListInfo.Method)),
		Status:     consts.STATUS_VALID,
		UpdatedAt:  now,
		CreatedAt:  now,
	}
	if retcode, err := whiteList.InsertWhiteListNoLock(&o); err != nil {
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

// 4. 白名单列表
// @router / [GET]
func (t *WhiteListController) GetWhiteLists() {
	type WhiteListInfo struct {
		Method      string `json:"method"`
		Name        string `json:"name"`
		Uri         string `json:"uri"`
		WhiteListId int    `json:"white_list_id"`
	}
	var (
		whiteListInfos []*WhiteListInfo = []*WhiteListInfo{}
	)
	pageIndex, _ := t.GetInt64("page_index", 1)
	pageSize, _ := t.GetInt64("page_size", 100)
	whiteList, count, realCount, retcode, err := models.GetWhiteLists(pageIndex-1, pageSize)
	if err != nil {
		Logger.Error(err.Error())
		t.Data["json"] = map[string]interface{}{
			"err_code": retcode,
			"err_msg":  errors.Cause(err).Error(),
		}
		t.ServeJSON()
		return
	}
	for index := 0; whiteList != nil && index < len(whiteList); index++ {
		whiteListInfos = append(whiteListInfos, &WhiteListInfo{
			Method:      models.GetMethodNameByType(int(whiteList[index].MethodType)),
			Name:        whiteList[index].Name,
			Uri:         whiteList[index].Url,
			WhiteListId: whiteList[index].Id,
		})
	}
	t.Data["json"] = map[string]interface{}{
		"err_code":   0,
		"err_msg":    "",
		"count":      count,
		"real_count": realCount,
		"white_list": whiteListInfos,
	}
	t.ServeJSON()
	return
}
