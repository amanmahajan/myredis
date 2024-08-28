package core

var KeySpaceStats [4]map[string]int

func UpdateDbStats(dbIndex int, key string, value int) {
	KeySpaceStats[dbIndex][key] = value
}
