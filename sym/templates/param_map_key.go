package templates

type ParamMapKey struct {
	Map *HCLParamMap
	Key string
}

func (pmf *ParamMapKey) Value() string {
	return pmf.Map.Params[pmf.Key]
}
