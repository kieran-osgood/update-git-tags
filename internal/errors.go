package internal

func HandleError(err error) {
	if err != nil {
		Error(err.Error())
		panic(err)
	}
	return
}
