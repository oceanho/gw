package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/tester/biz"
	"github.com/oceanho/gw/web/apps/tester/dto"
)

func CreateMyTester(c *gw.Context) {
	obj := &dto.MyTester{}
	err := c.Bind(obj)
	if err != nil {
		c.Err400Msg(4000, fmt.Sprintf("invalid request params, err: %v", err))
		return
	}
	err = biz.CreateMyTester(c.Store.GetDbStore(), obj)
	if err != nil {
		c.Err500Msg(5000, err)
	} else {
		c.OK(obj.ID)
	}
}

func QueryMyTester(c *gw.Context) {
	query := &gw.QueryExpr{}
	if c.Bind(query) != nil {
		return
	}
	objs := make([]dto.MyTester, 0)
	err := biz.QueryMyTester(c.Store.GetDbStore(), &objs)
	c.JSON(err, objs)
}

func GetTester(c *gw.Context) {
	c.OK(struct {
		RequestID string
	}{
		RequestID: c.RequestID,
	})
}

func GetTester400(c *gw.Context) {
	c.Err400(4000)
}

func GetTester400WithCustomErr(c *gw.Context) {
	c.Err400Msg(4000, "Custom 400 Err")
}

func GetTester400WithCustomPayload(c *gw.Context) {
	c.Err400Payload(4000, "Custom 400 Payload")
}

func GetTester400WithCustomPayloadErr(c *gw.Context) {
	c.Err400PayloadMsg(4000, "Custom 400 Err and Payload", gin.H{"A": "a"})
}

func GetTester401(c *gw.Context) {
	c.Err401(4001)
}

func GetTester401WithCustomErr(c *gw.Context) {
	c.Err401Msg(4001, "Custom 401 Err")
}

func GetTester401WithCustomPayload(c *gw.Context) {
	c.Err401Payload(4001, "Custom 401 Payload")
}

func GetTester401WithCustomPayloadErr(c *gw.Context) {
	c.Err401PayloadMsg(4001, "Custom 401 Err and Payload", gin.H{"A": "a"})
}

func GetTester403(c *gw.Context) {
	c.Err403(4003)
}

func GetTester403WithCustomErr(c *gw.Context) {
	c.Err403Msg(4003, "Custom 403 Err")
}

func GetTester403WithCustomPayload(c *gw.Context) {
	c.Err403Payload(4003, "Custom 403 Payload")
}

func GetTester403WithCustomPayloadErr(c *gw.Context) {
	c.Err403PayloadMsg(4003, "Custom 403 Err and Payload", gin.H{"A": "a"})
}

func GetTester404(c *gw.Context) {
	c.Err404(4004)
}

func GetTester404WithCustomErr(c *gw.Context) {
	c.Err404Msg(4004, "Custom 404 Err")
}

func GetTester404WithCustomPayload(c *gw.Context) {
	c.Err404Payload(4004, "Custom 404 Payload")
}

func GetTester404WithCustomPayloadErr(c *gw.Context) {
	c.Err404PayloadMsg(4004, "Custom 404 Err and Payload", gin.H{"A": "a"})
}

func GetTester500(c *gw.Context) {
	c.Err500(5000)
}

func GetTester500WithCustomErr(c *gw.Context) {
	c.Err500Msg(5000, "Custom 500 Err")
}

func GetTester500WithCustomPayload(c *gw.Context) {
	c.Err500Payload(5000, "Custom 500 Payload")
}

func GetTester500WithCustomPayloadErr(c *gw.Context) {
	c.Err500PayloadMsg(5000, "Custom 500 Err and Payload", gin.H{"A": "a"})
}
