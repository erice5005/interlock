package interlock

import (
	"encoding/json"
	"log"
	"time"
)

type MSGType int64

const (
	PING       MSGType = 0
	PONG       MSGType = 1
	DATA       MSGType = 2
	CONNECT    MSGType = 3
	DISCONNECT MSGType = 4
)

type DataFrame struct {
	SrcID   string
	Kind    MSGType
	Content interface{}
	Created time.Time
}

const DELIM byte = '\n'

func (df *DataFrame) ToBytes() []byte {
	buf, _ := json.Marshal(df)
	buf = append(buf, DELIM) //Customized Delimiter
	return buf
}

func ParseDataFrame(raw []byte) *DataFrame {
	var out *DataFrame
	err := json.Unmarshal(raw, &out)
	if err != nil {
		log.Printf("Got: %v\n", string(raw))
		log.Printf("Parse out: %v\n", err)
		return &DataFrame{}
	}
	return out
}

type DataPipe struct {
	Input  chan DataFrame
	Output chan DataFrame
}

func NewDataPipe() *DataPipe {
	return &DataPipe{
		Input:  make(chan DataFrame, 10),
		Output: make(chan DataFrame, 10),
	}
}
