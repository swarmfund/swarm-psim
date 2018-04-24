package resources

// TODO Comment
// TODO Consider moving the type into listener package
type TransactionEvent struct {
	Transaction *Transaction
	Meta        PageMeta
}
