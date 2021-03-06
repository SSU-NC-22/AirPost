// AirPost
package main

import (
	"log"

	"github.com/eunnseo/AirPost/application/dataService/sql"
	"github.com/eunnseo/AirPost/application/docs"
	"github.com/eunnseo/AirPost/application/domain/model"
	"github.com/eunnseo/AirPost/application/domain/repository"
	"github.com/eunnseo/AirPost/application/rest/handler"
	"github.com/eunnseo/AirPost/application/setting"
	"github.com/eunnseo/AirPost/application/usecase"
	"github.com/eunnseo/AirPost/application/usecase/eventUsecase"
	"github.com/eunnseo/AirPost/application/usecase/registUsecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

func main() {
	sql.Setup()

	sir := sql.NewSinkRepo()
	ndr := sql.NewNodeRepo()
	lgr := sql.NewLogicRepo()
	lsr := sql.NewLogicServiceRepo()
	tpr := sql.NewTopicRepo()

	dlr := sql.NewDeliveryRepo()
	ptr := sql.NewPathRepo()
	sdr := sql.NewStationDroneRepo()

	ru := registUsecase.NewRegistUsecase(sir, ndr, lgr, lsr, tpr, dlr, ptr, sdr)
	eu := eventUsecase.NewEventUsecase(sir, lsr)

	h := handler.NewHandler(ru, eu)

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// swagger
	docs.SwaggerInfo.Title = "AirPost application API"
	docs.SwaggerInfo.Description = "This is a registration server for AirPost UI."
	docs.SwaggerInfo.Version = "0.1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	setRegistrationRoute(r, h)
	setEventRoute(r, h)
	initTopic(tpr)

	initDroneSink(sir, eu)
	initStationSink(sir, eu)
	initTagSink(sir, eu)

	log.Fatal(r.Run(setting.Appsetting.Server))
}

func setEventRoute(r *gin.Engine, h *handler.Handler) {
	event := r.Group("/event")
	{
		event.POST("", h.RegistLogicService)
	}
}

func setRegistrationRoute(r *gin.Engine, h *handler.Handler) {
	regist := r.Group("/regist")
	{

		sink := regist.Group("/sink")
		{
			sink.GET("", h.ListSinks)
			sink.POST("", h.RegistSink)
			sink.DELETE("/:id", h.UnregistSink)
		}
		node := regist.Group("/node")
		{
			node.GET("", h.ListNodes)
			node.GET("/:sinkid", h.ListNodesBySink)
			node.POST("", h.RegistNode)
			node.POST("/update", h.UpdateNodeLoc)
			node.DELETE("/:id", h.UnregistNode)
		}
		logic := regist.Group("/logic")
		{
			logic.GET("", h.ListLogics)
			logic.POST("", h.RegistLogic) // << ???????????????
			logic.DELETE("/:id", h.UnregistLogic)
		}
		logicService := regist.Group("/logic-service")
		{
			logicService.GET("", h.ListLogicServices)
			logicService.DELETE("/:id", h.UnregistLogicService)
		}
		topic := regist.Group("/topic")
		{
			topic.GET("", h.ListTopics)
			topic.POST("", h.RegistTopic)
			topic.DELETE("/:id", h.UnregistTopic)
		}
		delivery := regist.Group("/delivery")
		{
			delivery.GET("/:orderNum", h.GetDroneID)
			delivery.POST("", h.RegistDelivery)
		}
		tracking := regist.Group("/tracking")
		{
			tracking.GET("/:orderNum", h.GetTracking)
		}
	}  
}

func initTopic(tpr repository.TopicRepo) {
	if setting.Topicsetting.Name != "" {
		t := model.Topic{
			Name:         setting.Topicsetting.Name,
			Partitions:   setting.Topicsetting.Partitions,
			Replications: setting.Topicsetting.Replications,
		}
		tpr.Create(&t)
	}
}

func initDroneSink(sir repository.SinkRepo, eu usecase.EventUsecase) {
	if setting.DroneSinksetting.Name != "" {
		s := model.Sink{
			Name:		setting.DroneSinksetting.Name,
			Addr:		setting.DroneSinksetting.Addr,
			TopicID:	setting.DroneSinksetting.TopicID,
		}
		sir.Create(&s)
		eu.CreateSinkEvent(&s)
	}
}

func initStationSink(sir repository.SinkRepo, eu usecase.EventUsecase) {
	if setting.StationSinksetting.Name != "" {
		s := model.Sink{
			Name:		setting.StationSinksetting.Name,
			Addr:		setting.StationSinksetting.Addr,
			TopicID:	setting.StationSinksetting.TopicID,
		}
		sir.Create(&s)
		eu.CreateSinkEvent(&s)
	}
}

func initTagSink(sir repository.SinkRepo, eu usecase.EventUsecase) {
	if setting.TagSinksetting.Name != "" {
		s := model.Sink{
			Name:		setting.TagSinksetting.Name,
			Addr:		setting.TagSinksetting.Addr,
			TopicID:	setting.TagSinksetting.TopicID,
		}
		sir.Create(&s)
		eu.CreateSinkEvent(&s)
	}
}

