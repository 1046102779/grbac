package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Entities struct {
	Id        int       `orm:"column(entity_id);auto"`
	RegionId  int       `orm:"column(region_id);null"`
	Name      string    `orm:"column(name);size(300);null"`
	ThirdId   int       `orm:"column(third_id);null"`
	Status    int16     `orm:"column(status);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
}

func (t *Entities) TableName() string {
	return "entities"
}

func init() {
	orm.RegisterModel(new(Entities))
}
