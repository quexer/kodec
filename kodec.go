/**
 * common codec library for chat application
 */
package kodec

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

const (
	BAG_ID_MSG byte = iota //bag id in for message
	BAG_ID_CMD
	BAG_ID_ACK
)

func BuildAck(id string) *Ack {
	m := &Ack{
		Id: &id,
	}
	return m
}

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

func SetMessageId(m *Msg, id string) {
	m.Id = proto.String(id)
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
	case *Ack:
		bagId = BAG_ID_ACK
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
	case BAG_ID_ACK:
		ack := new(Ack)
		if err := proto.Unmarshal(data, ack); err != nil {
			return nil, fmt.Errorf("parse ack err %v", err)
		}
		return ack, nil
	case BAG_ID_CMD:
		cmd := new(Cmd)
		if err := proto.Unmarshal(data, cmd); err != nil {
			return nil, fmt.Errorf("parse cmd err %v", err)
		}
		return cmd, nil
	default:
		return nil, fmt.Errorf("warning: unknown bag id: %v", bagId)
	}
}
