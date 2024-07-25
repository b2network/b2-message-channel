package enums

type MessageType int64

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeCall
	MessageTypeSend
)

type MessageStatus int64

const (
	MessageStatusPending MessageStatus = iota
	MessageStatusBroadcast
	MessageStatusValid
	MessageStatusInvalid
)

type SignatureStatus int64

const (
	SignatureStatusPending SignatureStatus = iota
	SignatureStatusBroadcast
	SignatureStatusSuccess
	SignatureStatusFailed
	SignatureStatusInvalid
)
