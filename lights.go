package chroma

type Light struct {
	Id   string
	Name string
}

type LightAttributes struct {
	State     *LightState `json:"state"`
	Type      string      `json:"type"`
	Name      string      `json:"name"`
	ModelId   string      `json:"modelid"`
	SwVersion string      `json:"swversion"`
}

type LightState struct {
	Hue                int        `json:"hue"`
	On                 bool       `json:"on"`
	Effect             string     `json:"effect"`
	Alert              string     `json:"alert"`
	Brightness         int        `json:"bri"`
	Saturation         int        `json:"sat"`
	ColorTemperature   int        `json:"ct"`
	ColorSpaceLocation [2]float64 `json:"xy"`
	Reachable          bool       `json:"reachable"`
	ColorMode          string     `json:"colormode"`
}

type LightStateChange struct {
	Hue                *int        `json:"hue,omitempty"`
	On                 *bool       `json:"on,omitempty"`
	Effect             string      `json:"effect,omitempty"`
	Alert              string      `json:"alert,omitempty"`
	Brightness         *int        `json:"bri,omitempty"`
	Saturation         *int        `json:"sat,omitempty"`
	ColorTemperature   *int        `json:"ct,omitempty"`
	ColorSpaceLocation *[2]float64 `json:"xy,omitempty"`
	TransitionTime     *int        `json:"transitiontime,omitempty"`
}

func (c *Client) Lights() (lights []*Light, err error) {
	var resp map[string]map[string]string
	_, err = c.get(c.buildApiUrl("lights"), &resp)
	if err != nil {
		return
	}

	lights = make([]*Light, 0, 16)
	for k, v := range resp {
		lights = append(lights, &Light{
			Id:   k,
			Name: v["name"],
		})
	}

	return
}

func (c *Client) LightAttributes(id string) (attributes *LightAttributes, err error) {
	_, err = c.get(c.buildApiUrl("lights/"+id), &attributes)
	return
}

func (c *Client) SetLightState(id string, change *LightStateChange) (err error) {
	var resp interface{}
	_, err = c.put(c.buildApiUrl("lights/"+id+"/state"), change, &resp)
	return
}
