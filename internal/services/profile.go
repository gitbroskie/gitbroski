package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gitbroski/utils/logger"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

// ── Data structs ─────────────────────────────────────────────────────────────

type ProfileEntry struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Host     string `json:"host"`
}

type ProfileStore struct {
	Active   string                  `json:"active"`
	Profiles map[string]ProfileEntry `json:"profiles"`
}

// ── Config file path ──────────────────────────────────────────────────────────

func profileFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	dir := filepath.Join(home, ".gitbroski")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	return filepath.Join(dir, "profiles.json"), nil
}

// ── Read / Write ──────────────────────────────────────────────────────────────

func loadProfiles() (*ProfileStore, error) {
	path, err := profileFilePath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &ProfileStore{Profiles: make(map[string]ProfileEntry)}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read profiles file: %w", err)
	}
	var store ProfileStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("profiles file is corrupted: %w", err)
	}
	if store.Profiles == nil {
		store.Profiles = make(map[string]ProfileEntry)
	}
	return &store, nil
}

func saveProfiles(store *ProfileStore) error {
	path, err := profileFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("could not serialize profiles: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("could not save profiles file: %w", err)
	}
	return nil
}

// ── Entry point ───────────────────────────────────────────────────────────────

func Profile(args ...string) {
	if len(args) == 0 {
		printProfileUsage()
		return
	}

	switch args[0] {
	case "add":
		profileAdd()
	case "list", "ls":
		profileList()
	case "use":
		profileUse(args[1:])
	case "remove", "rm":
		profileRemove(args[1:])
	default:
		logger.Error("Unknown subcommand: " + args[0])
		printProfileUsage()
	}
}

// ── Subcommands ───────────────────────────────────────────────────────────────

// profile add — fully interactive
func profileAdd() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("Adding a new GitHub profile")
	fmt.Println("───────────────────────────")

	name := prompt(reader, "Profile name (e.g. work, personal): ", false)
	if name == "" {
		logger.Error("Profile name cannot be empty.")
		return
	}

	store, err := loadProfiles()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if _, exists := store.Profiles[name]; exists {
		logger.Error(fmt.Sprintf("Profile '%s' already exists. Remove it first with: profile remove %s", name, name))
		return
	}

	username := prompt(reader, "GitHub username: ", false)
	if username == "" {
		logger.Error("Username cannot be empty.")
		return
	}

	email := prompt(reader, "GitHub email: ", false)
	if email == "" {
		logger.Error("Email cannot be empty.")
		return
	}

	token := prompt(reader, "Personal Access Token (PAT): ", true)
	if token == "" {
		logger.Error("Token cannot be empty.")
		return
	}

	host := prompt(reader, "GitHub host [github.com]: ", false)
	if host == "" {
		host = "github.com"
	}

	store.Profiles[name] = ProfileEntry{
		Username: username,
		Email:    email,
		Token:    token,
		Host:     host,
	}

	// Auto-activate if first profile
	if store.Active == "" {
		store.Active = name
	}

	if err := saveProfiles(store); err != nil {
		logger.Error(err.Error())
		return
	}

	fmt.Println()
	logger.Text(fmt.Sprintf("✓ Profile '%s' saved (%s <%s> @ %s)", name, username, email, host))
	if store.Active == name {
		logger.Text("  Set as active profile.")
	}
}

// profile list
func profileList() {
	store, err := loadProfiles()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if len(store.Profiles) == 0 {
		logger.Text("No profiles found. Run 'profile add' to create one.")
		return
	}

	fmt.Println()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "  NAME\tUSERNAME\tEMAIL\tHOST")
	fmt.Fprintln(w, "  ----\t--------\t-----\t----")
	for name, p := range store.Profiles {
		prefix := "  "
		if name == store.Active {
			prefix = "* "
		}
		fmt.Fprintf(w, "%s%s\t%s\t%s\t%s\n", prefix, name, p.Username, p.Email, p.Host)
	}
	w.Flush()
	fmt.Println()
}

// profile use <name>
func profileUse(args []string) {
	if len(args) == 0 {
		logger.Error("Usage: profile use <name>")
		return
	}
	name := args[0]
	store, err := loadProfiles()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if _, exists := store.Profiles[name]; !exists {
		logger.Error(fmt.Sprintf("Profile '%s' not found. Run 'profile list' to see available profiles.", name))
		return
	}
	store.Active = name
	if err := saveProfiles(store); err != nil {
		logger.Error(err.Error())
		return
	}
	p := store.Profiles[name]
	logger.Text(fmt.Sprintf("✓ Switched to profile '%s' (%s <%s> @ %s)", name, p.Username, p.Email, p.Host))
}

// profile remove <name>
func profileRemove(args []string) {
	if len(args) == 0 {
		logger.Error("Usage: profile remove <name>")
		return
	}
	name := args[0]
	store, err := loadProfiles()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if _, exists := store.Profiles[name]; !exists {
		logger.Error(fmt.Sprintf("Profile '%s' not found.", name))
		return
	}
	delete(store.Profiles, name)

	if store.Active == name {
		store.Active = ""
		for next := range store.Profiles {
			store.Active = next
			break
		}
		if store.Active != "" {
			logger.Text(fmt.Sprintf("  Active profile switched to '%s'.", store.Active))
		} else {
			logger.Text("  No profiles remaining.")
		}
	}

	if err := saveProfiles(store); err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Text(fmt.Sprintf("✓ Profile '%s' removed.", name))
}

// ── Exported helper ───────────────────────────────────────────────────────────

// GetActiveProfile is used by other services (clone, push, etc.) to get the current token.
func GetActiveProfile() (*ProfileEntry, string, error) {
	store, err := loadProfiles()
	if err != nil {
		return nil, "", err
	}
	if store.Active == "" {
		return nil, "", fmt.Errorf("no active profile — run 'profile add' to get started")
	}
	p, ok := store.Profiles[store.Active]
	if !ok {
		return nil, "", fmt.Errorf("active profile '%s' not found in store", store.Active)
	}
	return &p, store.Active, nil
}

// ── Internal helpers ──────────────────────────────────────────────────────────

// prompt prints a label and reads input from stdin.
// If hidden=true, terminal echo is suppressed (used for PAT input).
func prompt(reader *bufio.Reader, label string, hidden bool) string {
	fmt.Print(label)
	if hidden {
		fmt.Print("\033[8m") // ANSI: hide text
	}
	input, _ := reader.ReadString('\n')
	if hidden {
		fmt.Print("\033[28m") // ANSI: show text
		fmt.Println()
	}
	return strings.TrimSpace(input)
}

func printProfileUsage() {
	logger.Text("Usage:")
	logger.Text("  profile add            — interactive setup")
	logger.Text("  profile list           — show all profiles")
	logger.Text("  profile use <name>     — switch active profile")
	logger.Text("  profile remove <name>  — delete a profile")
}
