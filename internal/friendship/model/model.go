package friendshipmodel

type FriendshipStatus = string

type Friendship struct {
	Id          int               `json:"id"`
	RequesterId int               `json:"requesterId"`
	ReceiverId  int               `json:"receiverId"`
	Status      *FriendshipStatus `json:"status"`
}

type FriendshipRequest struct {
	RequesterId int `json:"requesterId"`
}
