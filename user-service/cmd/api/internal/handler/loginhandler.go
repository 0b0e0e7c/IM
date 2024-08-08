package handler

import (
	"net/http"

	"IM/user-service/cmd/api/internal/logic"
	"IM/user-service/cmd/api/internal/svc"
	"IM/user-service/cmd/api/internal/types"

	"github.com/gin-gonic/gin"
)

func LoginHandlers(r *gin.Engine, ctx *svc.ServiceContext) {

	r.POST("/api/user/login", func(c *gin.Context) {
		var req types.UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l := logic.NewLoginLogic(c.Request.Context(), ctx)
		resp, err := l.Login(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, resp)
		}
	})
}
