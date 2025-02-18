package GameDatamodel

var UserItemList map[int]int

//var User map[]

func ReturnUserItemList() map[int]int {
	if UserItemList == nil {
		UserItemList = make(map[int]int) // ✅ 初始化 map，防止 nil 访问
	}
	return UserItemList
}
