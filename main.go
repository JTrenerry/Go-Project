package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "sync"
    "time"
    "os/signal"
    "strconv"
)

const ( 
    USAGE = `Usage: go run main.go [--url URL | --hiphop | --synth | --piano | --ambient ] [--volume VOLUME | --mute ]`
    HIPHOP = "https://www.youtube.com/watch?v=jfKfPfyJRdk"
    SYNTH = "https://www.youtube.com/watch?v=4xDzrJKXOOY"
    PIANO = "https://www.youtube.com/watch?v=4oStw0r33so"
    AMBIENT = "https://www.youtube.com/watch?v=S_MOd40zlYU"
)

var mpvProcess *os.Process

// Displays some art to chill to
func display(wg *sync.WaitGroup) {
    defer wg.Done()
    times := 0
    for true {
        time.Sleep(1 * time.Second)
        fmt.Print("\033[H\033[2J")
        fmt.Println(FISH[times%2])
        fmt.Println("Time spent: " + strconv.Itoa(times))
        times++
    }
}

// Play the audio stream from the given URL
func play(url string, volume string, wg *sync.WaitGroup) {
    defer wg.Done()

    // Use yt-dlp to get the best audio stream URL
    cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-g", url)
    out, err := cmd.Output()

    // Check if the audio stream URL failed to get
    if err != nil {
        log.Fatalf("Failed to get audio URL: %v", err)
    }
    
    audioURL := string(out)
    audioURL = audioURL[:len(audioURL)-1] // Remove the trailing newline

    // Use mpv to play the audio stream

    playCmd := exec.Command("mpv", audioURL, "--volume=" + volume)
	
    mpvProcess = playCmd.Process
    //playCmd.Stdout = os.Stdout
    //playCmd.Stderr = os.Stderr

    err = playCmd.Run()
}

func main() {
    var volume = "75" // Default volume
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

    var wg sync.WaitGroup

    wg.Add(2)

	sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, os.Kill)

    if !mute {
        go play(URL, volume, &wg)
    }

    go func() {
		// Wait for a signal
		sig := <-sigs
		log.Printf("Received signal: %v. Forwarding to MPV.", sig)

		// Forward the signal to the MPV process
		if mpvProcess != nil {
			if err := mpvProcess.Signal(sig); err != nil {
				log.Printf("Failed to forward signal to MPV: %v", err)
			}
		} else {
			log.Printf("MPV process not found.")
		}
        os.Exit(0)
	}()

    go display(&wg)

    wg.Wait()
}
