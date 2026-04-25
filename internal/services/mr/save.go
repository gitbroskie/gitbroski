package mr

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"gitbroski/utils/logger"
)

// Save saves a new MR URL to the store.
func Save(url string) {
	if url == "" {
		logger.Error("No URL provided")
		return
	}

	store, err := LoadStore()
	if err != nil {
		logger.Error("Failed to load MR store: " + err.Error())
		return
	}

	// Check if URL already exists
	for _, m := range store.MRs {
		if m.URL == url {
			logger.Text("MR already saved: " + url)
			return
		}
	}

	reader := bufio.NewReader(os.Stdin)

	// Prompt for title
	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	// Prompt for service
	fmt.Print("Service: ")
	service, _ := reader.ReadString('\n')
	service = strings.TrimSpace(service)

	newMR := MR{
		URL:     url,
		Title:   title,
		Service: service,
		SavedAt: time.Now(),
	}

	store.MRs = append(store.MRs, newMR)

	if err := SaveStore(store); err != nil {
		logger.Error("Failed to save MR: " + err.Error())
		return
	}

	logger.Text("Saved MR: " + url)
}
