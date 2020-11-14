package main

//Field is field struct
type Field struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}
type Request struct {
	Instruments []string               `json:"instruments"`
	Fields      []Field                `json:"fields"`
	Parameters  map[string]interface{} `json:"parameters"`
}
type W struct {
	Requests []Request `json:"requests"`
}

type EikonRequest struct {
	Entity struct {
		E string `json:"E"`
		W W      `json:"W"`
	} `json:"Entity"`

	ID string `json:"ID"`
}

func newEikonRequest(instruments []string, fields []string, period string) *EikonRequest {
	req := EikonRequest{}
	req.ID = "123"
	req.Entity.W.Requests = make([]Request, 0)
	req.Entity.E = "DataGrid_StandardAsync"
	ek := Request{
		Instruments: instruments,
	}
	if period != "" {
		ek.Parameters = map[string]interface{}{
			"Period": period,
		}
	}
	_fields := make([]Field, 0)
	for i := range fields {
		_fields = append(_fields, Field{Name: fields[i]})
	}
	ek.Fields = _fields
	req.Entity.W.Requests = append(req.Entity.W.Requests, ek)
	return &req
}
