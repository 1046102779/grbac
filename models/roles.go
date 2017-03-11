package models

import (
	"strings"
	"time"

	"github.com/1046102779/grbac/common/consts"
	. "github.com/1046102779/grbac/logger"

	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
)

type Role struct {
	RoleId    int       `orm:"column(role_id);auto" json:"role_id"`
	RegionId  int       `orm:"column(region_id);null"`
	Name      string    `orm:"column(name);size(50);null" json:"name"`
	Code      string    `orm:"column(code);size(20);null" json:"code"`
	Status    int16     `orm:"column(status);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
}

func (t *Role) TableName() string {
	return "roles"
}

func GetRoleByRoleCode(code string, o *orm.Ormer) (role Role, retCode int, err error) {
	if o == nil {
		oo := orm.NewOrm()
		o = &oo
	}

	role = Role{
		Code: code,
	}
	if err = (*o).Read(&role, "Code"); err != nil {
		if err == orm.ErrNoRows {
			retCode = consts.ERROR_CODE__ROLE__NOT_EXIST
		} else {
			retCode = consts.ERROR_CODE__DB__READ
		}
		err = errors.Wrap(err, "GetRoleByCode")
	}
	return
}

func (t *Role) UpdateRoleNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter UpdateRoleNoLock.", t.RoleId)
	defer Logger.Info("[%v] left UpdateRoleNoLock.", t.RoleId)
	if o == nil {
		err = errors.New("param `orm.Ormer` empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Update(t); err != nil {
		err = errors.Wrap(err, "UpdateRoleNoLock")
		retcode = consts.ERROR_CODE__DB__UPDATE
		return
	}
	return
}

func (t *Role) InsertRoleNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter InsertRoleNoLock.", t.Name)
	defer Logger.Info("[%v] left InsertRoleNoLock.", t.Name)
	if o == nil {
		err = errors.New("param `orm.Ormer` empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Insert(t); err != nil {
		err = errors.Wrap(err, "InsertRoleNoLock")
		retcode = consts.ERROR_CODE__DB__INSERT
		return
	}
	return
}

func (t *Role) ReadRoleNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter ReadRoleNoLock.", t.RoleId)
	defer Logger.Info("[%v] enter ReadRoleNoLock.", t.RoleId)
	if o == nil {
		err = errors.New("param `orm.Ormer` empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if err = (*o).Read(t); err != nil {
		err = errors.Wrap(err, "ReadRoleNoLock")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	return
}

func init() {
	orm.RegisterModel(new(Role))
}

// 获取<角色ID， 角色名称> 列表
func GetRoles() (roles []*Role, retcode int, err error) {
	Logger.Info("enter GetRoleInfos.")
	defer Logger.Info("left GetRoleInfos.")
	roles = []*Role{}
	o := orm.NewOrm()
	_, err = o.QueryTable((&Role{}).TableName()).Filter("status", consts.STATUS_VALID).All(&roles)
	if err != nil {
		err = errors.Wrap(err, "GetRoleInfos")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	return
}

func GetRoleName(roles []*Role, roleId int64) (name string, retcode int, err error) {
	Logger.Info("[%v] enter GetRoleName.", roleId)
	defer Logger.Info("[%v] left GetRoleName.", roleId)
	if roles == nil {
		return
	}
	for index := 0; index < len(roles); index++ {
		if roles[index].RoleId == int(roleId) {
			return roles[index].Name, 0, nil
		}
	}
	return
}

// 通过角色名称，获取角色ID
func GetRoleIdByName(name string) (id int, retcode int, err error) {
	Logger.Info("[%v] enter GetRoleIdByName.", name)
	defer Logger.Info("[%v] left GetRoleIdByName.", name)
	var (
		roles []Role = []Role{}
		num   int64
	)
	if strings.TrimSpace(name) == "" {
		err = errors.New("param `name` empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	o := orm.NewOrm()
	num, err = o.QueryTable((&Role{}).TableName()).Filter("name", name).All(&roles)
	if err != nil {
		err = errors.Wrap(err, "GetRoleIdByName")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	if num > 0 {
		id = roles[0].RoleId
	}
	return
}

func getRoleIdByCode(code string) (id int, retcode int, err error) {
	Logger.Info("[%v] enter getRoleIdByCode.", code)
	defer Logger.Info("[%v] left getRoleIdByCode.", code)
	var role *Role = new(Role)
	o := orm.NewOrm()
	err = o.QueryTable((&Role{}).TableName()).Filter("code__icontains", code).One(&role)
	if err != nil {
		err = errors.Wrap(err, "getRoleIdByCode")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	id = role.RoleId
	return
}
