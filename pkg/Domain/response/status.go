package response

type Status int

const (
	InternalServerError Status = iota
	DBQueryError
	DBExecutionError
	DBRowsError
	DBLastRowIdError

	RequestError
	RequestDoError
	ReadRequestError

	EmailAlreadyExists
	NickNameAlreadyExists
	UserNotFound
	UserDontExist

	SuccesfulCreation
	UserFound
	CreationFailure

	InvalidPermissions
	InvalidEmailFormat
	IncorrectPassword
	OriginNotAllowed
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

		"RequestError",
		"RequestDoError",
		"ReadRequestError",

		"EmailAlreadyExists",
		"NickNameAlreadyExists",
		"UserNotFound",
		"UserDontExist",

		"SuccesfulCreation",
		"UserFound",
		"CreationFailure",

		"InvalidPermissions",
		"InvalidEmailFormat",
		"IncorrectPassword",
		"OriginNotAllowed",
		"RequestTimeOut",
		"BadRequest",
		"Unknown",
	}[s]
}

func (s Status) Index() int {
	return int(s)
}
