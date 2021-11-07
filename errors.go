package main

func HandleError(err error) {
	if err != nil {
		Error(err)
		panic(err)
	}
	return
}

