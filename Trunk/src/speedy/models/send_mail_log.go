package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type SendMailLog struct {
	Id           int       `orm:"column(id);auto" json:"id"`
	SendGameid   int       `orm:"column(send_gameid)" json:"send_gameid"`
	SendServerid int       `orm:"column(send_serverid)" json:"send_serverid"`
	ToRoles      string    `orm:"column(to_roles);null" json:"to_roles"`
	MailType     int       `orm:"column(mail_type)" json:"mail_type"`
	MailSubtype  int       `orm:"column(mail_subtype)" json:"mail_subtype"`
	MailTitle    string    `orm:"column(mail_title)" json:"mail_title"`
	MailAppendix string    `orm:"column(mail_appendix);null" json:"mail_appendix"`
	SendType     int       `orm:"column(send_type)" json:"send_type"`
	SendReason   string    `orm:"column(send_reason)" json:"send_reason"`
	SendTime     time.Time `orm:"column(send_time);type(datetime)" json:"send_time"`
	Sender       string    `orm:"column(sender)" json:"sender"`
}

func (t *SendMailLog) TableName() string {
	return "send_mail_log"
}

func init() {
	orm.RegisterModel(new(SendMailLog))
}
func GetCountSendMailLog() (count int64, err error) {
	o := orm.NewOrm()
	count, err = o.QueryTable(new(SendMailLog)).Count()
	return
}

// AddSendMailLog insert a new SendMailLog into database and returns
// last inserted Id on success.
func AddSendMailLog(m *SendMailLog) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSendMailLogById retrieves SendMailLog by Id. Returns error if
// Id doesn't exist
func GetSendMailLogById(id int) (v *SendMailLog, err error) {
	o := orm.NewOrm()
	v = &SendMailLog{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSendMailLog retrieves all SendMailLog matches certain condition. Returns empty list if
// no records exist
func GetAllSendMailLog(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SendMailLog))
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

	var l []SendMailLog
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

// UpdateSendMailLog updates SendMailLog by Id and returns error if
// the record to be updated doesn't exist
func UpdateSendMailLogById(m *SendMailLog) (err error) {
	o := orm.NewOrm()
	v := SendMailLog{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSendMailLog deletes SendMailLog by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSendMailLog(id int) (err error) {
	o := orm.NewOrm()
	v := SendMailLog{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SendMailLog{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
