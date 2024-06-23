package binding

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
	MIMEYAML2             = "application/yaml"
	MIMETOML              = "application/toml"
)

type BindingBody interface {
	BindBody([]byte, any) error
}

var (
	JSON BindingBody = jsonBinding{}
)

func Default(contentType string) BindingBody {
	switch contentType {
	case MIMEJSON:
		return JSON
	default: // case MIMEPOSTForm:
		return nil
	}
}

//func validate(obj any) error {
//	if Validator == nil {
//		return nil
//	}
//	return Validator.ValidateStruct(obj)
//}
