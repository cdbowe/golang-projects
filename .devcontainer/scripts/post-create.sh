#!/bin/bash
set -e

echo "Running post-create setup..."
echo "  WORKSPACE_DIR: ${WORKSPACE_DIR}"

###########################################
# Git Safe Directory
###########################################

git config --global --add safe.directory "${WORKSPACE_DIR}" 2>/dev/null || true

###########################################
# Docker Socket
###########################################

if [ -S /var/run/docker.sock ]; then
  sudo chgrp docker /var/run/docker.sock 2>/dev/null || true
  sudo chmod g+rw /var/run/docker.sock 2>/dev/null || true
fi

###########################################
# Project Dependencies
###########################################

if [ -f "go.mod" ]; then
  echo "Running: go mod download..."
  go mod download
fi

###########################################
# Template Setup
###########################################

echo "Running: MAIN_DIR="${WORKTREE_MAIN_DIR:-${WORKSPACE_DIR}/main}" ......"
MAIN_DIR="${WORKTREE_MAIN_DIR:-${WORKSPACE_DIR}/main}"
mkdir -p "$MAIN_DIR"
# Must precede init to have any effect on the initial branch name.
git config --global init.defaultBranch main
if [ ! -e "$MAIN_DIR/.git" ]; then
  git -C "$MAIN_DIR" init -q
  # Dev Containers normally copies the host ~/.gitconfig in, but not always
  # (headless `devcontainer up`, or the feature turned off). Without an
  # identity the commit below fails, and set -e would abort the rest of
  # post-create — including the install.sh steps.
  if ! git -C "$MAIN_DIR" var GIT_AUTHOR_IDENT >/dev/null 2>&1; then
    git -C "$MAIN_DIR" config user.name "devcontainer"
    git -C "$MAIN_DIR" config user.email "devcontainer@localhost"
  fi
  git -C "$MAIN_DIR" commit -q --allow-empty -m "initial commit"
  # Redundant on git >= 2.28 given init.defaultBranch above; kept for older git,
  # and guarded so it can't fail the script when the branch is already main.
  if git -C "$MAIN_DIR" show-ref -q --verify refs/heads/master; then
    git -C "$MAIN_DIR" branch -m master main
  fi
fi

echo "Running: if [ -x /opt/claude-code-tools/install.sh ]; then bash /opt/claude-code-tools/install.sh --all --dir "${CLAUDE_CONFIG_DIR:-$HOME/.claude}"; fi..."
if [ -x /opt/claude-code-tools/install.sh ]; then bash /opt/claude-code-tools/install.sh --all --dir "${CLAUDE_CONFIG_DIR:-$HOME/.claude}"; fi

echo "Running: if [ -x /opt/claude-code-tools/install.sh ]; then bash /opt/claude-code-tools/install.sh --all --dir "${WORKSPACE_DIR}/.claude" --with-local; fi..."
if [ -x /opt/claude-code-tools/install.sh ]; then bash /opt/claude-code-tools/install.sh --all --dir "${WORKSPACE_DIR}/.claude" --with-local; fi

echo ""
echo "Post-create setup complete!"
