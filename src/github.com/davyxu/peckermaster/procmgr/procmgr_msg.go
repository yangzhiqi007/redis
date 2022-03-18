package procmgr

import (
	"fmt"
	"github.com/davyxu/pecker/client"
	"github.com/davyxu/peckermaster/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(router *gin.RouterGroup) {
	initRouter(router)
	initRouter2(router)

	router.POST("/proc_query", model.HandleAuthSession, func(ctx *gin.Context) {

		type Request struct {
			Server string
		}

		var param Request

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		addr := model.DB.GetAddress(param.Server)
		if addr == "" {
			ctx.String(http.StatusBadRequest, "invalid server name")
			return
		}

		text, err := client.ExecuteRemoteCommand(addr, `supervisorctl status`)
		if err != nil {
			ctx.String(http.StatusOK, "exec supervisor command failed: %s", err.Error())
			return
		}

		status := model.ParseStatusText(text)

		ctx.IndentedJSON(200, &status)
	})

	router.POST("/proc_ctl", model.HandleAuthSession, func(ctx *gin.Context) {

		type Request struct {
			Server   string
			NameList []string
			Command  string
		}

		var param Request

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		addr := model.DB.GetAddress(param.Server)
		if addr == "" {
			ctx.String(http.StatusBadRequest, "invalid pecker name")
			return
		}

		for _, name := range param.NameList {
			var supervisorCommand string
			switch param.Command {
			case "start":
				supervisorCommand = fmt.Sprintf("supervisorctl start %s", name)
			case "stop":
				supervisorCommand = fmt.Sprintf("supervisorctl stop %s", name)
			case "restart":
				supervisorCommand = fmt.Sprintf("supervisorctl restart %s", name)
			default:
				ctx.String(http.StatusBadRequest, "invalid command")
				return
			}

			_, err := client.ExecuteRemoteCommand(addr, supervisorCommand)
			if err != nil {
				ctx.String(http.StatusOK, "exec supervisor command failed: %s", err.Error())
				return
			}
		}

		ctx.Status(http.StatusOK)
	})

	router.POST("/proc_log", model.HandleAuthSession, func(ctx *gin.Context) {

		type Request struct {
			Server string
			Name   string
		}

		var param Request

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		addr := model.DB.GetAddress(param.Server)
		if addr == "" {
			ctx.String(http.StatusBadRequest, "invalid pecker name")
			return
		}

		supervisorCommand := fmt.Sprintf("supervisorctl tail %s", param.Name)
		text, err := client.ExecuteRemoteCommand(addr, supervisorCommand)
		if err != nil {
			ctx.String(http.StatusOK, "exec supervisor command failed: %s", err.Error())
			return
		}

		ctx.IndentedJSON(http.StatusOK, text)
	})

}
