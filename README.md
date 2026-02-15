<h1 align="center">
	<br>
	⚡️
	Zap
	<br>
	<br>
</h1>

> A simple cross platform zombie process killer

Works on macOS, Linux, and Windows. A lightweight Go alternative to finding and killing network-bound processes.

## Install

### Download Binary (Recommended)

#### Linux / macOS
```sh
# Download (replace 'linux' with 'darwin' for macOS)
curl -LO https://github.com/brinzl/zap/releases/latest/download/zap-linux-amd64

# Make executable and move to PATH
chmod +x zap-linux-amd64
sudo mv zap-linux-amd64 /usr/local/bin/zap
```

#### Windows

Download `zap-windows-amd64.exe` from the [Releases page](https://github.com/brinzl/zap/releases), rename it to `zap.exe`, and add it to your PATH.

### Using Go

```sh
go install github.com/brinzl/zap@latest
```

## Usage

```
$ zap --help

zap helps you find and kill processes bound to network ports...

Usage:
  zap [flags]

Flags:
  -f, --force       Skip confirmation prompt
  -h, --help        help for zap
  -p, --port int    Kill process by port
  -v, --version     version for zap
```

### Interactive Mode

Run `zap` without arguments to see an interactive list of all listening processes:

```sh
$ zap
```

Use arrow keys to navigate, type to search by process name, and press Enter to select. You'll be asked to confirm before killing.

### Kill by Port

Kill a process running on a specific port:

```sh
# With confirmation prompt
$ zap -p 3000

# Skip confirmation
$ zap -p 3000 --force
$ zap -p 3000 -f
```

### Examples

```sh
# Interactive: browse and kill processes
$ zap

# Kill process on port 8080 (with confirmation)
$ zap --port 8080

# Kill process on port 3000 (skip confirmation)
$ zap -p 3000 -f
```
