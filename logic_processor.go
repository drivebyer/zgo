package zgo

// LogicProcessor hold XXXXXLogicProcessor which is servive layer defined.
// The first letter of Hander must be upper-case, since implement LogicProcessor
// in service layer, always in two different package.
type LogicProcessor interface {
	Handler(*Connection, []byte)
}
