package handles

import (
	"net/http"
	"nodestore/app/db"
	"nodestore/app/module"
	"nodestore/app/response"
	"time"
)

type Node struct{}

// ==================== 节点管理 ====================
func (Node) AddNodeHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	addr := r.URL.Query().Get("addr")
	if addr == "" {
		response.Failed(r, 400, "addr empty")
		// http.Error(w, "addr empty", 400)
		return
	}

	now := time.Now()
	err = db.DB.Exec(`INSERT INTO nodes (user_id,address,create_time) VALUES (?,?,?)`,
		claims.UserID, addr, now).Error
	if err != nil {
		response.Failed(r, 500, err.Error())
		// http.Error(w, err.Error(), 500)
		return
	}
	response.Success(r, "node added", nil)
	// w.Write([]byte("node added"))
}

func (Node) RemoveNodeHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	addr := r.URL.Query().Get("addr")
	err = db.DB.Exec(`DELETE FROM nodes WHERE user_id=? AND address=?`, claims.UserID, addr).Error
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	response.Success(r, "node removed", nil)
	// w.Write([]byte("node removed"))
}

func (Node) ListNodesHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	var nodes []string
	err = db.DB.Raw(`SELECT address FROM nodes WHERE user_id=?`, claims.UserID).Scan(&nodes).Error
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	response.Success(r, "ok", nodes)
	// json.NewEncoder(w).Encode(map[string]interface{}{"nodes": nodes})
}
