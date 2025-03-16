# Gator 🐊

This project is called "Gator", you know, because aggreGATOR 🐊. Anyhow, it is a CLI tool built in Golang and PostgreSQL that allows users to:

* Add RSS feeds from across the internet to be collected;
* Store the collected posts in a PostgreSQL database;
* Follow and unfollow RSS feeds that other users have added;
* View summaries of the aggregated posts in the terminal, with a link to the full post;

RSS feeds are a way for websites to publish updates to their content.
You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

# How to run it

Clone the repository in your preferred location destination. Run the command ```go build -d ./gator``` and use `./gator` followed by the command you would like to run. These are the possible commands:
* User management:
    * `register <username>`: Registers a new user and logs into the newly created account; 
    * `login <username>`: Logs into an account previously added via the `register` command;
    * `users`: Returns the list of registered users;
* Aggregation of new feeds:
    * `agg`: fetches the RSS feeds, parses them, and prints the posts to the console... All in a long-running loop that should be kept in a background process in a separated terminale than the one used to operate Gator;
    * `addfeed <feed_name> <feed_url>`: Adds a new RSS feed and assigns it to the currently logged in person;
    * `feeds`: Shows the list of all the available feeds;
* Follows management:
    * `follow`: Adds an existing RSS feed to the currently logged in user;
    * `following`: Shows the list of all the Feeds being followed by a specific user;
    * `unfollow`: Removes an RSS feed from the currently logged in user;
* Red alert - Red button:
    * `reset`: Deletes and restarts with fresh databases; It shouldn't be used, so use it at your own risk;

# Project status

The project is at its initial stages and it will be evolving following the Blog Aggregator Course in BootDev.
It is a project built to make use of `Golang` with `SQL` (`PostgreSQL`) configured via `Goose` and translating the `SQL queries` to `Golang` via `sqlc`
It stores the currently logged person in a local
Further enhancement could be taken into consideration, however, the current focus is on completing the course. 
