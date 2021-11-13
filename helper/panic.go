package helper

func PanicIfError(err error, errType ...interface{}) {
	if err != nil {
		if len(errType) > 0 {
			panic(errType[0])
		}
		panic(err)
	}
}
