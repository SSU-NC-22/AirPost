package handler

import (
	"fmt"
	"net/http"

	"github.com/eunnseo/AirPost/application/domain/model"
	"github.com/gin-gonic/gin"
)

/**************************************************************/
/* Logic service handler                                      */
/**************************************************************/

func (h *Handler) RegistLogicService(c *gin.Context) {
	var l model.LogicService
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if l.Topic.Name == "" || l.Addr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("topic name, addr invalid").Error()})
		return
	}

	if err := h.eu.RegistLogicService(&l); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, l.Topic.Sinks)
}
