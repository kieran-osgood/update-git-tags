package internal

func ExamplePrintInfo() {
	PrintInfo("%s", "Let's PrintInfo")

	// Output: [1;34mLet's PrintInfo[0m
	//
}

func ExamplePrintNotice() {
	PrintNotice("%s", "PrintNotice will happen one day")

	// Output: [1;36mPrintNotice will happen one day[0m
	//
}

func ExamplePrintWarning() {
	PrintWarning("%s", "Now PrintWarning")

	// Output: [1;33mNow PrintWarning[0m
	//
}

func ExamplePrintError() {
	PrintError("%s", "Hope I never have to PrintError")

	// Output: [1;31mHope I never have to PrintError[0m
	//
}

func ExamplePrintDebug() {
	PrintDebug("%s", "I suppose debugging is inevitable!")

	// Output: [0;36mI suppose debugging is inevitable![0m
	//
}
