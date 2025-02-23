#  Lazyface

Lazyface is a Terminal User Interface (TUI) for Hugging Face, allowing users to execute Hugging Face commands through an intuitive and visually appealing interface. 🎨

## ✨ Features

- 🔍 Browse and search Hugging Face models and datasets from the terminal.
- 📥 Execute commands such as downloading models, pushing models to the Hub, and more.
- 🎮 Simplified navigation with keyboard shortcuts.
- 🎨 A sleek and modern interface using Bubble Tea and Lip Gloss.

## ⚙️ Requirements

Before installing lazyface, ensure you have the following dependencies installed:

-  [Go](https://go.dev/) (latest version recommended)
-  [Bubble Tea](https://github.com/charmbracelet/bubbletea) (TUI framework)
-  [Lip Gloss](https://github.com/charmbracelet/lipgloss) (for styling)
- [Hugging Face CLI](https://huggingface.co/docs/huggingface_hub/main/en/) 🤗

Refer to the [📖 Hugging Face CLI documentation](https://huggingface.co/docs/huggingface_hub/main/en/guides/cli) for instructions on how to install Hugging Face CLI.

## 📦 Installation

Clone the repository and build the project:

```sh
# Clone the repository
git clone https://github.com/yourusername/lazyface.git
cd lazyface

# Install dependencies
go mod tidy

# Build the project
go build -o lazyface
```

Alternatively, install it using:

```sh
go install github.com/yourusername/lazyface@latest
```

You can also download pre-built binaries from the [📥 GitHub Releases](https://github.com/0xRavenspar/Lazyface/releases) page.

## 🚀 Usage

Run the application:

```sh
./lazyface
```

### ⌨️ Keyboard Shortcuts

- 🔼/🔽 `Up/Down` - Navigate through options
- ✅ `Enter` - Select an option
- ❌ `q` - Quit the application

## 🤝 Contributing

We welcome contributions! 🎉 Feel free to open issues and pull requests.

1. 🍴 Fork the repository.
2. 🌿 Create a new branch: `git checkout -b feature-branch`
3. 📝 Commit your changes: `git commit -m "Add new feature"`
4. 📤 Push to the branch: `git push origin feature-branch`
5. 🔄 Open a Pull Request.

## 📜 License

lazyface is licensed under the [MIT License](https://github.com/0xRavenspar/Lazyface/blob/main/LICENSE).
