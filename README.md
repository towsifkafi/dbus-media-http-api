# DBus Mris MediaPlayer HTTP API

This HTTP API connects to the host's DBus to retrieve current media information. It exposes an endpoint where this information can be accessed, and various methods can be executed.

## Endpoints

<details open>
    <summary><code>GET</code> <code><b>/get</b></code> <code>Returns a list of all available Mris MediaPlayer2 objects.</code></summary>

#### Response
```json
{
  "error": false,
  "result": [
    {
      "service": "org.mpris.MediaPlayer2.APPLICATION",
      "mpris:artUrl": "BASE64 DATA",
      "mpris:length": 1230,
      "mpris:trackid": "TRACK ID",
      "playback_status": "Playing",
      "position": 1230,
      "volume": 1,
      "xesam:album": "",
      "xesam:artist": [
        "ARTIST"
      ],
      "xesam:title": "TITLE"
    }
  ]
}
```

</details>

<details open>
    <summary><code>GET</code> <code><b>/{action}/{service}</b></code> <code>Triggers a specified action in the provided service.</code></summary>

#### Valid actions: 
* `playpause`
* `play`
* `pause`
* `stop`
* `next`
* `previous`

##### Response
```json
{"error": false, "message": "ACTION executed successfully"}
```

</details>

> [!NOTE]
> If you have `AUTH=true` and `AUTH_KEY=KEY` set in your .env file, you will need to add an `Authorization` header with the authentication key to your requests.


## Usage
To run this HTTP API, you need a Linux machine with DBus installed. This project integrates with the DBus Mris Media Player for media player interaction. You will also need to compile this project.

To compile and run this project, follow these commands:
```bash
git clone https://github.com/Towsif12/dbus-media-http-api.git
cd dbus-media-http-api

go build -o dbus-media-api ./src
chmod +x dbus-media-api

./dbus-media-api
``` 

Optionally, you can use a `.env` file to configure the port, authentication, and authentication key:
```
PORT=10004
AUTH=true
AUTH_KEY=key123
```

## Development
The development process is similar to the usage instructions. Start by cloning the repository, and then you can modify the code in the `src/` directory.

To run the project during development, use the following command:
```sh
go run ./src
```

Feel free to submit a Pull Request (PR) or post an issue if you encounter any bugs or errors.

## Sources
These links were helpful resources throughout the development of this project:
- https://www.baeldung.com/linux/dbus
- https://github.com/godbus/dbus/
- https://github.com/joho/godotenv
- https://tutorialedge.net/golang/creating-restful-api-with-golang/
- https://specifications.freedesktop.org/mpris-spec/2.2/Player_Interface.html
