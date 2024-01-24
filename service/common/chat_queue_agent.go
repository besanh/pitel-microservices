package common

func CheckAgentExist(id string, agents []string) bool {
	for _, agent := range agents {
		if agent == id {
			return true
		}
	}
	return false
}
