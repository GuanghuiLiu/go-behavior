package user_center

func start() {
	userCenter := newUserCenter("user_center")
	userCenter.Run(userCenter)
}
