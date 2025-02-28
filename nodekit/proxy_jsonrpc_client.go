package nodekit

import (
	"context"
	"strings"
)

const (
	Name            = "proxy"
	JSONRPCEndpoint = "/rpc"
)

type JSONRPCClient struct {
	requester *EndpointRequester
}

func NewJSONRPCClient(uri string) *JSONRPCClient {
	uri = strings.TrimSuffix(uri, "/")
	uri += JSONRPCEndpoint
	req := NewRequester(uri, Name)
	return &JSONRPCClient{requester: req}
}

type SubmitMsgTxArgs struct {
	Data []byte `json:"Data"`
}

type SubmitMsgTxReply struct {
	TxID string `json:"txId"`
}

func (j *JSONRPCClient) SubmitMsgTx(ctx context.Context, data []byte) (string, error) {
	resp := new(SubmitMsgTxReply)

	err := j.requester.SendRequest(ctx,
		"submitMsgTx",
		&SubmitMsgTxArgs{
			Data: data,
		},
		resp,
	)

	if err != nil {
		return "", err
	}

	return resp.TxID, nil
}
