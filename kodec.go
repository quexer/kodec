/**
 * common codec library for chat application
 */
package kodec

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
)

const (
	BAG_ID_MSG byte = iota //bag id in for message
	BAG_ID_CMD
)

func BuildCmd(tp Cmd_Type, desc string, tick int64) *Cmd {
	m := &Cmd{
		Tp:  &tp,
		Txt: proto.String(desc),
		Ct:  proto.Int64(tick),
	}
	return m
}

func BuildMeta(tp Meta_Type, txt string) *Meta {
	m := &Meta{
		Tp:  &tp,
		Txt: proto.String(txt),
	}
	return m
}

func BuildMessage(from int, to string, data []byte, tp Msg_Type, tick int64) *Msg {
	m := &Msg{
		From: proto.Int64(int64(from)),
		To:   proto.String(to),
		Tp:   &tp,
		D:    data,
		Ct:   proto.Int64(tick),
	}
	return m
}

func Boxing(m proto.Message) ([]byte, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}

	var bagId byte
	switch m.(type) {
	case *Msg:
		bagId = BAG_ID_MSG
	case *Cmd:
		bagId = BAG_ID_CMD
	default:
		return nil, fmt.Errorf("unknown frame")
	}
	return append([]byte{bagId}, b...), nil
}

func Unboxing(data []byte) (interface{}, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("warning: bad frame")
	}
	bagId := data[0]
	data = data[1:]

	switch bagId {
	case BAG_ID_MSG:
		if len(data) < 1 {
			return nil, fmt.Errorf("warning: bad message frame")
		}
		upMsg := new(Msg)
		if err := proto.Unmarshal(data, upMsg); err != nil {
			return nil, fmt.Errorf("parse msg err %v", err)
		}

		return upMsg, nil
	default:
		return nil, fmt.Errorf("warning: unknown bag id: %v", bagId)
	}
}
