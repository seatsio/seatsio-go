package events

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/shared"
)

type Channels struct {
	Client *req.Client
}

type CreateChannelParams struct {
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Index   int32    `json:"index"`
	Objects []string `json:"objects"`
}

type CreateChannelParamsOption func(Params *CreateChannelParams)

type UpdateChannelParams struct {
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Objects []string `json:"objects"`
}

type changeChannelObjectsRequest struct {
	Objects []string `json:"objects"`
}

type replaceChannelsRequest struct {
	Channels []Channel `json:"channels"`
}

type channelSupportNS struct{}

var ChannelSupport channelSupportNS

func (channels *Channels) Create(eventKey string, params *CreateChannelParams) error {
	result, err := channels.Client.R().
		SetBody(params).
		SetPathParam("key", eventKey).
		Post("/events/{key}/channels")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) CreateMultiple(eventKey string, params *[]CreateChannelParams) error {
	result, err := channels.Client.R().
		SetBody(params).
		SetPathParam("key", eventKey).
		Post("/events/{key}/channels")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) Update(eventKey string, channelKey string, params UpdateChannelParams) error {
	result, err := channels.Client.R().
		SetBody(params).
		SetPathParam("eventKey", eventKey).
		SetPathParam("channelKey", channelKey).
		Post("/events/{eventKey}/channels/{channelKey}")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) Delete(eventKey string, channelKey string) error {
	result, err := channels.Client.R().
		SetPathParam("eventKey", eventKey).
		SetPathParam("channelKey", channelKey).
		Delete("/events/{eventKey}/channels/{channelKey}")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) AddObjects(eventKey string, channelKey string, objects []string) error {
	result, err := channels.Client.R().
		SetBody(changeChannelObjectsRequest{objects}).
		SetPathParam("eventKey", eventKey).
		SetPathParam("channelKey", channelKey).
		Post("/events/{eventKey}/channels/{channelKey}/objects")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) RemoveObjects(eventKey string, channelKey string, objects []string) error {
	result, err := channels.Client.R().
		SetBody(changeChannelObjectsRequest{objects}).
		SetPathParam("eventKey", eventKey).
		SetPathParam("channelKey", channelKey).
		Delete("/events/{eventKey}/channels/{channelKey}/objects")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) Replace(eventKey string, newChannels []Channel) error {
	result, err := channels.Client.R().
		SetBody(replaceChannelsRequest{newChannels}).
		SetPathParam("eventKey", eventKey).
		Post("/events/{eventKey}/channels/replace")
	return shared.AssertOkNoBody(result, err)
}

func (channelSupportNS) WithIndex(index int32) CreateChannelParamsOption {
	return func(params *CreateChannelParams) {
		params.Index = index
	}
}

func (channelSupportNS) WithObjects(objects []string) CreateChannelParamsOption {
	return func(params *CreateChannelParams) {
		params.Objects = objects
	}
}

func NewCreateChannelParams(key string, name string, color string, opts ...CreateChannelParamsOption) *CreateChannelParams {
	params := &CreateChannelParams{key, name, color, 0, []string{}}
	for _, opt := range opts {
		opt(params)
	}
	return params
}

func (channels *Channels) NewUpdateChannelParams(channel Channel) UpdateChannelParams {
	return UpdateChannelParams{Name: channel.Name, Color: channel.Color, Objects: channel.Objects}
}
