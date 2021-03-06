package logic

import (
	"log"
	"time"

	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

type ValueElement struct {
	BaseElement
	Value string `json:"value"`
	Range []struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"range"`
}

func (ve *ValueElement) Exec(d *model.LogicData) {
	
	log.Println("!!!!in ValueElement.Exec !!!!")
	v, ok := d.Values[ve.Value]
	if !ok {
		return
	}
	isRange := false
	
	for _, rg := range ve.Range {
		if rg.Min <= v && v < rg.Max {
			isRange = true
		}
	}
	if isRange {
		
		ve.BaseElement.Exec(d)
	}
}

type TimeElement struct {
	BaseElement
	Range []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"range"`
}

func (te *TimeElement) Exec(d *model.LogicData) {
	
	log.Println("!!!!in TimeElement.Exec !!!!")
	timeFmt := "15:04:05"
	isRange := false
	
	t, _ := time.Parse(timeFmt, d.Timestamp.Format(timeFmt))
	for _, rg := range te.Range {
		st, _ := time.Parse("15:04:05", rg.Start)
		et, _ := time.Parse("15:04:05", rg.End)
		if st.After(t) && et.Before(t) {
			isRange = true
		}
	}
	if isRange || len(te.Range) == 0 {
		te.BaseElement.Exec(d)
	}
}

type ArrivalElement struct { // 도착 알림을 위한 filter
	BaseElement
	Done float64 `json:"done"`
}

func (ae *ArrivalElement) Exec(d *model.LogicData) {
	
	log.Println("!!!!in ArrivalElement.Exec !!!!")
	done, ok := d.Values["done"]
	if !ok {
		log.Println("in ArrivalElement.Exec, not ok")
		return
	}
	isArrived := false

	if ae.Done == done {
		isArrived = true
	}
	if isArrived {
		ae.BaseElement.Exec(d)
	}
}
