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
}

func DeleteGroup(groupId string, account string) bool {
	return db.CheckAndRemoveGroup(groupId, account)
}
