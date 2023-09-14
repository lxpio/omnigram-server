package utils

var (
	SUCCESS = Response{Code: 200, Message: `sucesss`}

	ErrSession = Response{Code: 401, Message: `session is invalid`}

	ErrSessionTimeout = Response{Code: 401, Message: `session is timeout`}
	ErrUnauthorized   = Response{Code: 401, Message: `unauthorized`}

	ErrAPIKey = Response{Code: 401, Message: `api key is invalid`}

	ErrInnerServer = Response{Code: 500}
	ErrReqArgs     = Response{Code: 400, Message: `req args error`}

	ErrNoFound             = Response{Code: 404}
	ModelNotExistsErr      = Response{Code: 1000}
	CallSteamCompletionErr = Response{Code: 1001}

	ErrScannerIsRunning = Response{Code: 1002, Message: `scanner is running`}
	ErrScanPathNotExist = Response{Code: 1003, Message: `scan path not exist`}
	ErrSaveFile         = Response{Code: 1004, Message: `save file error`}
	ErrParseEpubFile    = Response{Code: 1005, Message: `parse epub file error`}

	ErrSaveToken   = Response{Code: 1006, Message: `save token error`}
	ErrDeleteToken = Response{Code: 1007, Message: `delete token error`}
	ErrGetTokens   = Response{Code: 1008, Message: `get tokens failed`}
)

const (
	//config 目录
	ConfigBucket = `config`
)
