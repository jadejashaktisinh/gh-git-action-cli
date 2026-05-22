# 🚀 gh-git-action-cli

**Execute GitHub Actions workflows 100% locally.**

`gh-git-action-cli` is a high-performance GitHub CLI extension that intercepts and executes GitHub Actions workflows directly on your development machine. It eliminates reliance on cloud runner infrastructure, natively injects `.env` files without cloud exposure, and acts as a fast, private automation runner.

---

## 🌟 Key Features

- **Native Execution:** Runs steps directly on your host OS for maximum speed—no waiting for cloud queues.
- **Zero-Leak Secret Injection:** Injects local `.env` files safely into your steps without ever sending data to the cloud.
- **Git Interceptor:** Automatically block broken pushes by running local CI checks via a generated git pre-push hook.
- **Interactive TUI:** A beautiful terminal interface to browse and select local workflows/jobs.
- **Offline History:** Audit every local run via an embedded `bbolt` database.
- **Zero-Token Setup:** Works seamlessly with your existing `gh` CLI authentication.

---

## 🛠️ Installation

### Prerequisites
- [Go](https://go.dev/doc/install) 1.21+
- [GitHub CLI](https://cli.github.com/) (Optional, but recommended as an extension)

### Build from Source
```bash
git clone https://github.com/your-username/gh-git-action-cli.git
cd gh-git-action-cli
go build -o gh-git-action-cli main.go
```

To make it available globally:
```bash
cp gh-git-action-cli ~/.local/bin/  # Or any directory in your PATH
```

---

## 🚀 Usage

### 1. Interactive Mode
Simply run the tool to launch the TUI dashboard:
```bash
./gh-git-action-cli
```

### 2. Headless Mode (CLI)
Execute a specific job from a workflow file:
```bash
./gh-git-action-cli --job test .github/workflows/ci.yml
```

### 3. Secret Injection
Inject a local environment profile:
```bash
./gh-git-action-cli --job deploy --env-file .env.production
```

### 4. Git Hook Integration
Install a pre-push hook to ensure CI passes locally before pushing:
```bash
./gh-git-action-cli init-hook
```

### 5. Audit History
View your local execution metrics:
```bash
./gh-git-action-cli history
```

---

## 💾 Local History Schema
Every run is logged locally into a private database on your machine:
- **Timestamp:** When the run occurred.
- **Workflow:** Which file was executed.
- **Job:** The specific job targeted.
- **Mode:** Native-shell vs Container.
- **Status:** Passed or Failed.

---

## 📄 License
This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing
Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.
