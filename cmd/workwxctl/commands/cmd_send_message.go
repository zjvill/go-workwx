package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/zjvill/go-workwx/v2"
)

func cmdSendMessage(c *cli.Context) error {
	cfg := mustGetConfig(c)
	isSafe := c.Bool(flagSafe)
	toUsers := c.StringSlice(flagToUser)
	toParties := c.StringSlice(flagToParty)
	toTags := c.StringSlice(flagToTag)
	toChat := c.String(flagToChat)
	content := c.Args().Get(0)
	msgtype := c.String(flagMessageType)

	mediaID := c.String(flagMediaID)
	thumbMediaID := c.String(flagThumbMediaID)
	description := c.String(flagDescription)
	title := c.String(flagTitle)
	author := c.String(flagAuthor)
	url := c.String(flagURL)
	picURL := c.String(flagPicURL)
	buttonText := c.String(flagButtonText)
	// sourceContentURL := c.String(flagSourceContentURL)
	digest := c.String(flagDigest)

	app := cfg.MakeWorkwxApp()

	recipient := workwx.Recipient{
		UserIDs:  toUsers,
		PartyIDs: toParties,
		TagIDs:   toTags,
		ChatID:   toChat,
	}

	if msgtype == "" {
		// default to text
		msgtype = string(workwx.MessageTypeText)
	}

	var err error
	switch msgtype {
	case string(workwx.MessageTypeText):
		err = app.SendTextMessage(&recipient, content, isSafe)
	case string(workwx.MessageTypeImage):
		err = app.SendImageMessage(&recipient, mediaID, isSafe)
	case string(workwx.MessageTypeVoice):
		err = app.SendVoiceMessage(&recipient, mediaID, isSafe)
	case string(workwx.MessageTypeVideo):
		err = app.SendVideoMessage(
			&recipient,
			mediaID,
			description,
			title,
			isSafe,
		)
	case "file":
		err = app.SendFileMessage(&recipient, mediaID, isSafe)
	case "textcard":
		err = app.SendTextCardMessage(
			&recipient,
			title,
			description,
			url,
			buttonText,
			isSafe,
		)
	case "news":
		err = app.SendNewsMessage(
			&recipient,
			[]workwx.Article{
				workwx.Article{
					Title:       title,
					Description: description,
					URL:         url,
					PicURL:      picURL,
					AppID:       "",
					PagePath:    "",
				},
			},
			isSafe,
		)
	case "mpnews":
		err = app.SendMPNewsMessage(
			&recipient,
			[]workwx.MPArticle{workwx.MPArticle{
				Title:            title,
				ThumbMediaID:     thumbMediaID,
				Author:           author,
				ContentSourceURL: content,
				Content:          content,
				Digest:           digest,
			}},
			isSafe,
		)
	default:
		fmt.Printf("unrecognized message type: %s\n", msgtype)
		panic("unrecognized message type")
	}

	return err
}
