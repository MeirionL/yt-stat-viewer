# yt-stat-viewer

A full stack application where users can search the names of YouTube channels to have some statistics about those channels displayed to them. They can also authenticate their own Google account in order to allow the app to show some extra information with regards to their own channels recent livestreaming. This will register them as a user, making sure they don't have to log in again later, though they can logout whenever they want. It uses Golang for the backend and React for the frontend.

## Installation

Inside a Go module:

```bash
go get github.com/MeirionL/yt-stat-viewer
```

## How to start

Inside the project directory create a `.env` file with the following 5 values:

- A `PORT` address number
- A `DB_URL` for a PostgreSQL database in this format: "postgres://<postgres_username>:<postgres_username_password>@localhost:<server_port>/<database_name>?sslmode_disable"
- A `YOUTUBE_API_KEY` that can be obtained by following the instructions in this link:
- Both a `GOOGLE_CLIENT_ID` and a `GOOGLE_CLIENT_SECRET` that can be obtained by following the instructions in this link:

Start the backend server:

```bash
go build && ./yt-stat-viewer
```

Start the front end in a new terminal:

```bash
cd client

yarn dev
```

Now you can access the application by visitng `http://localhost:5173`

## Usage

Once you've opened the app, you can use the top search bar to search a YouTube channels name. The app will then add that channel to the list of displayed channel names, where you can click on a certain name to view its channel stats. The stats that get displayed are the channels total subscriber count, its total video count, and its total view count. If the channel has logged in, then it will also show the start time of the channels most recent livestream, if it's live currently, and the title of its most recently started stream.

## Handling authentication

To authenticate your own YouTube channel, make sure you're logged into the appropriate Google account before clicking the "Login" button. Afterwards, your accounts unique ID will be displayed in the url of the page in the format:

    localhost:5173/?id=<ID_VALUE>

You will need to enter the ID_VALUE into the search bar above the logout button before clicking it if you wish to withdraw authorization rights to your accounts recent stream data.

## Contributing

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
