package models

import (
	"strings"
	"time"

	"github.com/1046102779/grbac/common/consts"
	. "github.com/1046102779/grbac/logger"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
)

type RoleFunctionRelationships struct {
	Id         int       `orm:"column(role_function_relationship_id);auto"`
	RoleId     int       `orm:"column(role_id);null"`
	FunctionId int       `orm:"column(function_id);null"`
	RegionId   int       `orm:"column(region_id);null"`
	Status     int16     `orm:"column(status);null"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(datetime);null"`
	CreatedAt  time.Time `orm:"column(created_at);type(datetime);null"`
}

func (t *RoleFunctionRelationships) TableName() string {
	return "role_function_relationships"
}

func (t *RoleFunctionRelationships) ReadRoleFunctionNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter ReadRoleFunctionNoLock.", t.Id)
	defer Logger.Info("[%v] left ReadRoleFunctionNoLock.", t.Id)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if err = (*o).Read(t); err != nil {
		err = errors.Wrap(err, "ReadRoleFunctionNoLock")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	return
}

func (t *RoleFunctionRelationships) UpdateRoleFunctionNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter UpdateRoleFunctionNoLock.", t.Id)
	defer Logger.Info("[%v] left UpdateRoleFunctionNoLock.", t.Id)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Update(t); err != nil {
		err = errors.Wrap(err, "UpdateRoleFunctionNoLock")
		retcode = consts.ERROR_CODE__DB__UPDATE
		return
	}
	return
}

func (t *RoleFunctionRelationships) InsertRoleFunctionNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v.%v] enter InsertRoleFunctionNoLock.", t.RoleId, t.FunctionId)
	defer Logger.Info("[%v.%v] left InsertRoleFunctionNoLock.", t.RoleId, t.FunctionId)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Insert(t); err != nil {
		err = errors.Wrap(err, "InsertRoleFunctionNoLock")
		retcode = consts.ERROR_CODE__DB__INSERT
		return
	}
	return
}

func init() {
	orm.RegisterModel(new(RoleFunctionRelationships))
}

func GetRoleFunctions(pageIndex int64, pageSize int64, roleId int, searchKey string) (roleFunctions []*RoleFunctionRelationships, count int64, realCount int64, retcode int, err error) {
	Logger.Info("[%v.%v] enter GetRoleFunctions.", roleId, searchKey)
	defer Logger.Info("[%v.%v] left GetRoleFunctions.", roleId, searchKey)
	var (
		functionIds []int            = []int{}
		functions   []*FunctionInfos = []*FunctionInfos{}
	)
	o := orm.NewOrm()
	qs := o.QueryTable((&RoleFunctionRelationships{}).TableName()).Filter("status", consts.STATUS_VALID)
	if roleId > 0 {
		qs = qs.Filter("role_id", roleId)
	}
	if strings.TrimSpace(searchKey) != "" {
		// search_key可以是 功能名称或者URI
		functions, _, _, retcode, err = GetFunctions(0, 10000, searchKey)
		if err != nil {
			err = errors.Wrap(err, "GetRoleFunctions")
			return
		}
		for index := 0; index < len(functions); index++ {
			functionIds = append(functionIds, functions[index].Id)
		}
		if len(functionIds) > 0 {
			qs = qs.Filter("function_id__in", functionIds)
		} else {
			// 搜索结果为空
			count = 0
			realCount = 0
			return
		}
	}
	count, _ = qs.Count()
	realCount, _ = qs.Limit(pageSize, pageIndex*pageSize).All(&roleFunctions)
	return
}

// 根据角色ID， 获取功能ID列表
func GetFuncIdsByRoleId(id int) (funcIds []int, retcode int, err error) {
	Logger.Info("[%v] enter GetFuncIdsByRoleId.", id)
	defer Logger.Info("[%v] left GetFuncIdsByRoleId.", id)
	var (
		roleFuncs []RoleFunctionRelationships = []RoleFunctionRelationships{}
		num       int64
	)
	o := orm.NewOrm()
	num, err = o.QueryTable((&RoleFunctionRelationships{}).TableName()).Filter("role_id", id).Filter("status", consts.STATUS_VALID).All(&roleFuncs)
	if err != nil {
		err = errors.Wrap(err, "GetFuncIdsByRoleId")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	if num > 0 {
		for index := 0; index < int(num); index++ {
			funcIds = append(funcIds, roleFuncs[index].FunctionId)
		}
	}
	return
}
