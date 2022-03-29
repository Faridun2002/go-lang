package pkg

func ForError(err error) {
	if err != nil {
		panic(err)
	}
}