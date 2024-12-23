package responder

type MetaField struct {
	Meta MetaObject `json:"meta,omitempty"`
}

type MetaObject struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewMetaField(status int, msg string) *MetaField {
	return &MetaField{
		Meta: MetaObject{
			Status:  status,
			Message: msg,
		},
	}
}

type DataField struct {
	Data any `json:"data"`
}

type ErrorsArray struct {
	Errors []ErrorObject `json:"errors,omitempty"`
}

type ErrorObject struct {
	Code   *string `json:"code,omitempty"`
	Title  string  `json:"title"`
	Detail *string `json:"detail,omitempty"`
}

func NewErrorsArray(errors ...ErrorObject) *ErrorsArray {
	return &ErrorsArray{Errors: errors}
}

type Response struct {
	*MetaField
	*ErrorsArray
}
