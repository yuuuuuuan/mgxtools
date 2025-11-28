# 🚀 mgxtools — High-Speed Recursive File Downloader for mgx.dev

`mgxtools` is a command-line utility designed to recursively scan, build directory trees, and rapid-download files from `mgx.dev` chat workspaces with up to **20 concurrent workers** and fully animated **per-file progress bars**.

This tool is ideal for users who need to export or mirror an entire mgx.dev workspace—including nested folders and large files—quickly and reliably.

---

## ✨ Features

### 🔍 Recursive Directory Tree Builder
- Automatically traverses folders returned by  
  `https://mgx.dev/api/v1/files?path=<...>`
- Builds a complete in-memory JSON file tree.

### ⚡ Ultra-Fast Parallel Downloader
- Up to **20 simultaneous file downloads**
- Efficient worker-pool concurrency
- Automatically creates nested directories

### 📊 Real-Time Progress Bars
- Individual progress bar per file
- Shows transferred size (KiB), percentage, and filename
- Uses `mpb` for smooth, multi-line terminal rendering

### 🔑 Secure Authorization Handling
- The mgx.dev API token (`authorization`) is stored in a local `.env` file
- No secrets or tokens are hardcoded in the source code

### 🎯 Clean Module Architecture
- `internal/api` → mgx.dev API client  
- `internal/tree` → recursive file-tree builder  
- `internal/downloader` → concurrent download manager  
- Root CLI → accepts the chat session ID and output folder  

---

## 📦 Installation

### 1. Clone the repository
```bash
git clone https://github.com/yourname/mgxtools.git
cd mgxtools
```

### 2. Install dependencies
```bash
go mod tidy
```

---

## 🔧 Configuration

### Create a `.env` file:

```
MGX_AUTH=your_mgxdev_authorization_token_here
MGX_BASE=https://mgx.dev
```

### Required environment variables

| Variable | Description |
|---------|-------------|
| `MGX_AUTH` | mgx.dev JWT authorization token |
| `MGX_BASE` | Base URL (usually `https://mgx.dev`) |

---

## 🧪 Usage

The CLI accepts a **mgx.dev chat workspace ID** (like `ac4a88ea71c14d088ab3557312439f50`) and a local output directory.

### Example:

```bash
go run . ac4a88ea71c14d088ab3557312439f50 ./backup
```

This will:

1. Build the complete file tree  
2. Display the JSON tree  
3. Download all files concurrently with real-time progress bars  
4. Save everything into `./backup`

---

## 📁 Project Structure

```
mgxtools/
│
├── cmd/
│   └── mgxtools/        → CLI entrypoint
│
├── internal/
│   ├── api/             → mgx.dev API wrapper
│   ├── tree/            → recursive directory scanner
│   └── downloader/      → concurrent download system
│
├── .env                 → your mgx.dev token (ignored by git)
├── go.mod
└── README.md
```

---

## 🏗 How It Works

### 1. Build Directory Tree  
`tree.BuildTree()` recursively requests:

```
GET /api/v1/files?path=<encoded-path>
```

### 2. Spawn Concurrent Downloads  
Up to **20 workers** run in parallel using a semaphore.

### 3. Per-File Progress Bars  
Each file uses:

- `mpb.Progress`
- `mpb.Bar`
- Filename / Size / Percentage decorators

### 4. Save to Disk  
Directory structure is preserved exactly as on mgx.dev.

---

## ⚠️ Notes

- Only valid mgx.dev authorization tokens are supported  
- API behavior may change if mgx.dev updates endpoints  
- Animated progress bars require a real terminal  

---

## 🤝 Contributing

Pull requests are welcome.

Before submitting changes:

```bash
go fmt ./...
go vet ./...
golangci-lint run
```

---

## 📄 License

MIT License.  
You may freely modify or distribute this tool.
