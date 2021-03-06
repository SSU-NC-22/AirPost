package logicService

import (
	"errors"
	"fmt"
	"log"

	"github.com/eunnseo/AirPost/logic-core/domain/model"
	"github.com/eunnseo/AirPost/logic-core/logicService/logic"
)

type logicService struct {
	mux
}

type mux struct {
	chTable map[int]map[int]chan model.LogicData
}

func NewLogicService() *logicService {
	return &logicService{
		mux{
			chTable: make(map[int]map[int]chan model.LogicData),
		},
	}
}

func (m *mux) CreateAndStartLogic(l *model.Logic) error {
	listen := make(chan model.LogicData, 100)
	
	lchs, ok := m.chTable[l.NodeID]
	if !ok {
		m.chTable[l.NodeID] = make(map[int]chan model.LogicData)
		lchs, ok = m.chTable[l.NodeID]
		if !ok {
			return errors.New("cannot make logic channel")
		}
	}
	if _, ok := lchs[l.ID]; ok {
		close(listen)
		return errors.New("already exist logic event")
	}
	lchs[l.ID] = listen

	elems, err := logic.BuildLogic(l)
	
	if err != nil {
		log.Println("end BuildLogic, return error")
		return err
	}
	go func() {
		for d := range listen {
			elems.Exec(&d)
		}
	}()

	return nil
}

func (m *mux) RemoveLogic(nid, lid int) error {
	ch, ok := m.chTable[nid][lid]
	if !ok {
		return fmt.Errorf("GetLogicChans : cannot find listen channels")
	}
	close(ch)
	delete(m.chTable[nid], lid)
	if len(m.chTable[nid]) == 0 {
		delete(m.chTable, nid)
	}
	return nil
}

func (m *mux) GetLogicChans(nid int) (map[int]chan model.LogicData, error) {
	lchs, ok := m.chTable[nid]
	if !ok || len(lchs) == 0 {
		return nil, fmt.Errorf("GetLogicChans : cannot find listen channels")
	}
	return lchs, nil
}
