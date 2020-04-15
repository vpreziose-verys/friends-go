package ev

import (
	"fmt"
)

type X byte
type Op byte

const (
	XGatePub X = iota
	XGatePrv
	XAccount
	XStatus
	XProvider
	XManager
)
const (
	OpGet Op = iota
	OpPut
	OpDel
	OpExp
)

func (x X) String() string  { return xString[x] }
func (x Op) String() string { return opString[x] }

type E struct {
	Dst X
	Did string
	Op  Op
	Src X
	Sid string
}

func (e E) String() string {
	const format = `"dst":"%s","did":"%s","op":"%s","src":"%s","sid":"%s","data":%s`
	return fmt.Sprintf(format, e.Dst, e.Did, e.Op, e.Src, e.Sid, e.dataset())
}

func (e E) dataset() string {
	return "{}"
}

var valid = [...]E{
	{XAccount, "", OpPut, XStatus, ""},
	{XAccount, "", OpDel, XStatus, ""},
	{XAccount, "", OpExp, XStatus, ""},
}

var (
	xString = [...]string{
		XGatePub: "gate-pub",
		XGatePrv: "gate-prv",
		XAccount: "account",
		XStatus:  "status",
	}
	opString = [...]string{
		OpGet: "get",
		OpPut: "put",
		OpDel: "del",
		OpExp: "exp",
	}
)
