package handler

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/eunnseo/AirPost/application/adapter"
	"github.com/eunnseo/AirPost/application/domain/model"
	"github.com/gin-gonic/gin"
)

const (
	DRONE   = 1 // drone sink id
	STATION = 2 // station sink id
	TAG 	= 3 // tag sink id
)

/**************************************************************/
/* Sink handler                                               */
/**************************************************************/
// ListSinks ...
// @Summary List sink node(raspi info)
// @Description get sinks list
// @Tags sink
// @Param  page query int false "page num"
// @Param  size query int false "page size(row)"
// @Produce  json
// @Success 200 {array} model.Sink "default, return all sinks."
// @Success 201 {object} adapter.SinkPage "if page query is exist, return pagenation result. pages only valid when page is 1."
// @Router /regist/sink [get]
func (h *Handler) ListSinks(c *gin.Context) {
	var (
		err   error
		sinks []model.Sink
		page  adapter.Page
		pages int
	)
	if c.Bind(&page); page.IsBinded() {
		if page.Size == 0 {
			page.Size = 10
		}
		if sinks, err = h.ru.GetSinksPage(page); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if page.Page == 1 {
			pages = h.ru.GetSinkPageCount(page.Size)
		}

		for i, sink := range sinks {
			nodes, err := h.ru.GetNodesBySinkID(sink.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			sinks[i].Nodes = append(sink.Nodes, nodes...)
		}
		c.JSON(http.StatusOK, gin.H{"sinks": sinks, "pages": pages})
		return
	} else {
		sinks, err := h.ru.GetSinks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for i, sink := range sinks {
			nodes, err := h.ru.GetNodesBySinkID(sink.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			sinks[i].Nodes = append(sink.Nodes, nodes...)
		}
		c.JSON(http.StatusOK, sinks)
		return
	}
}

// RegistSink ...
// @Summary Add sink node(raspi info)
// @Description Add sink node
// @Tags sink
// @Param  sink body model.Sink true "name, address(only ip address, don't include port)"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Sink "include topic info"
// @Router /regist/sink [post]
func (h *Handler) RegistSink(c *gin.Context) {
	var sink model.Sink

	if err := c.ShouldBindJSON(&sink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ru.RegistSink(&sink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.eu.CreateSinkEvent(&sink)
	c.JSON(http.StatusOK, sink)
}

// UnregistSink ...
// @Summary Delete sink node(raspi info)
// @Description Delete sink node
// @Tags sink
// @Param  id path int true "sink's id"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Sink "include topic, nodes info"
// @Router /regist/sink [delete]
func (h *Handler) UnregistSink(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sink := model.Sink{ID: id}
	err = h.ru.UnregistSink(&sink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.eu.DeleteSinkEvent(&sink)
	c.JSON(http.StatusOK, sink)
}

/**************************************************************/
/* Node handler                                               */
/**************************************************************/
// ListNodes ...
// @Summary List sensor node
// @Description get nodes listh.eu.CreateNodeEvent(&node)
// @Tags node
// @Param  page query int false "page num"
// @Param  size query int false "page size(row)"
// @Param  sink query int false "sink filter"
// @Param  left query float32 false "location(longitude) filter"
// @Param  right query float32 false "location(longitude) filter"
// @Param  up query float32 false "location(Latitude) filter"
// @Param  down query float32 false "location(Latitude) filter"
// @Produce  json
// @Success 200 {array} model.Node "default, return all nodes. if location query is exist, return location filter result(square)."
// @Success 201 {object} adapter.NodePage "if page query is exist, return pagenation result. pages only valid when page is 1."
// @Router /regist/node [get]
func (h *Handler) ListNodes(c *gin.Context) {
	var (
		err    error
		nodes  []model.Node
		page   adapter.Page
		pages  int
		square adapter.Square
	)

	if c.Bind(&page); page.IsBinded() {
		if page.Size == 0 {
			page.Size = 10
		}
		if nodes, err = h.ru.GetNodesPage(page); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if page.Page == 1 {
			pages = h.ru.GetNodePageCount(page)
		}
		c.JSON(http.StatusOK, gin.H{"nodes": nodes, "pages": pages})
		return
	} else if c.Bind((&square)); square.IsBinded() {
		if nodes, err = h.ru.GetNodesSquare(square); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, nodes)
		return
	} else {
		nodes, err := h.ru.GetNodes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, nodes)
		return
	}

}

// ListNodesBySink ...
func (h *Handler) ListNodesBySink(c *gin.Context) {
	sinkid, err := strconv.Atoi(c.Param("sinkid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nodes, err := h.ru.GetNodesBySinkID(sinkid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

// RegistNode ...
// @Summary Add sensor node
// @Description Add sensor node
// @Tags node
// @Param  node body model.Node true "name, lat, lng, sink_id"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Node "include sink, sink.topic, sensors, sensors.logics info"
// @Router /regist/node [post]
func (h *Handler) RegistNode(c *gin.Context) {
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if node.Type[:3] == "DRO" {
		stationid, err := strconv.Atoi(node.Type[4:])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		node.Type = node.Type[:3]

		err = h.ru.RegistNode(&node)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// drone ?????? ??? station_drone ??????
		sd := model.StationDrone{
			StationID: stationid,
			DroneID:   node.ID,
			Usable:    true,
		}

		err = h.ru.RegistStationDrone(&sd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// drone ?????????????????? ?????? ?????????????????? logic ?????? ??? ??????
		movingElem := adapter.Element{
			Elem: "moving",
			Arg: map[string]interface{} {
				"nid": node.ID,
			},
		}
		aMovingLogic := adapter.Logic{
			LogicName: "moving-" + strconv.Itoa(node.ID),
			Elems: []adapter.Element{movingElem},
			NodeID: node.ID,
			Node: node,
		}

		movingLogic, err := adapter.LogicToModel(&aMovingLogic)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = h.ru.RegistLogic(&movingLogic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.eu.CreateLogicEvent(&movingLogic)

	} else if node.Type[:3] == "STA" {
		err := h.ru.RegistNode(&node)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// regist night LED logic
		night := adapter.Range{Max: 512, Min: 0}
		rl := []adapter.Range{night}
		e1 := adapter.Element{
			Elem: "value",
			Arg:  map[string]interface{} {
				"value": "light",
				"range": rl,
			},
		}
		v := adapter.Value{Value: 1, Sleep: 0}
		vl := []adapter.Value{v}
		e2 := adapter.Element{
			Elem: "actuator",
			Arg:  map[string]interface{} {
				"nid": "STA" + strconv.Itoa(node.ID),
				"name": "LED ON",
				"values": vl,
			},
		}
		aLedLogic := adapter.Logic{
			LogicName: "LED ON station#" + strconv.Itoa(node.ID),
			Elems: []adapter.Element{e1, e2},
			NodeID: node.ID,
			Node: node,
		}

		nightLedLogic, err := adapter.LogicToModel(&aLedLogic)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = h.ru.RegistLogic(&nightLedLogic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.eu.CreateLogicEvent(&nightLedLogic)

		// regist day LED logic
		day := adapter.Range{Max: 1024, Min: 512}
		rl = []adapter.Range{day}
		e1 = adapter.Element{
			Elem: "value",
			Arg:  map[string]interface{} {
				"value": "light",
				"range": rl,
			},
		}
		v = adapter.Value{Value: 0, Sleep: 0}
		vl = []adapter.Value{v}
		e2 = adapter.Element{
			Elem: "actuator",
			Arg:  map[string]interface{} {
				"nid": "STA" + strconv.Itoa(node.ID),
				"name": "LED OFF",
				"values": vl,
			},
		}
		aLedLogic = adapter.Logic{
			LogicName: "LED OFF station#" + strconv.Itoa(node.ID),
			Elems: []adapter.Element{e1, e2},
			NodeID: node.ID,
			Node: node,
		}

		dayLedLogic, err := adapter.LogicToModel(&aLedLogic)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = h.ru.RegistLogic(&dayLedLogic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.eu.CreateLogicEvent(&dayLedLogic)

		// regist path between station and tag
		tags, err := h.ru.GetNodesBySinkID(TAG)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, tag := range tags {
			/* 
			// TODO: using naver map api
			start := node.LocLon + "," + node.LocLat
			goal := tag.LocLon + "," + tag.LocLat
			client := resty.New()
			resp, err := client.R().
				SetQueryString("start=" + start + "&goal=" + goal).
				SetHeader("X-NCP-APIGW-API-KEY-ID", "").
				SetHeader("X-NCP-APIGW-API-KEY", "").
				Get("https://naveropenapi.apigw.ntruss.com/map-direction/v1/driving")
			if err != nil {
				panic(err)
			}
			*/

			// calc distance
			powLon := math.Pow((node.LocLon - tag.LocLon), 2)
			powLat := math.Pow((node.LocLat - tag.LocLat), 2)
			dist := math.Pow((powLon + powLat), 0.5)
			
			path := model.Path{
				StationID: node.ID,
				TagID: tag.ID,
				Path: "",
				Distance: dist,
			}
			if err := h.ru.RegistPath(&path); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		// resp.Body.Close()		
	} else if node.Type[:3] == "TAG" {
		err := h.ru.RegistNode(&node)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// regist path between station and tag
		stations, err := h.ru.GetNodesBySinkID(STATION)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, station := range stations {
			// calc distance
			powLon := math.Pow((node.LocLon - station.LocLon), 2)
			powLat := math.Pow((node.LocLat - station.LocLat), 2)
			dist := math.Pow((powLon + powLat), 0.5)
			
			path := model.Path{
				StationID: station.ID,
				TagID: node.ID,
				Path: "",
				Distance: dist,
			}
			if err := h.ru.RegistPath(&path); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	h.eu.CreateNodeEvent(&node)
	go h.eu.PostToSink(node.SinkID)
	c.JSON(http.StatusOK, node)
}

// UnregistNode ...
// @Summary Delete sensor node
// @Description Delete sensor node
// @Tags node
// @Param  id path int true "node's id"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Node "include sink, sink.topic info"
// @Router /regist/node [delete]
func (h *Handler) UnregistNode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if drone, ?????? ????????? station_drone ????????? ?????? ??????
	// if station, station_drone ????????? ????????? ?????? ?????? ?????????
	node, _ := h.ru.GetNodeByID(id)
	if node.Type == "DRO" {
		sd := model.StationDrone{DroneID: id}
		h.ru.UnregistStationDroneByDroneID(&sd)
	} else if node.Type == "STA" {
		sdl, _ := h.ru.GetStationDroneByStationID(id)
		if len(sdl) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	
	err = h.ru.UnregistNode(node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.eu.DeleteNodeEvent(node)
	go h.eu.PostToSink(node.SinkID)
	c.JSON(http.StatusOK, node)
}

func (h *Handler) UpdateNodeLoc(c *gin.Context) {
	var new struct{
		Nid      int              `json:"nid"`
		Location adapter.Location `json:"location"`
	}
	if err := c.ShouldBindJSON(&new); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node, err := h.ru.GetNodeByID(new.Nid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.ru.UpdateNodeLoc(node, &new.Location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, new)
}

/**************************************************************/
/* Logic handler                                              */
/**************************************************************/
// ListLogics ...
// @Summary List logics info
// @Description get logics list
// @Tags logic
// @Produce  json
// @Success 200 {array} model.Logic "return all logics info."
// @Router /regist/logic [get]
func (h *Handler) ListLogics(c *gin.Context) {
	logics, err := h.ru.GetLogics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	aLogics := adapter.LogicsToAdapter(logics)
	c.JSON(http.StatusOK, aLogics)
}

// RegistLogic ...
// @Summary Add logic info
// @Description Add logic info
// @Tags logic
// @Param  logic body adapter.Logic true "logic_name, elems"
// @Accept  json
// @Produce  json
// @Success 200 {object} adapter.Logic "include sensor info"
// @Router /regist/logic [post]
func (h *Handler) RegistLogic(c *gin.Context) {
	var aLogic adapter.Logic
	if err := c.ShouldBindJSON(&aLogic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logic, err := adapter.LogicToModel(&aLogic)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.ru.RegistLogic(&logic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resLogic, _ := adapter.LogicToAdapter(&logic)
	h.eu.CreateLogicEvent(&logic)
	c.JSON(http.StatusOK, resLogic)
}

// UnregistLogic ...
// @Summary Delete logic
// @Description Delete logic
// @Tags logic
// @Param  id path int true "logic's id"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Logic "include sensor info"
// @Router /regist/logic [delete]
func (h *Handler) UnregistLogic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logic := model.Logic{ID: id}

	err = h.ru.UnregistLogic(&logic)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resLogic, _ := adapter.LogicToAdapter(&logic)
	h.eu.DeleteLogicEvent(&logic)
	c.JSON(http.StatusOK, resLogic)
}

/**************************************************************/
/* Logic service handler                                      */
/**************************************************************/
// ListLogicServices ...
// @Summary List LogicServices info
// @Description get LogicServices list
// @Tags LogicService
// @Produce  json
// @Success 200 {array} model.LogicService "return all logics info."
// @Router /regist/logic-service [get]
func (h *Handler) ListLogicServices(c *gin.Context) {
	logicServices, err := h.ru.GetLogicServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logicServices)
}

// UnregistLogicService ...
// @Summary Delete LogicService
// @Description Delete LogicService
// @Tags logicService
// @Param  id path int true "logicSerivce's id"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Logic "include topic info"
// @Router /regist/logic-service [delete]
func (h *Handler) UnregistLogicService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logicService := model.LogicService{ID: id}

	err = h.ru.UnregistLogicService(&logicService)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logicService)
}

/**************************************************************/
/* Topic handler                                              */
/**************************************************************/
// ListTopics ...
// @Summary List topics info
// @Description get topics list
// @Tags topic
// @Produce  json
// @Success 200 {array} model.Topic "return all topics info."
// @Router /regist/topic [get]
func (h *Handler) ListTopics(c *gin.Context) {
	topics, err := h.ru.GetTopics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, topics)
}

// RegistTopic ...
// @Summary Add topic info
// @Description Add topic info
// @Tags topic
// @Param  logic body model.Logic true "name, partitions, replications"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Topic
// @Router /regist/topic [post]
func (h *Handler) RegistTopic(c *gin.Context) {
	var topic model.Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ru.RegistTopic(&topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, topic)
}

// UnregistTopic ...
// @Summary Delete topic(kafka topic for logicservices)
// @Description Delete topic(kafka topic for logicservices)
// @Tags topic
// @Param  id path int true "topic's id"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Topic "include logicService info"
// @Router /regist/topic [delete]
func (h *Handler) UnregistTopic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topic := model.Topic{ID: id}

	err = h.ru.UnregistTopic(&topic)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, topic)
}

/**************************************************************/
/* Delivery service handler                                   */
/**************************************************************/
// RegistDelivery ...
// @Summary Add delivery info
// @Description Add delivery info
// @Tags delivery
// @Param  
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Delivery
// @Router /regist/delivery [post]
func (h *Handler) RegistDelivery(c *gin.Context) {
	var delivery model.Delivery
	delivery.CreatedAt = time.Now()

	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SrcStation??? ????????? ??? ???????????? ????????? ????????? ??????
	sdl, err := h.ru.GetStationDroneByStationID(delivery.SrcStationID)
	if err != nil || len(sdl) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	droneid := -1
	for _, sd := range(sdl) {
		if sd.Usable {
			droneid = sd.DroneID
			break
		}
	}
	if droneid == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Regist Delivery with DroneID and Drone
	drone, err := h.ru.GetNodeByID(droneid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	delivery.Drone = *drone
	delivery.DroneID = droneid

	err = h.ru.RegistDelivery(&delivery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// destTag??? ?????? ????????? destStation??? ??????
	destStation, err := h.ru.GetShortestPathStation(delivery.DestTagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// SrcStation??? ????????? drone unregist, DestStation??? ????????? drone regist
	if delivery.SrcStationID != destStation.ID {
		sd := model.StationDrone{
			StationID: delivery.SrcStationID,
			DroneID:   droneid,
		}
		if err := h.ru.UnregistStationDrone(&sd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sd = model.StationDrone{
			StationID: destStation.ID,
			DroneID:   droneid,
			Usable:    true,
		}
		if err := h.ru.RegistStationDrone(&sd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// srcStation, destTag ????????? ?????????
	srcStation, err := h.ru.GetNodeByID(delivery.SrcStationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	destTag, err := h.ru.GetNodeByID(delivery.DestTagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/* ???????????? path ???????????? logic ?????? ??? ?????? */
	srcStationLoc := []float64{srcStation.LocLat, srcStation.LocLon, 1, 22} // lat, lon, alt, ??????(22: ??????)
	tagLoc := []float64{destTag.LocLat, destTag.LocLon, 1, 16} // lat, lon, alt, ??????(16: ??????)
	destStationLoc := []float64{destStation.LocLat, destStation.LocLon, 0, 21} // lat, lon, alt, ??????(21: ??????)

	path := [][]float64{}
	path = append(path, srcStationLoc)
	path = append(path, tagLoc)
	path = append(path, destStationLoc)

	droneElem := adapter.Element{
		Elem: "drone",
		Arg:  map[string]interface{} {
			"nid":    "DRO" + strconv.Itoa(delivery.DroneID),
			"values": path,
			"tagidx": 1, // TODO
		},
	}
	aDroneLogic := adapter.Logic{
		LogicName: "drone-" + delivery.OrderNum,
		Elems: []adapter.Element{droneElem},
		NodeID: delivery.DroneID,
		Node: delivery.Drone,
	}

	droneLogic, err := adapter.LogicToModel(&aDroneLogic)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// err = h.ru.RegistLogic(&droneLogic)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	h.eu.CreateLogicEvent(&droneLogic)

	/* ?????? ????????? ?????? logic ?????? ??? ?????? */
	e1 := adapter.Element{
		Elem: "arrival",
		Arg: map[string]interface{}{
			"done": 1,
		},
	}
	e2 := adapter.Element{
		Elem: "alarm",
		Arg: map[string]interface{}{
			"email":       delivery.Email,
			"ordernum":    delivery.OrderNum,
			"src_station": srcStation.Name,
			"dest_tag":    destTag.Name,
			"src_name":    delivery.SrcName,
			"dest_name":   delivery.DestName,
		},
	}
	aAlarmLogic := adapter.Logic{
		LogicName: "alarm-" + delivery.OrderNum,
		Elems: []adapter.Element{e1, e2},
		NodeID: delivery.DroneID,
	}

	alarmLogic, err := adapter.LogicToModel(&aAlarmLogic)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// err = h.ru.RegistLogic(&alarmLogic)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	h.eu.CreateLogicEvent(&alarmLogic)

	// go h.eu.CreateDeliveryEvent(&delivery)
	c.JSON(http.StatusOK, delivery)
}

func (h *Handler) GetDroneID(c *gin.Context) {
	ordernum, err := strconv.Atoi(c.Param("orderNum"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery, err := h.ru.GetDeliveryByOrderNum(ordernum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"drone_id": delivery.DroneID})
}

/**************************************************************/
/* Tracking service handler                                   */
/**************************************************************/
func (h *Handler) GetTracking(c *gin.Context) {
	ordernum, err := strconv.Atoi(c.Param("orderNum"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery, err := h.ru.GetDeliveryByOrderNum(ordernum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	src, err := h.ru.GetNodeByID(delivery.SrcStationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dest, err := h.ru.GetNodeByID(delivery.DestTagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	drone, err := h.ru.GetNodeByID(delivery.DroneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tracking := model.Tracking{
		DroneNid: delivery.DroneID,
		SrcLat:   src.LocLat,
		SrcLng:	  src.LocLon,
		DestLat:  dest.LocLat,
		DestLng:  dest.LocLon,
		DroneLat: drone.LocLat,
		DroneLng: drone.LocLon,
	}
	c.JSON(http.StatusOK, tracking)
}
