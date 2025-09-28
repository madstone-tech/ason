#!/bin/bash

# ※ Ason Autocompletion Installation Script
# This script installs shell completion for Ason across different shells

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_NAME="ason"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}※ Ason Autocompletion Installer${NC}"
    echo "=================================="
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Check if ason binary exists
check_binary() {
    if ! command -v "$BINARY_NAME" &> /dev/null; then
        print_error "ason binary not found in PATH"
        print_info "Please install ason first or add it to your PATH"
        exit 1
    fi
    print_success "Found ason binary at $(which $BINARY_NAME)"
}

# Install bash completion
install_bash_completion() {
    print_info "Installing bash completion..."

    # Try system-wide completion directory
    if [[ -d /usr/share/bash-completion/completions ]]; then
        COMPLETION_DIR="/usr/share/bash-completion/completions"
    elif [[ -d /etc/bash_completion.d ]]; then
        COMPLETION_DIR="/etc/bash_completion.d"
    elif [[ -d /usr/local/share/bash-completion/completions ]]; then
        COMPLETION_DIR="/usr/local/share/bash-completion/completions"
    else
        # Fallback to user directory
        COMPLETION_DIR="$HOME/.local/share/bash-completion/completions"
        mkdir -p "$COMPLETION_DIR"
    fi

    if [[ -w "$COMPLETION_DIR" ]] || [[ "$COMPLETION_DIR" == "$HOME"* ]]; then
        $BINARY_NAME completion bash > "$COMPLETION_DIR/ason"
        print_success "Bash completion installed to $COMPLETION_DIR/ason"
    else
        print_warning "No write permission to $COMPLETION_DIR"
        print_info "Run with sudo or manually run: $BINARY_NAME completion bash > ~/.local/share/bash-completion/completions/ason"
        return 1
    fi
}

# Install zsh completion
install_zsh_completion() {
    print_info "Installing zsh completion..."

    # Check for oh-my-zsh first
    if [[ -d "$HOME/.oh-my-zsh" ]]; then
        # Oh My Zsh uses custom completions directory
        COMPLETION_DIR="$HOME/.oh-my-zsh/custom/plugins/ason"
        mkdir -p "$COMPLETION_DIR"

        # Create plugin structure for Oh My Zsh
        $BINARY_NAME completion zsh > "$COMPLETION_DIR/_ason"

        # Create plugin file
        cat > "$COMPLETION_DIR/ason.plugin.zsh" << 'EOF'
# Ason completion plugin for Oh My Zsh
# Add this plugin to your .zshrc plugins list: plugins=(... ason)

# The completion file _ason in this directory will be automatically loaded
EOF

        print_success "Oh My Zsh plugin created at $COMPLETION_DIR"
        print_info "Add 'ason' to your plugins list in ~/.zshrc: plugins=(... ason)"

    elif [[ -d "$HOME/.oh-my-zsh/completions" ]]; then
        # Alternative Oh My Zsh completions directory
        COMPLETION_DIR="$HOME/.oh-my-zsh/completions"
        mkdir -p "$COMPLETION_DIR"
        $BINARY_NAME completion zsh > "$COMPLETION_DIR/_ason"
        print_success "Oh My Zsh completion installed to $COMPLETION_DIR/_ason"

    else
        # Standard zsh completion directory
        COMPLETION_DIR="${ZDOTDIR:-$HOME}/.zfunc"
        mkdir -p "$COMPLETION_DIR"

        # Add to fpath if not already there
        if [[ -f "${ZDOTDIR:-$HOME}/.zshrc" ]]; then
            if ! grep -q "fpath=.*$COMPLETION_DIR" "${ZDOTDIR:-$HOME}/.zshrc"; then
                echo "" >> "${ZDOTDIR:-$HOME}/.zshrc"
                echo "# Ason completion" >> "${ZDOTDIR:-$HOME}/.zshrc"
                echo "fpath=($COMPLETION_DIR \$fpath)" >> "${ZDOTDIR:-$HOME}/.zshrc"
                echo "autoload -U compinit" >> "${ZDOTDIR:-$HOME}/.zshrc"
                echo "compinit" >> "${ZDOTDIR:-$HOME}/.zshrc"
                print_info "Added completion setup to .zshrc"
            fi
        fi

        $BINARY_NAME completion zsh > "$COMPLETION_DIR/_ason"
        print_success "Zsh completion installed to $COMPLETION_DIR/_ason"
    fi
}

# Install fish completion
install_fish_completion() {
    print_info "Installing fish completion..."

    COMPLETION_DIR="$HOME/.config/fish/completions"
    mkdir -p "$COMPLETION_DIR"

    $BINARY_NAME completion fish > "$COMPLETION_DIR/ason.fish"
    print_success "Fish completion installed to $COMPLETION_DIR/ason.fish"
}

# Detect available shells and install completion
install_for_shell() {
    local shell="$1"

    case "$shell" in
        "bash")
            if command -v bash &> /dev/null; then
                install_bash_completion
            else
                print_warning "Bash not found, skipping bash completion"
            fi
            ;;
        "zsh")
            if command -v zsh &> /dev/null; then
                install_zsh_completion
            else
                print_warning "Zsh not found, skipping zsh completion"
            fi
            ;;
        "fish")
            if command -v fish &> /dev/null; then
                install_fish_completion
            else
                print_warning "Fish not found, skipping fish completion"
            fi
            ;;
        "all")
            install_for_shell "bash"
            install_for_shell "zsh"
            install_for_shell "fish"
            ;;
        *)
            print_error "Unknown shell: $shell"
            print_info "Supported shells: bash, zsh, fish, all"
            exit 1
            ;;
    esac
}

# Main function
main() {
    print_header

    # Check if binary exists
    check_binary

    # Determine what to install
    if [[ $# -eq 0 ]]; then
        # Auto-detect current shell
        current_shell=$(basename "$SHELL")
        print_info "Auto-detected shell: $current_shell"
        install_for_shell "$current_shell"
    else
        # Install for specified shell(s)
        for shell in "$@"; do
            install_for_shell "$shell"
        done
    fi

    echo
    print_success "Installation complete!"
    print_info "You may need to restart your shell or run 'source ~/.bashrc' (or equivalent) to enable completion"

    echo
    print_info "To test completion, try typing:"
    echo "  ason <TAB>"
    echo "  ason new <TAB>"
    echo "  ason remove <TAB>"
}

# Show help
show_help() {
    echo "※ Ason Autocompletion Installer"
    echo
    echo "Usage: $0 [SHELL...]"
    echo
    echo "SHELL can be:"
    echo "  bash    Install bash completion"
    echo "  zsh     Install zsh completion"
    echo "  fish    Install fish completion"
    echo "  all     Install completion for all available shells"
    echo
    echo "If no shell is specified, auto-detects current shell."
    echo
    echo "Examples:"
    echo "  $0              # Install for current shell"
    echo "  $0 bash         # Install bash completion only"
    echo "  $0 bash zsh     # Install for bash and zsh"
    echo "  $0 all          # Install for all available shells"
}

# Handle arguments
if [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    show_help
    exit 0
fi

# Run main function
main "$@"