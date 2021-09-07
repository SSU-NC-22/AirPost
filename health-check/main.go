// AirPost
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eunnseo/AirPost/health-check/dataService/memory"
	"github.com/eunnseo/AirPost/health-check/setting"
	"github.com/eunnseo/AirPost/health-check/usecase/healthCheckUC"
	"github.com/eunnseo/AirPost/health-check/usecase/websocketUC"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// TODO: ip 및 port 등등 세팅에 넣어서 처리하기
func main() {
	sr := memory.NewStatusRepo()

	event := make(chan interface{}, 10)
	_ = healthCheckUC.NewHealthCheckUsecase(sr, event)

	wu := websocketUC.NewWebsocketUsecase(event)

	r := gin.New()

	r.GET("/health-check", func(c *gin.Context) {
		log.Println("GET /health-check")
		listen := make(chan interface{})
		wu.Register(listen)
		defer wu.Unregister(listen)

		// conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
		var upgrader = websocket.Upgrader{}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("upgrade: %s", err.Error())
		}
		fmt.Println("connect websocket!")

		for data := range listen {
			log.Printf("read %v\n", data)
			err := conn.WriteJSON(data)
			if err != nil {
				log.Println("conn.WriteJSON error : ", err)
			}
			log.Println("after conn.WriteJSON")
		}
		fmt.Println("disconnect websocket!")
	})

	go log.Fatal(r.Run(setting.Healthsetting.Server))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

}
