package utils

var (
	SUCCESS = Response{Code: 200, Message: `sucesss`}

	ErrInnerServer = Response{Code: 500}
	ErrReqArgs     = Response{Code: 400, Message: `req args error`}

	ErrNoFound             = Response{Code: 404}
	ModelNotExistsErr      = Response{Code: 1000}
	CallSteamCompletionErr = Response{Code: 1001}

	ErrScannerIsRunning = Response{Code: 1002, Message: `scanner is running`}
	ErrScanPathNotExist = Response{Code: 1003, Message: `scan path not exist`}
	ErrSaveFile         = Response{Code: 1004, Message: `save file error`}
	ErrParseEpubFile    = Response{Code: 1005, Message: `parse epub file error`}
)

const (
	//config 目录
	ConfigBucket = `config`
)
