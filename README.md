# Git-Broski 

**Broski for your Git!**  
A CLI tool to perform various manual tasks with single commands

## Installation

### Prerequisites
- Go 1.22+ installed

### Setup
```bash
# Clone the repo
git clone https://github.com/gitbroskie/gitbroski.git
cd gitbroski

# Build
go build -o gitbroski ./cmd

# Create symlink for global access
sudo ln -s $(pwd)/gitbroski /usr/local/bin/gitbroski
```

Now you can run `gitbroski` from anywhere.

## Usage

### Open Remote Repository
Jump from terminal to your GitHub/GitLab repo page:
```bash
gitbroski open
```

### Auto-Generate .gitignore
Create a .gitignore file for your project:
```bash
gitbroski ignore <language>
```
Supported: `python` (more coming soon)

### Empty Commit
Trigger CI/CD without code changes:
```bash
gitbroski empty commit <message>
```

### MR/PR Manager
Save and manage your merge requests:
```bash
# Save an MR/PR
gitbroski mr save <url>

# List saved MRs (interactive)
gitbroski mr list
```

In the list view:
- `↑/↓` navigate
- `enter` open in browser
- `d` delete
- `h` auth help
- `q` quit

## Contributing
1. Fork and clone
2. `go mod tidy`
3. `go build -o gitbroski ./cmd`
4. Make changes and submit a PR

## License
MIT
