package cast

const (
	fHeader int = 1 << iota
	fParam
	fResponse
	fTiming
	fStd = fHeader | fParam | fResponse | fTiming
)
