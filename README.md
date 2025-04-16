# Blog Aggre(Gator)
This is a Blog aggregator or, for short, Gator. It uses Postgres and Go to run the program, so they must be installed first. 

### Installation
The program can be easily installed using `go install`.

To run the program, you will need to create a SQL dabatase and, most importantly, create a config file (I named it `.gatorconfig.json`) and store it in the home directory. 
The structure is as follows:
```json
{"db_url":"<your_connection_string>","current_user_name":"Lumian"}
```
You may look how to get your connection string in online tutorials. 

### How to use it
The program is run via the CLI and currently supports the following commands: 
- `register`: Registers a user. Expects an argument called username, e.g `register Lumian`
- `login`: Logins a user. Expects an argument, the username.
- `reset`: Deletes all registered users from the database.
- `users`: Prints all users registered in the database.
- `agg`: Expects a single argument in time-like format, e.g `1s`, `1m`, `1h`. It fetches all feeds from a given RSS.
- `addfeed`: Expects two arguments, a name and a url. Adds a feed to the database with the given name.
- `follow`: Expects a single argument, a url. Simulates the current user following said feed.
- `following`: Shows all feeds that the current user is following.
- `unfollow`: Self-explanatory.
- `browse`: Shows the number of posts that are stored in the database (each feed may contain many posts). Has `limit` as an optional argument which limits the number of posts shown, defaults to `2`.
