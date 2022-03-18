package procmgr

import (
	"github.com/davyxu/pecker/client"
	"github.com/davyxu/peckermaster/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func initRouter(router *gin.RouterGroup) {

	router.POST("/task_query", model.HandleAuthSession, func(ctx *gin.Context) {

		ctx.IndentedJSON(http.StatusOK, model.DB.TaskList)
	})

	router.POST("/task_create", model.HandleAuthSession, func(ctx *gin.Context) {

		var param model.Task

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		if param.Name == "" {
			ctx.String(http.StatusBadRequest, "Name empty")
			return
		}

		if model.DB.TaskByName(param.Name) != nil {
			ctx.String(http.StatusBadRequest, "Task name exists")
			return
		}

		addr := model.DB.GetAddress(param.Server)
		if addr == "" {
			ctx.String(http.StatusBadRequest, "invalid server")
			return
		}

		model.DB.AddTask(&param)
		model.DB.Save()

		ctx.Status(http.StatusOK)
	})

	router.POST("/task_delete", model.HandleAuthSession, func(ctx *gin.Context) {

		type Request struct {
			Name string
		}

		var param Request

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		task := model.DB.TaskByName(param.Name)
		if task == nil {
			ctx.String(http.StatusBadRequest, "Task not found")
			return
		}

		model.DB.DeleteTask(param.Name)
		model.DB.Save()

	})

	router.POST("/task_update", model.HandleAuthSession, func(ctx *gin.Context) {

		type Request struct {
			Name string
			Task model.Task
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

		if model.DB.TaskByName(param.Name) == nil {
			ctx.String(http.StatusBadRequest, "Task not exists")
			return
		}

		addr := model.DB.GetAddress(param.Task.Server)
		if addr == "" {
			ctx.String(http.StatusBadRequest, "invalid server")
			return
		}

		model.DB.UpdateTask(param.Name, &param.Task)
		model.DB.Save()

		ctx.Status(http.StatusOK)
	})

	router.POST("/task_exec", model.HandleAuthSession, func(ctx *gin.Context) {

		type Request struct {
			Name string
		}

		var param Request

		if err := ctx.BindJSON(&param); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		task := model.DB.TaskByName(param.Name)
		if task == nil {
			ctx.String(http.StatusBadRequest, "Task not exists")
			return
		}

		addr := model.DB.GetAddress(task.Server)
		if addr == "" {
			ctx.String(http.StatusBadRequest, "invalid server")
			return
		}

		respond, err := client.ExecuteRemoteCommandText(addr, "script", strings.NewReader(task.Code), "file")
		if err != nil {
			ctx.String(http.StatusOK, "exec shell failed: %s\n%s", err.Error(), respond)
			return
		}

		ctx.IndentedJSON(http.StatusOK, respond)
	})
}
