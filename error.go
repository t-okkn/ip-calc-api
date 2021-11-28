package main

type ErrorMessage struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

var (
	errCannotConnectDB = ErrorMessage{
		Code   : "E001",
		Message: "DBと接続できません",
	}

	errFailedOperateData = ErrorMessage{
		Code   : "E002",
		Message: "データの操作に失敗しました",
	}

	errFailedGetData = ErrorMessage{
		Code   : "E003",
		Message: "データの取得に失敗しました",
	}

	errInvalidRequestedURL = ErrorMessage{
		Code   : "E100",
		Message: "リクエストされたURLが不正です",
	}

	errInvalidRequestedData = ErrorMessage{
		Code   : "E101",
		Message: "リクエストされたデータが不正です",
	}
)