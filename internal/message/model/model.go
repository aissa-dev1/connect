package messagemodel

type ReadStatus = string

type Message struct {
	Id         int         `json:"id"`
	SenderId   int         `json:"senderId"`
	ReceiverId int         `json:"receiverId"`
	Text       string      `json:"text"`
	ReadStatus *ReadStatus `json:"readStatus"`
}
