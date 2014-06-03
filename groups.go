package chroma

type Group struct {
	Id   string
	Name string
}

type GroupAttributes struct {
	// Action *Action `json:"action"` -- not implemented yet, low priority

	Name   string   `json:"name"`
	Lights []string `json:"lights"`
}

type GroupStateChange struct {
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

func (c *Client) Groups() (groups []*Group, err error) {
	var resp map[string]map[string]string
	_, err = c.get(c.buildApiUrl("groups"), &resp)
	if err != nil {
		return
	}

	groups = make([]*Group, 0, 16)
	for k, v := range resp {
		groups = append(groups, &Group{
			Id:   k,
			Name: v["name"],
		})
	}

	return
}

func (c *Client) GroupAttributes(id string) (attributes *GroupAttributes, err error) {
	_, err = c.get(c.buildApiUrl("groups/"+id), &attributes)
	return
}

func (c *Client) SetGroupState(id string, change *GroupStateChange) (err error) {
	var resp interface{}
	_, err = c.put(c.buildApiUrl("groups/"+id+"/action"), change, &resp)
	return
}
