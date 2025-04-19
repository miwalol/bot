# Miwa Discord Bot

This repository contains the code of the [Miwa Discord server](https://discord.gg/miwa) bot.

## Features

The bot is used to:
- Send a welcome message to joining members
- Thanks to the user when someone boosted the server
- Have stats of the website with the `/stats` command

## Project Structure

- `events/`: Contains the event handlers for the bot.
- `models/`: Contains the database models.
- `utils/`: Contains utility functions, such as the database connection.
- `main.go`: Self-explanatory, the main file of the bot.
- `server.go`: A local server used by the Miwa API and website, to get users and presences. 

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.