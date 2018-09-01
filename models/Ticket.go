package models

import "time"

type Ticket struct {
	BaseModel
	Ticket string `json:"ticket"`
	TTL int	`json:"ttl"`
	User_id int `json:"user_id"`
	IsVerify bool `json:"is_verify"`
}


func (t Ticket)CreateTicket(sql string) error {
	if mtdt, err := GetMySQLDB().Prepare(sql); err != nil {
		return err
	}else{
		timeStamp := time.Now().Unix()
		if _, err = mtdt.Exec(t.Ticket, t.TTL, t.User_id, timeStamp, timeStamp); err != nil {
			return  err
		}else {
			return nil
		}
	}
}

func (t Ticket)FindOneByTicket(sql string) (*Ticket, error){

	if mtdt, err := GetMySQLDB().Prepare(sql); err != nil{

		return nil, err
	}else {

		var ticket Ticket
		Row :=  mtdt.QueryRow(t.Ticket)
		// ticket, ttl, created_at, updated_at, is_verify
		if err := Row.Scan(&ticket.Ticket, &ticket.TTL, &ticket.CreatedAt, &ticket.UpdatedAt, &ticket.IsVerify, &ticket.User_id);err !=nil{
			return nil, err
		}else{
			return &ticket, nil
		}
	}
}
// 失效ticket
func (t Ticket)InvalidTicket(sql string)  error{
	_, err := GetMySQLDB().Exec(sql)
	if err != nil {
		return err
	}else{
		return nil
	}
}