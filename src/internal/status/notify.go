package status

type NotifyStatus int16

const (
	NotifyPending NotifyStatus = iota
	NotifySent
	NotifyFailed
)
