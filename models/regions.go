package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Regions struct {
	Id        int       `orm:"column(region_id);auto"`
	Code      string    `orm:"column(code);size(50);null"`
	Name      string    `orm:"column(name);size(50);null"`
	Status    int16     `orm:"column(status);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
}

func (t *Regions) TableName() string {
	return "regions"
}

func init() {
	orm.RegisterModel(new(Regions))
}
