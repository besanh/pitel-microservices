package response

type PaginationBodyResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Total   int    `json:"total"`
}

type PaginationResponse[T any] struct {
	Body PaginationBodyResponse[T]
}

type GenericResponse[T any] struct {
	// Status int
	Body BodyResponse[T]
}

type BodyResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func OK[T any](data T, msgs ...string) (res *GenericResponse[T]) {
	msg := "success"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	res = &GenericResponse[T]{
		// Status: http.StatusOK,
		Body: BodyResponse[T]{
			Code:    "OK",
			Message: msg,
			Data:    data,
		},
	}
	return
}

func OK_Only(msgs ...string) (res *GenericResponse[any]) {
	msg := "success"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	res = &GenericResponse[any]{
		// Status: http.StatusOK,
		Body: BodyResponse[any]{
			Code:    "OK",
			Message: msg,
		},
	}
	return
}

func Pagination[T any](total int, data T, msgs ...string) (res *PaginationResponse[T]) {
	msg := "success"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	res = &PaginationResponse[T]{
		// Status: http.StatusOK,
		Body: PaginationBodyResponse[T]{
			Code:    "OK",
			Message: msg,
			Data:    data,
			Total:   total,
		},
	}
	return
}

type MediaResponse struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

type IdResponse struct {
	Id string `json:"id"`
}
