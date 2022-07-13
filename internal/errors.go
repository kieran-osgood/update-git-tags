package internal

func HandleError(err error) {
	if err != nil {
		PrintError(err.Error())
		panic(err)
	}
}
