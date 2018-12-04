package controllers

import "speedy/models"

type Prop struct {
	CheckRight
}

func (p *Prop) GetPropsList() {
	var reply Reply
	reply.Status = 500
	props, err := models.GetAllProps(nil, nil, nil, nil, 0, 0)
	if err != nil {
		reply.Data = err.Error()
		p.Data["json"] = &reply
		p.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"props": props,
	}
	p.Data["json"] = &reply
	p.ServeJSON()
}
