package esclient

type DebugHandler struct {
	requestData  []byte
	responseData []byte
	debugData    map[string]interface{}
	wasUpdated   bool
}

func NeWDebugHandler() *DebugHandler {
	return &DebugHandler{
		requestData:  []byte{},
		responseData: []byte{},
		debugData:    map[string]interface{}{},
	}
}

// >>>>>>>>>> Interface function
func (d *DebugHandler) Add(key string, val interface{}) {
	d.debugData[key] = val
}

func (d *DebugHandler) Get(key string) interface{} {
	if val, find := d.debugData[key]; find {
		return val
	}

	return nil
}

func (d *DebugHandler) All() map[string]interface{} {
	return d.debugData
}

func (d *DebugHandler) Request() []byte {
	return d.requestData
}

func (d *DebugHandler) Response() []byte {
	return d.responseData
}

func (d *DebugHandler) WasUpdated() bool {
	return d.wasUpdated
}

func (d *DebugHandler) SetRequest(b []byte) {
	d.requestData = b
	d.wasUpdated = true
}

func (d *DebugHandler) SetResponse(b []byte) {
	d.responseData = b
	d.wasUpdated = true
}
