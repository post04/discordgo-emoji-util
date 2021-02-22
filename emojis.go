package emutil

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	emojiRegex = regexp.MustCompile(`<(a|):[A-z 0-9]+:[0-9]+>`)
)

/*
FindEmoji - Find a specific emoji in an array of emojis
emojis: list of discordgo.Emoji
emojiName: string that represents the emoji name you're trying to find the list
castSensitive: if the emoji finding should be case sensitive or not
*/
func FindEmoji(emojis []*discordgo.Emoji, emojiName string, caseSensitive bool) *discordgo.Emoji {
	if caseSensitive {
		for _, em := range emojis {
			if em.Name == emojiName {
				return em
			}
		}
		return nil
	}
	for _, em := range emojis {
		if strings.ToLower(em.Name) == strings.ToLower(emojiName) {
			return em
		}
	}
	return nil
}

/*
EncodeEmojiByID - takes an emoji ID and returns the base64 encoded emoji image, used to add emojis. Returns empty string if something fails
emojiID: the id of the emoji in question
*/
func EncodeEmojiByID(emojiID string) string {
	return EncodeImageEmoji("https://cdn.discordapp.com/emojis/" + emojiID)
}

/*
EncodeImageEmoji - encodes an image link to base64, discord takes base64 to add an emoji so this will return what needs to be given to discord.
link: the url that the emoji resides at
*/
func EncodeImageEmoji(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString(body)
	return "data:image/png;base64," + encoded
}

// Emoji - contains the emojis id and name
type Emoji struct {
	ID   string
	Name string
}

/*
MatchEmojis - uses an emoji regex to match all emojis in a message, returns list of Emoji struct. Returns nil if something goes wrong or nothing is found!
messageContent: the string that we apply the regex to.
*/
func MatchEmojis(messageContent string) []*Emoji {
	var toReturn []*Emoji
	emojis := emojiRegex.FindAllString(messageContent, -1)
	if len(emojis) < 1 {
		return nil
	}
	for _, em := range emojis {
		p := strings.Split(em, ":")
		id := p[2]
		name := p[1]
		id = strings.ReplaceAll(id, ">", "")
		toReturn = append(toReturn, &Emoji{
			ID:   id,
			Name: name,
		})
	}
	return toReturn
}
