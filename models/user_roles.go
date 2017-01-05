package models

import (
	"errors"
	"reflect"
	"strings"
	"time"

	utils "github.com/1046102779/common"
	. "github.com/1046102779/grbac/logger"
	"github.com/astaxie/beego/orm"
)

type UserRoles struct {
	Id        int       `orm:"column(user_role_id);auto"`
	UserId    int       `orm:"column(user_id);null"`
	RoleId    int       `orm:"column(role_id);null"`
	RegionId  int       `orm:"column(region_id);null"`
	Status    int16     `orm:"column(status);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
}

func (t *UserRoles) TableName() string {
	return "user_roles"
}

func init() {
	orm.RegisterModel(new(UserRoles))
}

// 根据用户ID，获取角色ID列表
// input:  @param userId
// output: @param roleIds
func GetRoleIdsByUserId(id int) (roleIds []int, retcode int, err error) {
	Logger.Info("[%v] enter GetRoleIdsByUserId.", id)
	defer Logger.Info("[%v] left GetRoleIdsByUserId.", id)
	var (
		userRoles []UserRoles = []UserRoles{}
		num       int64
	)
	if id <= 0 {
		err = errors.New("param `user_id` empty")
		retcode = utils.SOURCE_DATA_ILLEGAL
		return
	}
	o := orm.NewOrm()
	num, err = o.QueryTable((&UserRoles{}).TableName()).Filter("user_id", id).Filter("status", utils.STATUS_VALID).All(&userRoles)
	if err != nil {
		Logger.Error(err.Error())
		return
	}
	if num > 0 {
		for index := 0; index < int(num); index++ {
			roleIds = append(roleIds, userRoles[index].RoleId)
		}
	}
	return
}

// GetAllUserRoles retrieves all UserRoles matches certain condition. Returns empty list if
// no records exist
func GetAllUserRoles(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserRoles))
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

	var l []UserRoles
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
