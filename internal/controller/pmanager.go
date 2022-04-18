package controller

import (
	"context"
	"encoding/json"
	"errors"
	v1 "go-ops/api/v1"
	"go-ops/peer"
	"time"

	"github.com/luxingwen/pnet/log"
	"github.com/luxingwen/pnet/stat"

	"github.com/gogo/protobuf/proto"
	"github.com/luxingwen/pnet/protos"
)

var PeerManagaer = new(peerManager)

type peerManager struct {
}

func (self *peerManager) GetNodes(ctx context.Context, nodeReq *v1.NodeReq) (res *v1.NodeRes, err error) {
	res = new(v1.NodeRes)
	if nodeReq.NodeId == "" {
		rns, err := peer.GetOspPeer().GetLocalNode().GetNeighbors(nil)
		if err != nil {
			return nil, err
		}
		for _, item := range rns {
			res.Nodes = append(res.Nodes, item.Node.Node)
		}
		return res, nil
	}

	msg := peer.GetOspPeer().GetLocalNode().NewGetNeighborsMessage(nodeReq.NodeId)

	resMsg, ok, err := peer.GetOspPeer().SendMessageSync(msg, protos.RELAY, time.Minute)
	if err != nil {
		return
	}

	if !ok {
		err = errors.New("获取节点信息失败")
		return
	}

	nodes := &protos.Neighbors{}

	err = proto.Unmarshal(resMsg.Message, nodes)
	if err != nil {
		log.Error("unmarshal err:", err)
		return
	}

	res.Nodes = nodes.Nodes
	return
}

func (self *peerManager) NodeConnect(ctx context.Context, req *v1.NodeConnectReq) (res *v1.NodeOpRes, err error) {
	res = new(v1.NodeOpRes)
	if req.NodeId == "" {
		req.NodeId = peer.GetOspPeer().GetLocalNode().GetId()

		_, _, err = peer.GetOspPeer().GetLocalNode().Connect(&protos.Node{Addr: req.RemoteAddr})
		if err != nil {
			res.Msg = "链接节点失败"
			return
		}

		res.Msg = "OK"
		return
	}

	msg, err := peer.GetOspPeer().GetLocalNode().NewConnnetNodeMessage(req.NodeId, &protos.Node{Addr: req.RemoteAddr})
	if err != nil {
		return
	}

	resMsg, ok, err := peer.GetOspPeer().SendMessageSync(msg, protos.RELAY, time.Minute)
	if err != nil {
		return
	}

	if !ok {
		err = errors.New("链接节点失败")
		return
	}

	res.Msg = string(resMsg.Message)
	return

}

func (self *peerManager) StopConnect(ctx context.Context, req *v1.NodeStopReq) (res *v1.NodeOpRes, err error) {
	res = new(v1.NodeOpRes)
	if req.NodeId == "" {
		rm, err := peer.GetOspPeer().GetLocalNode().GetNeighbors(nil)
		if err != nil {
			return res, err
		}
		v, ok := rm[req.RemoteId]
		if ok {
			v.Stop(errors.New("STOP_REMOTENODE"))
		} else {
			err = errors.New("没有找到节点:" + req.RemoteId)
			return res, err
		}

		res.Msg = "OK"
		return res, nil
	}

	msg, err := peer.GetOspPeer().GetLocalNode().NewStopConnnetNodeMessage(req.NodeId, &protos.Node{Id: req.RemoteId})
	if err != nil {
		return
	}

	resMsg, ok, err := peer.GetOspPeer().SendMessageSync(msg, protos.RELAY, time.Minute)
	if err != nil {
		return
	}

	if !ok {
		err = errors.New("停止链接节点失败")
		return
	}

	res.Msg = string(resMsg.Message)
	return
}

func (self *peerManager) PeerStat(ctx context.Context, req *v1.NodeStatReq) (res *v1.NodeStatRes, err error) {
	res = new(v1.NodeStatRes)
	if req.NodeId == "" {
		res.NodeStat = stat.GetNodeStat()
		res.PeerId = peer.GetOspPeer().GetLocalNode().GetId()
		return
	}

	resMsg, _, err := peer.GetOspPeer().SendMessageSync(peer.GetOspPeer().GetLocalNode().NewNodeStatMessage(req.NodeId), protos.RELAY, 0)
	if err != nil {
		return
	}

	err = json.Unmarshal(resMsg.Message, res)

	return
}
