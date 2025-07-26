package blockmodel

type Block struct {
	Id        int `json:"id"`
	BlockerId int `json:"blockerId"`
	BlockedId int `json:"blockedId"`
}
