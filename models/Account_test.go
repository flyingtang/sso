package models

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
)

func TestGeneraterTicket(t *testing.T){
	Convey("check ticket ", t, func() {
		rd := GeneraterTicket("admin")
			//So(rd, is)
			fmt.Println("rd := ", rd)
	})
}
