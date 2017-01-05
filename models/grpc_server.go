package models

import (
	"fmt"

	. "github.com/1046102779/common/utils"
	. "github.com/1046102779/grbac/logger"
	pb "github.com/1046102779/igrpc"
)

type grbacServer struct{}

func (t *grbacServer) LoadGrbacUserRel(in *pb.GrbacUserRel, out *pb.GrbacUserRel) (err error) {
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
		// 存储到redis中，三元表<用户ID-功能ID，实体ID集合={公司ID}>
		RedisClient.SAdd(fmt.Sprintf("YCFM_%d_%d", in.UserId, funcIds[index]), in.CompanyId)
	}
	return
}
