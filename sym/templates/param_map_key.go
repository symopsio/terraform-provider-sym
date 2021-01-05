package templates

type ParamMapKey struct {
	Map *HCLParamMap
	Key string
}

func (pmf *ParamMapKey) Value() string {
	return pmf.Map.Params[pmf.Key]
}

func (pmf *ParamMapKey) addDiagWithDetail(summary string, detail string) {
	pmf.Map.addDiagWithDetail(pmf.Key, summary, detail)
}

func (pmf *ParamMapKey) addDiag(summary string) {
	pmf.Map.addDiag(pmf.Key, summary)
}

func (pmf *ParamMapKey) checkError(summary string, err error) {
	pmf.Map.checkError(pmf.Key, summary, err)
}
