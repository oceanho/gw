package Api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Dto"
)

func GetUser(c *gw.Context) {
	dto := &Dto.UserDto{}
	if c.Bind(dto) != nil {
		return
	}
	//user := &db
	//db := c.Store().GetDbStoreByName("user-primary")
	//biz.CreateUser(db)
}
func CreateUser(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func ModifyUser(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func DeleteUser(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func QueryUser(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId(), c.Query("uid")),
	})
}
