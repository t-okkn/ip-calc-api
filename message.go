package main

type SuccessMessage struct {
	Code    string `json:"success"`
	Message string `json:"message"`
}

type ErrorMessage struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

var (
	sucUpdateDone = SuccessMessage {
		Code   : "S001",
		Message: "更新に成功しました",
	}

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

	errTheQuestionIsNotExist = ErrorMessage{
		Code   : "E120",
		Message: "リクエストされた問題番号はまだ存在していません",
	}

	errTheQuestionIsTerminated = ErrorMessage{
		Code   : "E121",
		Message: "既に回答済みの問題を更新することはできません",
	}
)