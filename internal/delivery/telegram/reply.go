package telegram

import (
	"errors"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

// reply keeps a single bot message per chat. A button press edits the message
// it came from; a command sends a new message, then deletes the command and
// the previous bot reply.
// Deletion errors are ignored: messages older than 48h or missing admin
// rights in groups must not break the response.
func (s *Server) reply(c tele.Context, text string, markup *tele.ReplyMarkup) error {
	if c.Callback() != nil {
		return s.editReply(c, text, markup)
	}

	msg, err := s.bot.Send(c.Recipient(), text, markup, tele.ModeHTML)
	if err != nil {
		return err
	}

	if c.Message() != nil {
		_ = c.Delete()
	}

	s.rememberReply(msg)
	return nil
}

func (s *Server) editReply(c tele.Context, text string, markup *tele.ReplyMarkup) error {
	_ = c.Respond()

	// Edited messages accept only inline keyboards; without one the inline
	// buttons are removed.
	opts := []interface{}{tele.ModeHTML}
	if markup != nil && len(markup.InlineKeyboard) > 0 {
		opts = append(opts, markup)
	}

	err := c.Edit(text, opts...)
	if err != nil && !errors.Is(err, tele.ErrSameMessageContent) && !errors.Is(err, tele.ErrMessageNotModified) {
		return err
	}

	if msg := c.Callback().Message; msg != nil {
		s.rememberReply(msg)
	}
	return nil
}

// rememberReply stores msg as the chat's latest reply and deletes the
// previous one if it is a different message.
func (s *Server) rememberReply(msg *tele.Message) {
	current := tele.StoredMessage{
		MessageID: strconv.Itoa(msg.ID),
		ChatID:    msg.Chat.ID,
	}
	prev, loaded := s.lastReplies.Swap(msg.Chat.ID, current)
	if loaded && prev.(tele.StoredMessage) != current {
		_ = s.bot.Delete(prev.(tele.StoredMessage))
	}
}
