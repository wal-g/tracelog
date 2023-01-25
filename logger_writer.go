package tracelog

type LoggerWriter interface {
	Log(Fields)
}
