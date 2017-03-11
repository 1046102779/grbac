package models

import (
	"fmt"
	"strings"
	"time"

	"git.kissdata.com/ycfm/common/utils"

	"github.com/1046102779/grbac/common/consts"
	"github.com/1046102779/grbac/conf"
	pb "github.com/1046102779/grbac/igrpc"
	. "github.com/1046102779/grbac/logger"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
)

type UserRole struct {
	UserRoleId int       `orm:"column(user_role_id);auto" json:"user_role_id"`
	UserId     int       `orm:"column(user_id);null" json:"user_id"`
	RoleId     int       `orm:"column(role_id);null" json:"role_id"`
	RegionId   int       `orm:"column(region_id);null" json:"region_id"`
	Status     int16     `orm:"column(status);null" json:"status"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(datetime);null" json:"updated_at"`
	CreatedAt  time.Time `orm:"column(created_at);type(datetime);null" json:"created_at"`
}

func (t *UserRole) TableName() string {
	return "user_roles"
}

func AddUserRole(userId, roleId int) (id int, retcode int, err error) {
	Logger.Info("[%v.%v] enter AddUserRole.", userId, roleId)
	defer Logger.Info("[%v.%v] left AddUserRole.", userId, roleId)
	var (
		userRole *UserRole
	)
	userRole, retcode, err = GetUserRoleByRoleIdAndUserId(userId, roleId)
	if err != nil {
		err = errors.Wrap(err, "AddUserRole")
		return
	}
	if userRole != nil && userRole.UserRoleId > 0 {
		return userRole.UserRoleId, 0, nil
	}
	o := orm.NewOrm()
	now := time.Now()
	userRole = &UserRole{
		UserId:    userId,
		RoleId:    roleId,
		Status:    consts.STATUS_VALID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if retcode, err = userRole.InsertUserRoleNoLock(&o); err != nil {
		err = errors.Wrap(err, "AddUserRole")
		return
	}
	return userRole.UserRoleId, 0, nil
}

func GetUserRoleByRoleIdAndUserId(userId int, roleId int) (userRole *UserRole, retcode int, err error) {
	Logger.Info("[%v.%v] enter GetUserRoleByRoleIdAndUserId.", userId, roleId)
	defer Logger.Info("[%v.%v] left GetUserRoleByRoleIdAndUserId.", userId, roleId)
	var (
		userRoles []*UserRole
		num       int64 = 0
	)
	o := orm.NewOrm()
	num, err = o.QueryTable((&UserRole{}).TableName()).Filter("user_id", userId).Filter("role_id", roleId).Filter("status", consts.STATUS_VALID).All(&userRoles)
	if err != nil {
		err = errors.Wrap(err, "GetUserRoleByRoleIdAndUserId")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	if num > 0 {
		return userRoles[0], 0, nil
	}
	return
}

func (t *UserRole) ReadUserRoleNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter ReadUserRoleNoLock.", t.UserRoleId)
	defer Logger.Info("[%v] enter ReadUserRoleNoLock.", t.UserRoleId)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if err = (*o).Read(t); err != nil {
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	return
}

func (t *UserRole) UpdateUserRoleNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v.%v] enter UpdateUserRoleNoLock.", t.UserId, t.RoleId)
	defer Logger.Info("[%v.%v] left UpdateUserRoleNoLock.", t.UserId, t.RoleId)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Update(t); err != nil {
		retcode = consts.ERROR_CODE__DB__UPDATE
		return
	}
	return
}

func (t *UserRole) InsertUserRoleNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v.%v] enter InsertUserRoleNoLock.", t.UserId, t.RoleId)
	defer Logger.Info("[%v.%v] left InsertUserRoleNoLock.", t.UserId, t.RoleId)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Insert(t); err != nil {
		err = errors.Wrap(err, "InsertUserRoleNoLock.")
		retcode = consts.ERROR_CODE__DB__INSERT
		return
	}
	return
}

func init() {
	orm.RegisterModel(new(UserRole))
}

// 根据用户ID，获取角色ID列表
// input:  @param userId
// output: @param roleIds
func GetRoleIdsByUserId(id int) (roleIds []int, retcode int, err error) {
	Logger.Info("[%v] enter GetRoleIdsByUserId.", id)
	defer Logger.Info("[%v] left GetRoleIdsByUserId.", id)
	var (
		userRoles []UserRole = []UserRole{}
		num       int64
	)
	if id <= 0 {
		err = errors.New("param `user_id` empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	o := orm.NewOrm()
	num, err = o.QueryTable((&UserRole{}).TableName()).Filter("user_id", id).Filter("status", consts.STATUS_VALID).All(&userRoles)
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

type UserRoleInfos struct {
	Mobile                     string `json:"mobile"`
	Name                       string `json:"name"`
	RoleFunctionRelationshipId int64  `json:"role_function_relationship_id"`
	UserId                     int64  `json:"user_id"`
}

// 3. 用户与角色关系列表和搜索
// @param page_index: 页码(1+), page_size: 页面大小，role_id：角色ID，search_key={手机号，用户ID}
// 难点：1. 搜索手机号，跨服务rpc
//		 2. 手机号或者用户ID，需要组织两个列表，然后求Union集
//		 3. 分页
func GetUserRoles(pageIndex, pageSize int64, roleId int, searchKey string) (userRoleInfos []*UserRoleInfos, count int64, realCount int64, retcode int, err error) {
	Logger.Info("[%v.%v] enter GetUserRoles.", roleId, searchKey)
	defer Logger.Info("[%v.%v] left GetUserRoles.", roleId, searchKey)
	var (
		userIdFuzzy     int        = 0
		userIds         []int      = []int{}
		users                      = &pb.Users{}
		userRoles       []UserRole = []UserRole{}
		index, subIndex int        = 0, 0
		roles           []*Role
	)
	cond := orm.NewCondition()
	o := orm.NewOrm()
	qs := o.QueryTable((&UserRole{}).TableName())
	//cond = cond.And("status", consts.STATUS_VALID)
	if roleId > 0 {
		cond = cond.And("role_id", roleId)
	}
	if strings.TrimSpace(searchKey) != "" {
		// Union集
		if len(searchKey) <= 5 {
			// 如果长度小于5，可以搜索用户ID
			userIdFuzzy = utils.ConvertStrToInt(searchKey)
			cond = cond.Or("user_id__icontains", userIdFuzzy) // 条件一: OR查询
		}
		// rpc 服务，通过手机模糊匹配，获取<user_id, mobile>列表
		user := &pb.User{
			Mobile: searchKey,
		}
		conf.AccountClient.Call(fmt.Sprintf("%s.%s", "accounts", "GetUsersByFuzzyMobile"), user, users)
		for index = 0; users != nil && index < len(users.Users); index++ {
			userIds = append(userIds, int(users.Users[index].UserId))
		}
		if len(userIds) > 0 {
			cond = cond.Or("user_id__in", userIds) // 条件二：OR查询
		}
	}
	// 获取总记录数
	count, _ = qs.SetCond(cond).Filter("status", consts.STATUS_VALID).Count()
	// 获取分页实际页面填充大小，和页面数据
	realCount, _ = qs.SetCond(cond).Filter("status", consts.STATUS_VALID).Limit(pageSize, pageSize*pageIndex).All(&userRoles)

	// 现在还需要填充两种数据：1.手机号码；2.职位名称
	// 1. 获取用户手机号, 填充前端页面数据, 批量获取
	users = &pb.Users{}
	for index = 0; index < len(userRoles); index++ {
		users.Users = append(users.Users, &pb.User{
			UserId: int64(userRoles[index].UserId),
		})
	}
	conf.AccountClient.Call(fmt.Sprintf("%s.%s", "accounts", "GetUsersByUserIds"), users, users)
	for index = 0; index < len(userRoles); index++ {
		for subIndex = 0; users != nil && subIndex < len(users.Users); subIndex++ {
			if userRoles[index].UserId == int(users.Users[subIndex].UserId) {
				userRoleInfos = append(userRoleInfos, &UserRoleInfos{
					Mobile: users.Users[subIndex].Mobile,
					RoleFunctionRelationshipId: int64(userRoles[index].UserRoleId),
					UserId: int64(userRoles[index].UserId),
				})
			}
		}
	}
	// 2. 获取职位名称列表
	if roles, retcode, err = GetRoles(); err != nil {
		err = errors.Wrap(err, "GetUserRoles")
		return
	}
	for index = 0; index < len(userRoleInfos); index++ {
		userRoleInfos[index].Name, _, _ = GetRoleName(roles, int64(userRoles[index].RoleId))
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
	num, err = o.QueryTable((&UserRole{}).TableName()).Filter("user_id", userId).Filter("role_id", roleId).Filter("status", consts.STATUS_VALID).Count()
	if err != nil {
		err = errors.Wrap(err, "isExistUserRoleByUserIdAndRoleId")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	isExist = false
	if num > 0 {
		isExist = true
	}
	return
}
