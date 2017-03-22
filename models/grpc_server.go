// rpcx服务列表
// 1. 加载用户关系数据到redis中
// 2. 新增系统账户，修改系统或者员工角色时，需要修改用户与角色之间的关系
package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"

	"github.com/1046102779/grbac/common/consts"
	"github.com/1046102779/grbac/conf"
	pb "github.com/1046102779/grbac/igrpc"
	. "github.com/1046102779/grbac/logger"
)

type GrbacServer struct{}

func (t *GrbacServer) LoadGrbacUserRel(in *pb.GrbacUserRel, out *pb.GrbacUserRel) (err error) {
	Logger.Info("[%v] enter LoadGrbacUserRel.", in.UserId)
	defer Logger.Info("[%v] left LoadGrbacUserRel.", in.UserId)
	var (
		roleIds []int    // 用户在角色ID列表
		funcIds []int    // 用户在功能ID列表
		keys    []string // redis正则表达式获取的keys列表
	)
	defer func() {
		err = nil
	}()
	if roleIds, _, err = GetRoleIdsByUserId(int(in.UserId)); err != nil {
		Logger.Error(err.Error())
		return
	}
	// 清空用户权限数据，然后重新建立
	if keys, err = conf.Redis__Client.Keys(fmt.Sprintf("YCFM_%d_*", in.UserId)); err != nil {
		Logger.Error(err.Error())
		return
	}
	fmt.Println("keys: ", keys)
	if err = conf.Redis__Client.DelKeys(keys); err != nil {
		Logger.Error(err.Error())
	}
	// 重新建立
	for index := 0; roleIds != nil && index < len(roleIds); index++ {
		funcIds, _, err = GetFuncIdsByRoleId(roleIds[index])
		// 添加key，三元表<用户ID-功能ID，实体ID集合={公司ID}>
		if err != nil {
			Logger.Error(err.Error())
		}
		for subIndex := 0; funcIds != nil && subIndex < len(funcIds); subIndex++ {
			conf.Redis__Client.SAdd(fmt.Sprintf("YCFM_%d_%d", in.UserId, funcIds[index]), in.CompanyId)
		}
	}
	return
}

// 2. 新增系统账户，修改系统或者员工角色时，需要修改用户与角色之间的关系
func (t *GrbacServer) ModifyUserRole(in *pb.GrbacUserRel, out *pb.GrbacUserRel) (err error) {
	Logger.Info("[%v] enter ModifyUserRole.", in.CompanyId)
	defer Logger.Info("[%v] left ModifyUserRole.", in.CompanyId)
	var (
		isExist bool = false
		roleId  int
	)
	defer func() {
		err = nil
	}()
	if in.CompanyId <= 0 || in.UserId <= 0 || in.Code == "" {
		return
	}
	if roleId, _, err = getRoleIdByCode(in.Code); err != nil {
		Logger.Error(err.Error())
		return
	}
	if roleId <= 0 {
		return
	}
	if isExist, _, err = isExistUserRoleByUserIdAndRoleId(int(in.UserId), roleId); err != nil {
		Logger.Error(err.Error())
		return
	}
	if !isExist {
		// 不存在，则添加用户与角色关系记录
		now := time.Now()
		o := orm.NewOrm()
		userRole := &UserRole{
			UserId:    int(in.UserId),
			RoleId:    roleId,
			Status:    consts.STATUS_VALID,
			UpdatedAt: now,
			CreatedAt: now,
		}
		if _, err = userRole.InsertUserRoleNoLock(&o); err != nil {
			Logger.Error(err.Error())
			return
		}
	}
	return
}

func (s *GrbacServer) AddUserRole(in *pb.UserRole, out *pb.UserRole) (err error) {
	defer func() { err = nil }()
	var userRoleId int = 0
	if userRoleId, _, err = AddUserRole(int(in.UserId), int(in.RoleId)); err != nil {
		err = errors.Wrap(err, "AddUserRole, GrbacServer")
		return
	}

	out = in
	out.UserRoleId = int32(userRoleId)
	return
}

func (s *GrbacServer) DelUserRole(in *pb.UserRole, out *pb.UserRole) (err error) {
	Logger.Info("[%v.%v] enter DelUserRole.", in.UserId, in.RoleId)
	defer Logger.Info("[%v.%v] enter DelUserRole.", in.UserId, in.RoleId)
	defer func() { err = nil }()
	var (
		userRole *UserRole
		roleIds  []int // 用户在角色ID列表
		funcIds  []int // 用户在功能ID列表
	)
	o := orm.NewOrm()
	now := time.Now()
	userRole, _, err = GetUserRoleByRoleIdAndUserId(int(in.UserId), int(in.RoleId))
	if err != nil {
		err = errors.Wrap(err, "DelUserRole")
		return
	}
	if userRole != nil && userRole.UserRoleId > 0 {
		// 删除redis中的key
		if roleIds, _, err = GetRoleIdsByUserId(int(in.UserId)); err != nil {
			Logger.Error(err.Error())
			return
		}
		for index := 0; roleIds != nil && index < len(roleIds); index++ {
			funcIds, _, err = GetFuncIdsByRoleId(roleIds[index])
			if err != nil {
				Logger.Error(err.Error())
			}
			for subIndex := 0; funcIds != nil && subIndex < len(funcIds); subIndex++ {
				err = conf.Redis__Client.Del(fmt.Sprintf("YCFM_%d_%d", in.UserId, funcIds[subIndex]))
				if err != nil {
					Logger.Error(err.Error())
				}
			}
		}
		userRole.Status = consts.STATUS_DELETED
		userRole.UpdatedAt = now
		if _, err = userRole.UpdateUserRoleNoLock(&o); err != nil {
			err = errors.Wrap(err, "DelUserRole")
			return
		}
	}
	out = in
	return
}

func (s *GrbacServer) GetRoleByRoleCode(in *pb.String, out *pb.Role) (err error) {
	defer func() { err = nil }()
	code := in.Value

	role := Role{}
	if role, _, err = GetRoleByRoleCode(code, nil); err != nil {
		err = errors.Wrap(err, "GetRoleByRoleCode, GrbacServer")
		return
	}

	out = &pb.Role{
		RoleId:   int32(role.RoleId),
		RegionId: int32(role.RegionId),
		Code:     role.Code,
		Name:     role.Name,
		Status:   int32(role.Status),
	}

	return
}
