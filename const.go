package api

var (
	ReqArgsErr             = Response{Code: 400}
	ModelNotExistsErr      = Response{Code: 1000}
	CallSteamCompletionErr = Response{Code: 1001}
)
