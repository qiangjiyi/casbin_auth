package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type Role struct {
	Id          int `orm:"column(id);auto"`
	RoleName    string
	Description string
	CreateTime  time.Time
	UpdateTime  time.Time
}

func (t *Role) TableName() string {
	return "role"
}

func init() {
	orm.RegisterModel(new(Role))
}

// AddRole insert a new Role into database and returns
// last inserted Id on success
func AddRole(m *Role) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRoleById retrieves Role by Id. Returns errors if
// Id doesn't exist
func GetRoleById(id int) (v *Role, err error) {
	o := orm.NewOrm()
	v = &Role{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateRoleById updates Role by Id and returns errors if
// the record to be updated doesn't exist
func UpdateRoleById(m *Role) (err error) {
	o := orm.NewOrm()
	v := Role{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, "role_name", "description", "update_time"); err == nil {
			logs.GetLogger().Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRole deletes Role by Id and returns errors if
// the record to be deleted doesn't exist
func DeleteRole(id int) (err error) {
	o := orm.NewOrm()
	v := Role{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Role{Id: id}); err == nil {
			logs.GetLogger().Println("Number of records deleted in database:", num)
		}
	}
	return
}
