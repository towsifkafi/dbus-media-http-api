package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/gorilla/mux"
)

func getPlaying(w http.ResponseWriter, r *http.Request) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		respondWithError(w, "Failed to connect to session bus")
		return
	}
	defer conn.Close()

	var list []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&list)
	if err != nil {
		respondWithError(w, "Failed to list DBus names")
		return
	}

	results := make([]map[string]interface{}, 0)

	for _, name := range list {
		if strings.HasPrefix(name, "org.mpris.MediaPlayer2.") {
			mediaPlayerInfo, err := fetchMediaPlayerInfo(conn, name)
			if err == nil {
				results = append(results, mediaPlayerInfo)
			}
			// note: If err != nil, it'll simply skip the player.
		}
	}

	response := map[string]interface{}{
		"error":  false,
		"result": results,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, "Failed to marshal metadata to JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func fetchMediaPlayerInfo(conn *dbus.Conn, name string) (map[string]interface{}, error) {
	obj := conn.Object(name, "/org/mpris/MediaPlayer2")

	metadataJSON := make(map[string]interface{})

	if metadata, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata"); err == nil {
		var metadataMap map[string]dbus.Variant
		if err := metadata.Store(&metadataMap); err == nil {
			for k, v := range metadataMap {
				if k == "mpris:artUrl" {
					if artUrl, ok := v.Value().(string); ok {
						base64Data, err := fetchBase64EncodedArt(artUrl)
						if err == nil {
							metadataJSON[k] = base64Data
						}
					}
				} else {
					metadataJSON[k] = v.Value()
				}
			}
			metadataJSON["service"] = name
		}
	}

	addPropertyIfExists(obj, metadataJSON, "org.mpris.MediaPlayer2.Player.PlaybackStatus", "playback_status")
	addPropertyIfExists(obj, metadataJSON, "org.mpris.MediaPlayer2.Player.LoopStatus", "loop_status")
	addPropertyIfExists(obj, metadataJSON, "org.mpris.MediaPlayer2.Player.Shuffle", "shuffle")
	addPropertyIfExists(obj, metadataJSON, "org.mpris.MediaPlayer2.Player.Volume", "volume")
	addPropertyIfExists(obj, metadataJSON, "org.mpris.MediaPlayer2.Player.Position", "position")

	return metadataJSON, nil
}

func addPropertyIfExists(obj dbus.BusObject, metadataJSON map[string]interface{}, propertyName string, jsonKey string) {
	if prop, err := obj.GetProperty(propertyName); err == nil {
		metadataJSON[jsonKey] = prop.Value()
	}
}

func controlMediaPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]
	action := vars["action"]

	if !strings.HasPrefix(serviceName, "org.mpris.MediaPlayer2.") {
		respondWithError(w, "Invalid service name")
		return
	}

	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		respondWithError(w, "Failed to connect to session bus")
		return
	}
	defer conn.Close()

	obj := conn.Object(serviceName, "/org/mpris/MediaPlayer2")

	if strings.ToLower(action) == "playpause" {
		action = "PlayPause"
	}

	method := fmt.Sprintf("org.mpris.MediaPlayer2.Player.%s", strings.Title(action))

	call := obj.Call(method, 0)
	if call.Err != nil {
		fmt.Println(call.Err)
		respondWithError(w, "Failed to execute "+action)
		return
	}

	response := map[string]interface{}{
		"error":   false,
		"message": action + " executed successfully",
	}
	jsonData, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
