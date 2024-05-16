package api

type ShowHiddenSpec struct {
	_VoidApiSpec
}

func (s *ShowHiddenSpec) Init(password string) *ShowHiddenSpec {
	s._VoidApiSpec.Init("https://115.com/?ct=hiddenfiles&ac=switching")
	s.form.Set("show", "1").
		Set("valid_type", "1").
		Set("safe_pwd", password)
	return s
}

type HideHiddenSpec struct {
	_VoidApiSpec
}

func (s *HideHiddenSpec) Init() *HideHiddenSpec {
	s._VoidApiSpec.Init("https://115.com/?ct=hiddenfiles&ac=switching")
	s.form.Set("show", "0").
		Set("valid_type", "1").
		Set("safe_pwd", "")
	return s
}
