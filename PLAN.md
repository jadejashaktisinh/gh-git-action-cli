# 📋 Project Plan: gh-trigger (Offline Local Workflow Engine)

A precompiled Go-based GitHub CLI extension (`gh trigger`) that intercepts and executes GitHub Actions workflows 100% locally on the developer's computer. It eliminates reliance on cloud runner infrastructure, natively injects `.env` files without cloud exposure, and acts as a fast, private automation runner.

---

## 🛠️ System Architecture & Stack

- **CLI Framework:** [Cobra](https://github.com) (Commands, global flags, execution orchestration)
- **YAML Parser:** [go-yaml/yaml](https://gopkg.in) (Translates GitHub Workflow schemas into native Go execution models)
- **Execution Drivers:** 
  - **Native Engine:** Go `os/exec` (Runs steps directly on the host machine filesystem for maximum speed)
  - **Container Engine:** [Docker SDK for Go](https://github.com) (Optional layer to run workflows inside isolated clean environments)
- **Terminal UI (TUI):** [Bubble Tea Engine](https://github.com) (Interactive menus, live step logging, spinners)
- **Configuration & History Store:** [Viper](https://github.com) & [bbolt](https://github.com) (Tracks local execution metrics, telemetry, and execution logs completely offline)

---

## 🗺️ Project Phases & Checklists

### 🟥 Phase 1: Local Project Foundation
*Objective: Build an independent, offline-first codebase structure.*

- [ ] Initialize standard Go application layout with the core `main.go`.
- [ ] Run dependency injection setup in `go.mod`:
  ```bash
  go get ://github.com
  go get gopkg.in/yaml.v3
  go get ://github.com
  go get ://github.com
  ```
- [ ] Establish application folder layout: `/cmd`, `/parser`, `/runner`, `/db`.

---

### 🟨 Phase 2: Workflow Parsing & Secret Injection Engine
*Objective: Read local workspace data safely without sending files to the cloud.*

- [ ] **The YAML Parser (`/parser/workflow.go`):**
  - [ ] Implement structs mapped to GitHub Action specifications (`jobs`, `steps`, `run`, `env`).
  - [ ] Code a discovery loop that finds and processes files in `.github/workflows/`.
- [ ] **Zero-Leak Secret Injector (`/parser/env.go`):**
  - [ ] Build a `.env` loader that parses files safely on the host machine.
  - [ ] Map variable values straight to step environment bindings without editing the workspace files.

---

### 🟩 Phase 3: The Native Execution Layer
*Objective: Run steps directly on the developer's local OS with zero cloud queues.*

- [ ] **Shell Execution Pipeline (`/runner/shell.go`):**
  - [ ] Implement Go `os/exec` orchestration to execute `run:` strings natively.
  - [ ] Bind real-time `Stdout` and `Stderr` streaming so output appears immediately in the terminal.
  - [ ] Ensure execution exits and alerts the user if any single step returns a non-zero exit code.
- [ ] **Step Router Layer:**
  - [ ] Implement a system that skips cloud-exclusive steps (like third-party actions) and securely passes environment tables directly down to native scripts.

---

### 🟦 Phase 4: Developer UX & TUI Diagnostics
*Objective: Provide visual feedback, interactive options, and guardrails.*

- [ ] **Targeted Execution Flags:**
  - [ ] Implement `--job <name>` to execute individual jobs independently.
  - [ ] Implement `--env-file <path>` to dynamically inject variable profiles.
- [ ] **Bubble Tea Dashboard Interface:**
  - [ ] Build an interactive terminal screen that lets users select a local workflow file to run.
  - [ ] Display a live checklist UI mapping step states (`⏳ Pending`, `🔄 Running`, `✅ Passed`, `❌ Failed`).
- [ ] **Smart Local Analytics (bbolt Engine):**
  - [ ] Log metrics for every local execution run (Total runtime, broken steps, run frequency) to compile performance data.

---

### 🟪 Phase 5: The Ultimate Dev Hook (Git Interceptor)
*Objective: Block broken pushes before they ever leave the computer.*

- [ ] **Git Hook Generation Engine (`gh trigger init-hook`):**
  - [ ] Write logic that creates a local `.git/hooks/pre-push` file automatically.
  - [ ] Configure the hook to run your local tool (`gh trigger --job test`) automatically whenever a developer executes `git push`.
  - [ ] If local checks fail, reject the push to ensure the remote repository never receives broken code.

---

## 💾 Local Execution Record Layout (bbolt JSON)
All telemetry and logs stay locked to your personal device inside your single-file embedded database:

```json
{
  "local_run_id": "b4f2c8d1-9012",
  "timestamp": "2026-05-22T12:15:00Z",
  "workflow_file": "ci-checks.yml",
  "target_job": "build-and-test",
  "execution_mode": "native-shell",
  "env_source": ".env.local",
  "duration_ms": 14230,
  "status": "passed"
}
```
