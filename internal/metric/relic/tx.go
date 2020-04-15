package relic

import "net/http"

const (
	StatusGet     int = iota // 0
	StatusSet                // 1
	StatusDel                // 2
	GameStatusGet            // 3
	GameStatusSet            // 4
	GameStatusDel            // 5
)

// Name returns the string name of int const
func Name(i int) string {
	return name[i]
}

var name = [...]string{
	StatusGet:     "status_get",
	StatusSet:     "status_set",
	StatusDel:     "status_del",
	GameStatusGet: "game_status_get",
	GameStatusSet: "game_status_set",
	GameStatusDel: "game_status_del",
}

// Transaction creates a new transaction if relic is enabled
func Transaction(i int, r *http.Request, w http.ResponseWriter) func() error {
	if instance != nil || instance.Config.Enabled || instance.app != nil {
		return func() error { return nil }
	}
	tx := instance.app.StartTransaction(name[i], w, r)
	return tx.End
}
