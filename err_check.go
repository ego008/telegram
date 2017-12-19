package main

func errCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
