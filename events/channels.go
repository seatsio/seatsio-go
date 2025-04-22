package events

import (
	"context"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v11/shared"
)

type Channels struct {
	Client *req.Client
}

type CreateChannelParams struct {
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Index   int      `json:"index,omitempty"`
	Objects []string `json:"objects,omitempty"`
}

type CreateChannelParamsOption func(Params *CreateChannelParams)

type UpdateChannelParams struct {
	Name    string   `json:"name,omitempty"`
	Color   string   `json:"color,omitempty"`
	Objects []string `json:"objects,omitempty"`
}

type changeChannelObjectsRequest struct {
	Objects []string `json:"objects"`
}

type replaceChannelsRequest struct {
	Channels []CreateChannelParams `json:"channels"`
}

func (channels *Channels) Create(context context.Context, eventKey string, params ...*CreateChannelParams) error {
	result, err := channels.Client.R().
		SetContext(context).
		SetBody(params).
		SetPathParam("key", eventKey).
		Post("/events/{key}/channels")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) Update(context context.Context, eventKey string, channelKey string, params UpdateChannelParams) error {
	result, err := channels.Client.R().
		SetContext(context).
		SetBody(params).
		SetPathParam("eventKey", eventKey).
		SetPathParam("channelKey", channelKey).
		Post("/events/{eventKey}/channels/{channelKey}")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) Delete(context context.Context, eventKey string, channelKey string) error {
	result, err := channels.Client.R().
		SetContext(context).
		SetPathParam("eventKey", eventKey).
		SetPathParam("channelKey", channelKey).
		Delete("/events/{eventKey}/channels/{channelKey}")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) AddObjects(context context.Context, eventKey string, channelKey string, objects []string) error {
	result, err := channels.Client.R().
		SetContext(context).
		SetBody(changeChannelObjectsRequest{objects}).
		SetPathParam("eventKey", eventKey).
		SetPathParam("channelKey", channelKey).
		Post("/events/{eventKey}/channels/{channelKey}/objects")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) RemoveObjects(context context.Context, eventKey string, channelKey string, objects []string) error {
	result, err := channels.Client.R().
		SetContext(context).
		SetBody(changeChannelObjectsRequest{objects}).
		SetPathParam("eventKey", eventKey).
		SetPathParam("channelKey", channelKey).
		Delete("/events/{eventKey}/channels/{channelKey}/objects")
	return shared.AssertOkNoBody(result, err)
}

func (channels *Channels) Replace(context context.Context, eventKey string, newChannels ...CreateChannelParams) error {
	result, err := channels.Client.R().
		SetContext(context).
		SetBody(replaceChannelsRequest{newChannels}).
		SetPathParam("eventKey", eventKey).
		Post("/events/{eventKey}/channels/replace")
	return shared.AssertOkNoBody(result, err)
}
