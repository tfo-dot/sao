package discord

import (
	"sao/types"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type MessageBuffer struct {
	sync.Mutex
	messages       []types.DiscordMessageStruct
	lastUpdateTime time.Time
}

type MessageBufferManager struct {
	sync.Mutex
	buffers     map[string]*MessageBuffer
	flushPeriod time.Duration
}

func NewMessageBufferManager(flushPeriod time.Duration) *MessageBufferManager {
	manager := &MessageBufferManager{
		buffers:     make(map[string]*MessageBuffer),
		flushPeriod: flushPeriod,
	}

	go manager.periodicFlush()

	return manager
}

func (m *MessageBufferManager) Add(msg types.DiscordMessageStruct) {
	m.Lock()
	buffer, exists := m.buffers[msg.ChannelID]
	if !exists {
		buffer = &MessageBuffer{
			messages:       []types.DiscordMessageStruct{},
			lastUpdateTime: time.Now(),
		}
		m.buffers[msg.ChannelID] = buffer
	}
	m.Unlock()

	buffer.Lock()
	buffer.messages = append(buffer.messages, msg)
	buffer.lastUpdateTime = time.Now()
	buffer.Unlock()
}

func (m *MessageBufferManager) Flush(channelID string) {
	m.Lock()

	buffer, exists := m.buffers[channelID]

	if !exists {
		m.Unlock()
		return
	}

	delete(m.buffers, channelID)

	m.Unlock()

	buffer.Lock()
	defer buffer.Unlock()

	if len(buffer.messages) == 0 {
		return
	}

	isDM := buffer.messages[0].DM

	snowflakeID := snowflake.MustParse(channelID)

	if isDM {
		ch, err := (*Client).Rest().CreateDMChannel(snowflakeID)
		if err != nil {
			return
		}
		snowflakeID = ch.ID()
	}

	if len(buffer.messages) == 1 {
		(*Client).Rest().CreateMessage(snowflakeID, buffer.messages[0].MessageContent)

		return
	}

	var combinedContent discord.MessageCreate

	combinedContent = buffer.messages[0].MessageContent

	for i := 1; i < len(buffer.messages); i++ {
		msg := buffer.messages[i].MessageContent

		if len(msg.Embeds) > 0 {
			combinedContent.Embeds = append(combinedContent.Embeds, msg.Embeds...)
		}

		if len(msg.Components) > 0 {
			combinedContent.Components = append(combinedContent.Components, msg.Components...)
		}

		if msg.Content != "" && combinedContent.Content == "" {
			combinedContent.Content = msg.Content
		} else if msg.Content != "" && combinedContent.Content != "" {
			combinedContent.Content += "\n\n" + msg.Content
		}
	}

	(*Client).Rest().CreateMessage(snowflakeID, combinedContent)
}

func (m *MessageBufferManager) FlushAll() {
	m.Lock()

	channelsToFlush := make([]string, 0, len(m.buffers))

	for channelID := range m.buffers {
		channelsToFlush = append(channelsToFlush, channelID)
	}

	m.Unlock()

	for _, channelID := range channelsToFlush {
		m.Flush(channelID)
	}
}

func (m *MessageBufferManager) periodicFlush() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		channelsToFlush := make([]string, 0)

		m.Lock()

		for channelID, buffer := range m.buffers {
			buffer.Lock()
			if now.Sub(buffer.lastUpdateTime) >= m.flushPeriod {
				channelsToFlush = append(channelsToFlush, channelID)
			}
			buffer.Unlock()
		}

		m.Unlock()

		for _, channelID := range channelsToFlush {
			m.Flush(channelID)
		}
	}
}

func worldMessageListener() {
	bufferManager := NewMessageBufferManager(500 * time.Millisecond)

	for {
		msg, ok := <-World.DiscordChannel

		if !ok {
			return
		}

		switch msg.GetEvent() {
		case types.MSG_SEND:			
			bufferManager.Add(msg.GetData().(types.DiscordMessageStruct))
		case types.MSG_CHOICE:
			bufferManager.FlushAll()

			Choices = append(Choices, msg.GetData().(types.DiscordChoice))
		}
	}
}
