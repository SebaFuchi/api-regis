package response

type Status int

const (
	InternalServerError Status = iota
	DBQueryError
	DBExecutionError
	DBRowsError
	DBLastRowIdError
	DBScanError

	EmailAlreadyExists
	UserNotFound
	UserDontExist

	SuccesfulCreation
	SuccesfulLogin
	UserFound
	CreationFailure

	InvalidEmailFormat
	IncorrectPassword
	RequestTimeOut
	BadRequest
	Unknown
)

func (s Status) String() string {
	return [...]string{
		"InternalServerError",
		"DBQueryError",
		"DBExecutionError",
		"DBRowsError",
		"DBLastRowIdError",
		"DBScanError",

		"EmailAlreadyExists",
		"UserNotFound",
		"UserDontExist",

		"SuccesfulCreation",
		"SuccesfulLogin",
		"UserFound",
		"CreationFailure",

		"InvalidEmailFormat",
		"IncorrectPassword",
		"RequestTimeOut",
		"BadRequest",
		"Unknown",
	}[s]
}

func (s Status) Index() int {
	return int(s)
}
