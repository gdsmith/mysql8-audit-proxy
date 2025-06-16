package serverconfig

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/server"
)

type ConfigHandler struct {
	server.EmptyHandler
	db string
	*Manager
	err error
	res *mysql.Result
}

func (h *ConfigHandler) selectStmt(p *ParsedQuery) {
	h.res = &mysql.Result{}
	col, data, err := h.Select(p)
	if err != nil {
		h.err = err
		return
	}
	if len(col) == 0 {
		h.res = &mysql.Result{}
		return
	}
	r, _ := mysql.BuildSimpleResultset(col, data, false)
	h.res = &mysql.Result{Resultset: r}
}
func (h *ConfigHandler) insertStmt(p *ParsedQuery) {
	var n uint64
	n, h.err = h.Insert(p)
	h.res = &mysql.Result{AffectedRows: n}

}
func (h *ConfigHandler) updateStmt(p *ParsedQuery) {
	var n uint64
	n, h.err = h.Update(p)
	h.res = &mysql.Result{AffectedRows: n}
}
func (h *ConfigHandler) deleteStmt(p *ParsedQuery) {
	var n uint64
	n, h.err = h.Delete(p)
	h.res = &mysql.Result{AffectedRows: n}
}

func (h *ConfigHandler) handleQuery(query string) (*mysql.Result, error) {
	astNode, err := Parse(query)
	if err != nil {
		return nil, err
	}
	data := NewParsedQuery(
		map[string]func(p *ParsedQuery){
			SelectStmt: h.selectStmt,
			InsertStmt: h.insertStmt,
			UpdateStmt: h.updateStmt,
			DeleteStmt: h.deleteStmt,
		},
	)
	(*astNode).Accept(data)
	return h.res, h.err
}

func (h *ConfigHandler) HandleQuery(query string) (*mysql.Result, error) {
	return h.handleQuery(query)
}

func (h *ConfigHandler) HandleOtherCommand(cmd byte, data []byte) error {
	return mysql.NewError(mysql.ER_UNKNOWN_ERROR, fmt.Sprintf("command %d is not supported now", cmd))
}

func (h *ConfigHandler) UseDB(dbName string) error {
	h.db = dbName
	return nil
}

func (h *ConfigHandler) GetDB() string { return h.db }

func NewConfigHandler(m *Manager) *ConfigHandler {
	return &ConfigHandler{Manager: m}
}
