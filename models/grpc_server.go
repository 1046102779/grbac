// rpcx服务列表
// 1. 加载用户关系数据到redis中
// 2. 新增系统账户，修改系统或者员工角色时，需要修改用户与角色之间的关系
package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	utils "github.com/1046102779/common"
	. "github.com/1046102779/grbac/logger"
	pb "github.com/1046102779/igrpc"
)

type GrbacServer struct{}

func (t *GrbacServer) LoadGrbacUserRel(in *pb.GrbacUserRel, out *pb.GrbacUserRel) (err error) {
	Logger.Info("[%v] enter LoadGrbacUserRel.", in.UserId)
	defer Logger.Info("[%v] left LoadGrbacUserRel.", in.UserId)
	var (
		roleIds []int // 用户在角色ID列表
		funcIds []int // 用户在功能ID列表
	)
	defer func() {
		err = nil
	}()
	if roleIds, _, err = GetRoleIdsByUserId(int(in.UserId)); err != nil {
		Logger.Error(err.Error())
		return
	}
	for index := 0; roleIds != nil && index < len(roleIds); index++ {
		funcIds, _, err = GetFuncIdsByRoleId(roleIds[index])
		// 添加key，三元表<用户ID-功能ID，实体ID集合={公司ID}>
		if err != nil {
			Logger.Error(err.Error())
		}
		for subIndex := 0; funcIds != nil && subIndex < len(funcIds); subIndex++ {
			utils.RedisClient.SAdd(fmt.Sprintf("YCFM_%d_%d", in.UserId, funcIds[index]), in.CompanyId)
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
		userRole := &UserRoles{
			UserId:    int(in.UserId),
			RoleId:    roleId,
			Status:    utils.STATUS_VALID,
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
