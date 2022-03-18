package procmgr

import (
	"github.com/davyxu/peckermaster/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func initRouter2(router *gin.RouterGroup) {
	router.POST("/server_query", model.HandleAuthSession, func(ctx *gin.Context) {

		nameList := make([]string, 0)

		for _, v := range model.DB.ServerList {
			nameList = append(nameList, v.Name)
		}

		ctx.IndentedJSON(200, nameList)
	})

	router.POST("/server_querydetail", model.HandleAuthSession, func(ctx *gin.Context) {

		ctx.IndentedJSON(200, model.DB.ServerList)
	})

	router.POST("/server_add", model.HandleAuthSession, func(ctx *gin.Context) {

		var param model.Server

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		if param.Name == "" {
			ctx.String(http.StatusBadRequest, "Name empty")
			return
		}

		if model.DB.ServerByName(param.Name) != nil {
			ctx.String(http.StatusBadRequest, "Server name exists")
			return
		}

		model.DB.AddServer(&param)
		model.DB.Save()

		ctx.Status(http.StatusOK)
	})

	router.POST("/server_delete", model.HandleAuthSession, func(ctx *gin.Context) {

		type Request struct {
			Name string
		}

		var param Request

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		task := model.DB.ServerByName(param.Name)
		if task == nil {
			ctx.String(http.StatusBadRequest, "Server not found")
			return
		}

		model.DB.DeleteServer(param.Name)
		model.DB.Save()

	})

	router.POST("/server_update", model.HandleAuthSession, func(ctx *gin.Context) {

		type Request struct {
			Name   string
			Server model.Server
		}

		var param Request

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		if param.Name == "" {
			ctx.String(http.StatusBadRequest, "Name empty")
			return
		}

		if model.DB.ServerByName(param.Name) == nil {
			ctx.String(http.StatusBadRequest, "Server not exists")
			return
		}

		model.DB.UpdateServer(param.Name, &param.Server)
		model.DB.Save()

		ctx.Status(http.StatusOK)
	})
}
