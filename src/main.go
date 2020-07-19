package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const PostPerPage = 100

type Plugin struct {
	plugin.MattermostPlugin
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	command_args := strings.Fields(args.Command)
	trigger := strings.TrimPrefix(command_args[0], "/")

	switch trigger {
	case "channel_id":
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Channel id: " + args.ChannelId,
		}, nil
	case "user_id":
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "User id: " + args.UserId,
		}, nil
	case "team_id":
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Team id: " + args.TeamId,
		}, nil
	case "delete_posts":
		if !p.API.HasPermissionToChannel(args.UserId, args.ChannelId, model.PERMISSION_DELETE_OTHERS_POSTS) {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "Permission denied",
			}, nil
		}

		if command_args[1] != args.ChannelId {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "First argument must be channel id. First argument is " + command_args[1],
			}, nil
		}
		n, err := strconv.ParseFloat(command_args[2], 64)
		if err != nil {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "Second arguments must be a number",
			}, nil
		}

		cutoff := time.Now().Unix()*1000 - int64(n*24*60*60*1000)

		deleted := 0
		// Mattermost returns page in reverse chronological order. We first
		// find the first page with an old enough post.
		page := 0
		for {
			posts, err := p.API.GetPostsForChannel(args.ChannelId, page+1, PostPerPage)
			if err != nil {
				return &model.CommandResponse{
					ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
					Text:         "Failed to retrieve posts",
				}, nil
			}
			if len(posts.Order) == 0 {
				break
			}
			if posts.Posts[posts.Order[len(posts.Order)-1]].CreateAt < cutoff {
				break
			}
			page += 1
		}

		// Now page is the first page with a post we have to delete. First
		// delete all posts in subsequent pages
		for {
			posts, err := p.API.GetPostsForChannel(args.ChannelId, page+1, PostPerPage)
			if err != nil {
				return &model.CommandResponse{
					ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
					Text:         "Failed to retrieve posts",
				}, nil
			}
			if len(posts.Order) == 0 {
				break
			}
			for _, pid := range posts.Order {
				if e := p.API.DeletePost(pid); e != nil {
					return &model.CommandResponse{
						ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
						Text:         "Error deleting post " + pid,
					}, nil
				}
				deleted += 1
			}
		}

		// Now delete the rest in the final page
		posts, err := p.API.GetPostsForChannel(args.ChannelId, page, PostPerPage)
		for i := len(posts.Order) - 1; i >= 0; i-- {
			pid := posts.Order[i]
			if posts.Posts[pid].CreateAt < cutoff {
				if e := p.API.DeletePost(pid); e != nil {
					return &model.CommandResponse{
						ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
						Text:         "Error deleting post " + pid,
					}, nil
				}
				deleted += 1
			} else {
				break
			}
		}

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         strconv.Itoa(deleted) + " posts older than " + command_args[2] + " days deleted",
		}, nil
	}
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Unknown command: " + args.Command,
	}, nil
}

func (p *Plugin) OnActivate() error {
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "channel_id",
		AutoComplete:     true,
		AutoCompleteDesc: "Displays channel id",
	}); err != nil {
		return err
	}
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "user_id",
		AutoComplete:     true,
		AutoCompleteDesc: "Displays user id",
	}); err != nil {
		return err
	}
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "team_id",
		AutoComplete:     true,
		AutoCompleteDesc: "Displays team id",
	}); err != nil {
		return err
	}
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "delete_posts",
		AutoComplete:     true,
		AutoCompleteDesc: "Usage: /delete_posts <channel_id> <n>; This command deletes all posts more than n days old in this channel. The first argument is the channel_id obtained from /channel_id.",
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	plugin.ClientMain(&Plugin{})
}
