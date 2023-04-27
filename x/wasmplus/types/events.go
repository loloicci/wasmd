package types

const (
	// events from callable point
	CallablePointEventType = "wasm-callablepoint"
	// prefix for custom events from callable point
	CustomCallablePointEventPrefix = "wasm-callablepoint-"
)

// event attributes returned from contract execution
const (
	AttributeKeyCallstack = "_callstack"
	AttributeKeyCodeIDs   = "code_ids"
	AttributeKeyFeature   = "feature"
)
