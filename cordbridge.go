package cordbridge

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	Session *discordgo.Session
	GuildID string
}

func NewClient(token, guildID string) (*Client, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	
	return &Client{
		Session: dg,
		GuildID: guildID,
	}, nil
}

func (c *Client) Open() error {
	return c.Session.Open()
}

func (c *Client) Close() error {
	return c.Session.Close()
}

func (c *Client) CreateCategory(categoryName string) (*discordgo.Channel, error) {
	channels, err := c.Session.GuildChannels(c.GuildID)
	if err != nil {
		return nil, fmt.Errorf("error fetching channels: %w", err)
	}

	for _, ch := range channels {
		if ch.Name == categoryName && ch.Type == discordgo.ChannelTypeGuildCategory {
			return nil, fmt.Errorf("category already exists: %s", categoryName)
		}
	}

	category, err := c.Session.GuildChannelCreateComplex(c.GuildID, discordgo.GuildChannelCreateData{
		Name: categoryName,
		Type: discordgo.ChannelTypeGuildCategory,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating category: %w", err)
	}

	return category, nil
}

func (c *Client) CreateChannel(channelName, categoryID string) (*discordgo.Channel, error) {
	channel, err := c.Session.GuildChannelCreateComplex(c.GuildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating channel: %w", err)
	}
	return channel, nil
}

func (c *Client) FindChannelByName(channelName, categoryID string) (*discordgo.Channel, error) {
	channels, err := c.Session.GuildChannels(c.GuildID)
	if err != nil {
		return nil, err
	}

	for _, ch := range channels {
		if ch.Name == channelName && ch.ParentID == categoryID && ch.Type == discordgo.ChannelTypeGuildText {
			return ch, nil
		}
	}
	return nil, nil
}

func (c *Client) SendMessageToChannel(channelID, message string) error {
	if channelID == "" {
		return errors.New("channelID is empty")
	}
	_, err := c.Session.ChannelMessageSend(channelID, message)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}

func (c *Client) DeleteMessageByID(channelID, messageID string) error {
	if channelID == "" {
		return errors.New("channelID is empty")
	}
	return c.Session.ChannelMessageDelete(channelID, messageID)
}

func (c *Client) EditMessageByID(channelID, messageID, newContent string) error {
	if channelID == "" {
		return errors.New("channelID is empty")
	}
	_, err := c.Session.ChannelMessageEdit(channelID, messageID, newContent)
	if err != nil {
		return fmt.Errorf("error editing message: %w", err)
	}
	return nil
}

func (c *Client) ReadAllMessagesFromChannel(channelID string) ([]*discordgo.Message, error) {
	var allMessages []*discordgo.Message
	var lastID string

	for {
		msgs, err := c.Session.ChannelMessages(channelID, 100, lastID, "", "")
		if err != nil {
			return nil, err
		}
		if len(msgs) == 0 {
			break
		}

		allMessages = append(allMessages, msgs...)
		lastID = msgs[len(msgs)-1].ID

		time.Sleep(250 * time.Millisecond) // respect rate limits
	}

	return allMessages, nil
}

func (c *Client) ReadLastXMessagesFromChannel(channelID string, x int) ([]*discordgo.Message, error) {
	if x <= 0 {
		return nil, fmt.Errorf("invalid limit")
	}

	var allMessages []*discordgo.Message
	var lastID string
	remaining := x

	for remaining > 0 {
		fetchLimit := remaining
		if fetchLimit > 100 {
			fetchLimit = 100
		}

		msgs, err := c.Session.ChannelMessages(channelID, fetchLimit, lastID, "", "")
		if err != nil {
			return nil, err
		}
		if len(msgs) == 0 {
			break
		}

		allMessages = append(allMessages, msgs...)
		lastID = msgs[len(msgs)-1].ID
		remaining -= len(msgs)

		time.Sleep(250 * time.Millisecond) // respect rate limits
	}

	for i, j := 0, len(allMessages)-1; i < j; i, j = i+1, j-1 { //reverse
		allMessages[i], allMessages[j] = allMessages[j], allMessages[i]
	}

	return allMessages, nil
}