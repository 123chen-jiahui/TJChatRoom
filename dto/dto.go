package dto

type UserForCreationDto struct {
	Account  string `json:"account"`
	Passwd   string `json:"passwd"`
	NickName string `json:"nickName"`
}

type FriendDto struct {
	Account string `json:"account"`
}
