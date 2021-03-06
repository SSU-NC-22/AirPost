package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/domain/model"
	"github.com/eunnseo/AirPost/logic-core/setting"
)

const (
	GoogleSMTPServer = "smtp.gmail.com"
	from = "airpostsch@gmail.com"
	pass = ""
)

type EmailElement struct {
	BaseElement
	Email    string `json:"text"`
	Interval map[string]bool
}

func (ee *EmailElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in EmailElement.Exec !!!!")
	ok, exist := ee.Interval[d.Node.Name]

	if !exist {
		ee.Interval[d.Node.Name] = true
	}
	if ok {
		ee.Interval[d.Node.Name] = false

		to := []string{ee.Email}
		body := fmt.Sprintf("node(%s)", d.Node.Name)
		msg := "From: " + from + "\n" +
			"To: " + strings.Join(to, ",") + "\n" +
			"Subject: AirPost email\n" + body

		err := smtp.SendMail(GoogleSMTPServer + ":587",
			smtp.PlainAuth("", from, pass, GoogleSMTPServer),
			from, to, []byte(msg))

		if err != nil {
			fmt.Printf("smtp error: %s\n", err)
		} else {
			fmt.Println("Mail sent successfully")
		}

		tick := time.NewTicker(10 * time.Second)
		go func() {
			<-tick.C
			ee.Interval[d.Node.Name] = true
		}()
	}
	ee.BaseElement.Exec(d)
}


type ActuatorElement struct {
	BaseElement
	Name   string `json:"name"`
	Values []struct {
		Value int `json:"value"`
		Sleep int `json:"sleep"`
	} `json:"values"`
	Interval map[string]bool
}

type Actuator struct {
	Nid    string `json:"nid"`  // node id
	Name   string `json:"name"` // actuator name
	Values []struct {           // action values
		Value int `json:"value"`
		Sleep int `json:"sleep"`
	} `json:"values"`
}

func (ae *ActuatorElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in ActuatorElement.Exec !!!!")
	ok, exist := ae.Interval[d.Node.Name]
	if !exist {
		ae.Interval[d.Node.Name] = true
	}
	if ok {
		ae.Interval[d.Node.Name] = false
		go func() {
			
			res := Actuator{
				Nid:    "STA" + strconv.Itoa(d.Node.Nid),
				Name:   ae.Name,
				Values: ae.Values,
			}
					
			pbytes, _ := json.Marshal(res)
			buff := bytes.NewBuffer(pbytes)
			addr := (*adapter.AddrMap)[d.Node.Sid] // sink address
			log.Println("in Act.Exec, ?????? ??????: " + "http://" + addr.Addr + "/actuator" + " ????????????: " + string(pbytes))
			resp, err := http.Post("http://"+addr.Addr+"/actuator", "application/json", buff)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
		}()
		
		tick := time.NewTicker(1 * time.Second)
		go func() {
			<-tick.C
			ae.Interval[d.Node.Name] = true
		}()
	}
	ae.BaseElement.Exec(d)
}

type Drone struct {
	Nid    string 	   `json:"nid"` // drone node id
	Values [][]float64 `json:"values"`
	Tagidx int 		   `json:"tagidx"` // values ????????? tag??? ????????? index?????? (0~)
}

type DroneElement struct {
	BaseElement
	Nid      string 	 `json:"nid"`
	Values   [][]float64 `json:"values"`
	Tagidx   int 		 `json:"tagidx"`
	Sent	 bool		 `json:"sent"`
}

func (de *DroneElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in DroneElement.Exec !!!!")
			
	if !de.Sent {
		de.Sent = true
		go func() {
			res := Drone{
				Nid:    "DRO0",
				Values: de.Values,
				Tagidx: 1,
			}
					
			pbytes, _ := json.Marshal(res)
			buff := bytes.NewBuffer(pbytes)
			addr := (*adapter.AddrMap)[d.Node.Sid]
			log.Println("in Drone.Exec, ?????? ??????: " + "http://" + addr.Addr + "/drone" + " ????????????: " + string(pbytes))
			resp, err := http.Post("http://"+addr.Addr+"/drone", "application/json", buff)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
		}()
	}
	de.BaseElement.Exec(d)
}

type AlarmElement struct { // ?????? ????????? ?????? action
	BaseElement
	Email      string `json:"email"`
	OrderNum   string `json:"ordernum"`
	SrcStation string `json:"src_station"`
	DestTag    string `json:"dest_tag"`
	SrcName    string `json:"src_name"`
	DestName   string `json:"dest_name"`
}

func (ae *AlarmElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in AlarmElement.Exec !!!!")

	to := []string{ae.Email}
	subject := "AirPost ?????? ?????? - ????????????(" + ae.OrderNum + ")"
	body := "???????????? : " + ae.OrderNum + "\r\n" +
		"?????? ???????????? : " + ae.SrcStation + "\r\n" +
		"?????? ?????? : " + ae.DestTag + "\r\n" + "\r\n" +
		ae.DestName + "???, " + ae.SrcName + "?????? ???????????? ????????? ?????? ?????????????????????."

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" + body

	err := smtp.SendMail(GoogleSMTPServer + ":587",
		smtp.PlainAuth("", from, pass, GoogleSMTPServer),
		from, to, []byte(msg))

	if err != nil {
		log.Panicln("smtp send error: ", err)
	} else {
		log.Println("smtp send ok")
	}

	ae.BaseElement.Exec(d)
}

type MovingElement struct { // ????????????????????? ?????? ????????? ?????????????????? ?????? action
	BaseElement
	Nid int `json:"nid"`
}

type Moving struct {
	Nid int `json:"nid"` // drone node id
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
		Alt float64 `json:"alt"`
	} `json:"location"`
}

func (me *MovingElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in TrackingElement.Exec !!!!")
			
	go func() {
		res := Moving{
			Nid:      me.Nid,
			Location: struct{Lat float64 "json:\"lat\""; Lon float64 "json:\"lon\""; Alt float64 "json:\"alt\""}{
				Lat: d.Values["lat"],
				Lon: d.Values["long"],
				Alt: d.Values["alt"],
			},
		}
				
		pbytes, _ := json.Marshal(res)
		buff := bytes.NewBuffer(pbytes)
		log.Println("in Tracking.Exec, ?????? ??????: " + "http://" + setting.Appsetting.Server + "/regist/node/update" + " ????????????: " + string(pbytes))
		resp, err := http.Post("http://"+setting.Appsetting.Server+"/regist/node/update", "application/json", buff)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}()

	me.BaseElement.Exec(d)
}
