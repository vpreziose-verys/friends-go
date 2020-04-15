package client

import (
	"encoding/json"
	"fmt"
)

type Reply struct {
	PlatformReply
}

func (r *Reply) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &r.PlatformReply); err != nil {
		return err
	}
	if c := r.Platform.Code; c < 2000 {
		return PlatformError{ctx: "bad return code", Code: c}
	}
	return nil
}

type PlatformReply struct {
	Platform struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message,omitempty"`
		Body    interface{} `json:"response,omitempty"`
	} `json:"platform"`
}

type PlatformError struct {
	Code int
	Msg  string
	ctx  string
}

func (p PlatformError) Error() string {
	return fmt.Sprintf("platform: %s: %d", p.ctx, p.Code)
}
