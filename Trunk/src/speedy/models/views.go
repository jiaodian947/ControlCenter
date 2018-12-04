package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Views struct {
	Id           int       `orm:"column(id);auto" json:"id"`
	Title        string    `orm:"column(title);null" json:"title"`
	Category     int       `orm:"column(category);null" json:"category"`
	CategoryDesc string    `orm:"column(category_desc);null" json:"category_desc"`
	Icon         string    `orm:"column(icon);null" json:"icon"`
	Path         string    `orm:"column(path);null" json:"path"`
	CreateTime   time.Time `orm:"column(create_time);type(datetime);null" json:"create_time"`
	Hidden       int       `orm:"column(hidden)" json:"hidden"`
	ViewId       int       `orm:"column(view_id)" json:"view_id"`
}

func (t *Views) TableName() string {
	return "views"
}

func init() {
	orm.RegisterModel(new(Views))
}

// AddViews insert a new Views into database and returns
// last inserted Id on success.
func AddViews(m *Views) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetViewsById retrieves Views by Id. Returns error if
// Id doesn't exist
func GetViewsById(id int) (v *Views, err error) {
	o := orm.NewOrm()
	v = &Views{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllViews retrieves all Views matches certain condition. Returns empty list if
// no records exist
func GetAllViews(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Views))
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

	var l []Views
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

// UpdateViews updates Views by Id and returns error if
// the record to be updated doesn't exist
func UpdateViewsById(m *Views) (err error) {
	o := orm.NewOrm()
	v := Views{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteViews deletes Views by Id and returns error if
// the record to be deleted doesn't exist
func DeleteViews(id int) (err error) {
	o := orm.NewOrm()
	v := Views{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Views{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
