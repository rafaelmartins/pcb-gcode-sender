package response

type ResponseHandler interface {
	Supports(data string) bool
	Handle(data string) error
}

type ResponseHandlers []ResponseHandler

func (r ResponseHandlers) Lookup(data string) ResponseHandler {
	if len(data) == 0 {
		return nil
	}

	for _, handler := range r {
		if handler.Supports(data) {
			return handler
		}
	}

	return nil
}
