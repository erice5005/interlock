package interlock

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	var b64 string
	if df.Content != nil {
		b64 = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", df.Content)))

	}

	sendFrame := DataFrame{
		SrcID:   df.SrcID,
		Kind:    df.Kind,
		Content: b64,
		Created: df.Created,
	}

	buf, _ := json.Marshal(sendFrame)
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

	if out.Content != nil {
		cleanContent, err := base64.StdEncoding.DecodeString(out.Content.(string))
		if err != nil {
			log.Printf("b64 decode err: %v\n", err)

		} else {
			out.Content = string(cleanContent)
		}
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
