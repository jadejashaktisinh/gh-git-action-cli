# 📋 Project Plan: gh-git-action-cli (GitHub CLI Extension)

A precompiled Go-based GitHub CLI extension (`gh git-action-cli`) to trigger, track, and manage manual GitHub Actions (`workflow_dispatch`) with local configuration, local history tracking, and an interactive terminal UI.

---

## 🛠️ System Architecture & Stack

- **CLI Framework:** [Cobra](https://github.com) (Flags, commands, orchestration)
- **GitHub Integration:** [go-gh](https://github.com) (Native host authentication, zero-token setup)
- **Terminal UI (TUI):** [Bubble Tea Engine](https://github.com) (Interactive lists, prompts, spinners)
- **Configuration Engine:** [Viper](https://github.com) (Manages global configurations via `~/.config/gh-git-action-cli/config.yaml`)
- **Local History Database:** [bbolt](https://github.com) (Single-file embedded KV store for logging run execution metadata)

---

## 🗺️ Project Phases & Checklists

### 🟥 Phase 1: Project Setup & Repository Initialization
*Objective: Establish a unified Go module layout compliant with GitHub CLI Extension standards.*

- [x] Run `gh extension create --precompiled=go gh-git-action-cli` to scaffold the official project layout.
- [x] Initialize dependencies in `go.mod`.
- [x] Implement the core `main.go` entry point calling `cmd.Execute()`.
- [x] Create the directory structure: `/cmd`, `/config`, `/db`, and `/tui`.

---

### 🟨 Phase 2: Configuration & Database Layers
*Objective: Build persistent local state engines for tool configuration and history tracking.*

- [x] **Viper Storage Engine (`/config/config.go`):**
  - [x] Auto-create storage paths based on OS spec (`~/.config/gh-git-action-cli/`).
  - [x] Define settings schema (Default repository, polling interval, default flags).
- [x] **bbolt History Storage Engine (`/db/database.go`):**
  - [x] Initialize `history.db` embedded file.
  - [x] Create `runs` bucket.
  - [x] Implement `SaveRun(run Record)`.
  - [x] Implement `GetHistory() []Record`.

---

### 🟩 Phase 3: Core Commands & Cobra Integration
*Objective: Develop standard, automation-friendly flag behaviors using Cobra commands.*

- [x] **Root Command (`/cmd/root.go`):**
  - [x] Map execution syntax: `gh git-action-cli [workflow-file.yml]`.
  - [x] Bind flags: `--repo (-r)`, `--branch (-b)`, `--input (-i)`.
  - [x] Add routing condition: If *no* flags are provided, launch the interactive TUI. If flags are provided, bypass TUI for headless execution.
- [x] **History Command (`/cmd/history.go`):**
  - [x] Implement `gh git-action-cli history`.
  - [x] Pull records from bbolt and format a clean terminal matrix/table view.
- [x] **Status Command (`/cmd/status.go`):**
  - [x] Implement `gh git-action-cli status <run-id>`.
  - [x] Query real-time API states (`queued`, `in_progress`, `completed`).

---

### 🟦 Phase 4: Interactive TUI (Bubble Tea)
*Objective: Build an interactive, menu-driven terminal experience for manual workflows.*

- [x] **Selection Menus:**
  - [x] Fetch remote repositories if none specified.
  - [x] Parse `.github/workflows/` files via GitHub API to render an interactive workflow selector.
- [ ] **Dynamic Forms:**
  - [ ] Read `workflow_dispatch.inputs` schemas from target YAMLs.
  - [ ] Dynamically generate interactive text fields/select dropdowns based on YAML expectations.
- [ ] **Live Polling UI:**
  - [ ] Build a Bubble Tea spinner state linked to a background HTTP polling ticker.
  - [ ] Stream execution lifecycle success/failure indicators directly onto the terminal screen.

---

### 🟪 Phase 5: Build, Release & CI/CD Pipeline
*Objective: Automate builds for multiple operating systems and architectures so users can install it.*

- [x] Configure `.github/workflows/release.yml` using the canonical GitHub Extensions release template.
- [ ] Setup cross-compilation matrix matching:
  - [ ] `linux/amd64`, `linux/arm64`
  - [ ] `darwin/amd64`, `darwin/arm64`
  - [ ] `windows/amd64`
- [ ] Draft an extension installation validation test:
  ```bash
  gh extension install <your-username>/gh-git-action-cli
  gh git-action-cli --help
  ```

---

## 💾 Local Database Execution Schema (JSON Reference)
Every run executed through the engine logs into `history.db` with this data layout:

```json
{
  "run_id": 4829104822,
  "timestamp": "2026-05-22T11:47:00Z",
  "repository": "your-org/your-repo",
  "workflow": "deploy.yml",
  "branch": "main",
  "inputs": {
    "environment": "production",
    "debug_mode": "true"
  },
  "conclusion": "success"
}
```
