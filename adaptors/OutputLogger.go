package adaptors

type Writer interface{
	Write(string)
	RegisterOutputFile(string)
}
