package models

import (
	"reflect"
	"strings"
	"time"

	utils "github.com/1046102779/common"
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

func init() {
	orm.RegisterModel(new(RoleFunctionRelationships))
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
	num, err = o.QueryTable((&RoleFunctionRelationships{}).TableName()).Filter("role_id", id).Filter("status", utils.STATUS_VALID).All(&roleFuncs)
	if err != nil {
		err = errors.Wrap(err, "GetFuncIdsByRoleId")
		retcode = utils.DB_READ_ERROR
		return
	}
	if num > 0 {
		for index := 0; index < int(num); index++ {
			funcIds = append(funcIds, roleFuncs[index].FunctionId)
		}
	}
	return
}

// GetAllRoleFunctionRelationships retrieves all RoleFunctionRelationships matches certain condition. Returns empty list if
// no records exist
func GetAllRoleFunctionRelationships(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RoleFunctionRelationships))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []RoleFunctionRelationships
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}
