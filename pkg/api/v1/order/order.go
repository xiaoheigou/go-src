package order

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)


func GetOrder(c *gin.Context) {

	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	session.Options(sessions.Options{MaxAge:3600})
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count++
	}
	session.Set("count", count)
	session.Save()
	m := make(map[string]interface{})
	json.Unmarshal([]byte(""),&m)
	c.JSON(200, gin.H{"count": count})
	//order := models.Order{
	//	BuyerId: 12,
	//}
	//
	//
	//utils.DB.Create(&order)
	//c.JSON(200, order)
}
