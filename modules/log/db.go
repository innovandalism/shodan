package log

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"time"
)

// DBUpsertGuild adds or updates a guild in the database
func DBUpsertGuild(db *sql.DB, guild *discordgo.Guild) error {
	stmt, err := db.Prepare("INSERT INTO discord_guild (id, name) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET name = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(guild.ID, guild.Name)
	return err
}

// DBUpsertChannel adds or updates a channel in the database
func DBUpsertChannel(db *sql.DB, channel *discordgo.Channel) error {
	stmt, err := db.Prepare("INSERT INTO discord_channel (id, guild, name) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET guild = $2,name = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(channel.ID, channel.GuildID, channel.Name)
	return err
}

// DBUpsertUser adds or updates a user in the database
func DBUpsertUser(db *sql.DB, user *discordgo.User) error {
	stmt, err := db.Prepare("INSERT INTO discord_user (id, name, discriminator) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET name = $2,discriminator = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.ID, user.Username, user.Discriminator)
	return err
}

// DBUpsertMessage adds a message in the database
func DBUpsertMessage(db *sql.DB, message *discordgo.Message) error {
	stmt, err := db.Prepare("INSERT INTO discord_message (id, channel, userid, message, time) VALUES ($1, $2, $3, $4,$5)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(message.ID, message.ChannelID, message.Author.ID, message.Content, time.Now().UTC())
	return err
}
