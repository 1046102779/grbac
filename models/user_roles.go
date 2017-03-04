package models

import (
	"time"

	utils "github.com/1046102779/common"
	. "github.com/1046102779/grbac/logger"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
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

func (t *UserRoles) InsertUserRoleNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v.%v] enter InsertUserRoleNoLock.", t.UserId, t.RoleId)
	defer Logger.Info("[%v.%v] left InsertUserRoleNoLock.", t.UserId, t.RoleId)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = utils.DB_DATA_ILLEGAL
		return
	}
	if _, err = (*o).Insert(t); err != nil {
		err = errors.Wrap(err, "InsertUserRoleNoLock.")
		retcode = utils.DB_INSERT_ERROR
		return
	}
	return
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

func isExistUserRoleByUserIdAndRoleId(userId int, roleId int) (isExist bool, retcode int, err error) {
	Logger.Info("[%v.%v] enter isExistUserRoleByUserIdAndRoleId.", userId, roleId)
	defer Logger.Info("[%v.%v] enter isExistUserRoleByUserIdAndRoleId.", userId, roleId)
	var (
		num int64
	)
	o := orm.NewOrm()
	num, err = o.QueryTable((&UserRoles{}).TableName()).Filter("user_id", userId).Filter("role_id", roleId).Filter("status", utils.STATUS_VALID).Count()
	if err != nil {
		err = errors.Wrap(err, "isExistUserRoleByUserIdAndRoleId")
		retcode = utils.DB_READ_ERROR
		return
	}
	isExist = false
	if num > 0 {
		isExist = true
	}
	return
}
