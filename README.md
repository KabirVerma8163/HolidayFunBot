# HolidayFunBot

A simple golang program for a discord bot that makes Secret Santa surprises, updates coming later, just a fun project for now.
Link to invite the bot [here!](https://discord.com/api/oauth2/authorize?client_id=1186890517817085983&permissions=137439415360&scope=bot)

Made with discordgo, great library found [here](https://github.com/bwmarrin/discordgo/tree/master/examples) with lots of great go samples for learning. Once you get the hang of it, really easy to make more stuff with simple code completion suggestions. Was a fun project, might make another bot in go later. 


## Commands

1. !SecretSanta (case insensitive):  
usage: `!SecretSanta <Event Name>`  
Sends a message to everyone in the server, telling them who their secret Santa is, doesn't support more than 1000 server members (limitations with the discordgo library)

