package name_server

var defaultAddr = "192.168.1.36:9989"

func getNode(uid string, isRepair bool) (string, error) {
	if uid == NullString {
		return getHealthiestNode(isRepair)
	}
	addr, err := clusterData.getKeyAddr(uid, isRepair)
	if err != nil {
		return getHealthiestNode(isRepair)
	}
	return addr, nil
}

func getNodeByPass(name, passWord string, isRepair bool) (string, error) {
	uid, err := getUid(name, passWord)
	if err != nil {
		return NullString, err
	}
	return getNode(uid, isRepair)
}

func getUid(name, passWord string) (string, error) {
	return name, nil
}

func getHealthiestNode(isRepair bool) (string, error) {
	addr, err := clusterData.getHealthiestNodeAddr(isRepair)
	if err != nil || addr == NullString {
		return defaultAddr, err
	}
	return addr, nil
}
