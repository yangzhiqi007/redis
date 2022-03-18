package login

import (
	"github.com/davyxu/peckermaster/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(router *gin.RouterGroup) {

	router.POST("/login", func(ctx *gin.Context) {

		type LoginREQ struct {
			Account  string
			Password string
		}

		type LoginACK struct {
			Token   string
			Account string
			Name    string
		}

		var param LoginREQ

		err := ctx.BindJSON(&param)
		if err != nil {
			ctx.String(http.StatusBadRequest, "bind json failed %s", err.Error())
			return
		}

		model.DB.RemoveTimeoutSession()

		if param.Account == "admin" && param.Password == model.AdminPassword {

			ses := model.DB.GenSession("admin")
			model.DB.Save()

			ctx.IndentedJSON(http.StatusOK, &LoginACK{
				Token:   ses.Token,
				Account: ses.Account,
				Name:    ses.Name,
			})

		} else {
			ctx.String(http.StatusBadRequest, "account or password error")
		}
	})
}
