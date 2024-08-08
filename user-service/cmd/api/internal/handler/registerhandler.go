package handler

import (
	"net/http"

	"IM/user-service/cmd/api/internal/logic"
	"IM/user-service/cmd/api/internal/svc"
	"IM/user-service/cmd/api/internal/types"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(r *gin.Engine, ctx *svc.ServiceContext) {
	r.POST("/api/user/register", func(c *gin.Context) {
		var req types.UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l := logic.NewRegisterLogic(c.Request.Context(), ctx)
		resp, err := l.Register(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, resp)
		}
	})

}
