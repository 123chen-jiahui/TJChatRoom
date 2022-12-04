package method

import (
	"github.com/db"
	"github.com/entity"
)

func AddGroup(group entity.Group) {
	db.InsertGroup(group)
}

func GetGroups(account string) []entity.Group {
	return db.GetGroups(account)
}

func AddMemberToGroup(groupId string, members []string) {
	db.AddMember(groupId, members)
}

func DeleteMembersFromGroup(groupId string, members []string) {
	db.PullMembers(groupId, members)
	// 删除群成员需要把被删除者的聊天记录删了，不然他再重新加入后依然能看到自己的聊天记录
	for _, m := range members {
		db.DeleteMessagesInGroup(m, groupId)
	}
}

func DeleteGroup(groupId string, account string) bool {
	return db.CheckAndRemoveGroup(groupId, account)
}
