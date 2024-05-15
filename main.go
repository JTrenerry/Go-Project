package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
)

const USAGE = `
 Usage: go run main.go [--url URL | --hiphop | --synth | --piano | --ambient ] [--volume VOLUME | --mute ]
`

const HIPHOP = "https://www.youtube.com/watch?v=jfKfPfyJRdk"
const SYNTH = "https://www.youtube.com/watch?v=4xDzrJKXOOY"
const PIANO = "https://www.youtube.com/watch?v=4oStw0r33so"
const AMBIENT = "https://www.youtube.com/watch?v=S_MOd40zlYU"

func play(url string, volume string) {

    // Use yt-dlp to get the best audio stream URL
    cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-g", url)
    out, err := cmd.Output()

    // Check if the audio stream URL failed to get
    if err != nil {
        log.Fatalf("Failed to get audio URL: %v", err)
    }

    log.Printf("Got audio URL: %s", url)
    
    audioURL := string(out)
    audioURL = audioURL[:len(audioURL)-1] // Remove the trailing newline

    // Use mpv to play the audio stream

    playCmd := exec.Command("mpv", audioURL, "--volume=" + volume)
    playCmd.Stdout = os.Stdout
    playCmd.Stderr = os.Stderr

    err = playCmd.Run()

    // Check if the audio stream failed to play
    if err != nil {
        log.Fatalf("Failed to play audio: %v", err)
    }
}

func main() {
    var volume = "100" // Default volume
    var URL = HIPHOP // Default to hiphop
    var mute = false

    // Parse command line arguments
    for i, arg := range os.Args {
        if i == 0 {
            continue
        }
        // Help
        if arg == "--help" {
            fmt.Println(USAGE)
            os.Exit(0)
        }
        // Volume Control
        if arg == "--volume" {
            if len(os.Args) < i+1 {
                log.Fatalf(USAGE)
            }
            volume = os.Args[i+1]
        }
        if arg == "--mute" {
            mute = true
        }

        // Video Control
        // In built video options
        if arg == "--hiphop" {
            URL = HIPHOP
        }
        if arg == "--synth" {
            URL = SYNTH
        }
        if arg == "--piano" {
            URL = PIANO
        }
        if arg == "--ambient" {
            URL = AMBIENT
        }

        // Custom video URL
        if arg == "--url" {
            if len(os.Args) < i+1 {
                log.Fatalf(USAGE)
            }
            URL = os.Args[i+1]
        }
    }

    if !mute {
        play(URL, volume)
    }
}

