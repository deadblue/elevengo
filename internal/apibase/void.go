package apibase

// VoidResult can be used as result type parameter in ApiSpec, when we do not
// care about result data from API.
type VoidResult struct{}

type VoidApiSpec struct {
	JsonApiSpec[VoidResult, BasicResp]
}
