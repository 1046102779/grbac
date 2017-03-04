package models

import (
	"time"

	utils "github.com/1046102779/common"
	. "github.com/1046102779/grbac/logger"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
)

type Roles struct {
	Id        int       `orm:"column(role_id);auto"`
	RegionId  int       `orm:"column(region_id);null"`
	Name      string    `orm:"column(name);size(50);null"`
	Status    int16     `orm:"column(status);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
}

func (t *Roles) TableName() string {
	return "roles"
}

func init() {
	orm.RegisterModel(new(Roles))
}

func getRoleIdByCode(code string) (id int, retcode int, err error) {
	Logger.Info("[%v] enter getRoleIdByCode.", code)
	defer Logger.Info("[%v] left getRoleIdByCode.", code)
	var role *Roles = new(Roles)
	o := orm.NewOrm()
	err = o.QueryTable((&Roles{}).TableName()).Filter("code__icontains", code).One(&role)
	if err != nil {
		err = errors.Wrap(err, "getRoleIdByCode")
		retcode = utils.DB_READ_ERROR
		return
	}
	id = role.Id
	return
}
