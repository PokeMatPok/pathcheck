![Pathcheck banner](https://raw.githubusercontent.com/PokeMatPok/pathcheck/master/pathcheck.png)

# Pathcheck: simple windows PATH utility

Pathcheck gives you easy access to the PATH state from your command line. It allows you to run commands from path even if they aren't loaded in your current session.

## Features

- **Audit** - Scan your PATH for duplicates, invalid entries, and potential security issues
- **Which** - Locate executables across both system and user PATH entries
- **Cast** - Execute commands directly from PATH, even when they're not available in your current shell
- **Diff** - Compare registry PATH entries against your active shell environment
- **Unique** - List all unique PATH entries in machine-readable format

## Installation

Download the latest release from the releases page, or build from source:
```bash
go build -o pathcheck.exe
```

## Usage
Pathcheck is an CLI tool, an CLI tool can just be called via powershell or cmd. Please make sure you navigated to the folder where the executable lives and run it directly:
```bash
.\pathcheck.exe <command> [args]
```
or register it to path and run it globally:
```bash
pathcheck <command> [args]
```

### Commands

#### audit
Performs a comprehensive scan of your PATH configuration:
- Identifies duplicate entries between system and user PATH
- Detects invalid or non-existent directories
- Flags potential security risks (like system directories in user PATH)
- Provides summary statistics and recommendations
```bash
pathcheck audit
```

#### which
Searches for an executable across all PATH entries and displays its location:
```bash
pathcheck which python
```

#### cast
Runs a command from PATH, even if it hasn't been loaded into your current shell session. Useful after installing new software without restarting your terminal:
```bash
pathcheck cast node --version
```

If multiple matches exist, you'll be prompted to select which one to execute.

#### diff
Compares the PATH entries stored in the Windows registry against what's currently active in your shell. Helps identify:
- Entries that require a shell restart to take effect
- Dynamically added entries that aren't persisted
- Environment variables that couldn't be resolved
```bash
pathcheck diff
```

#### unique
Outputs all unique PATH entries, one per line. Useful for scripting or further processing:
```bash
pathcheck unique
```

## Why Pathcheck?

Windows PATH management can be frustrating. Changes to PATH in the registry don't immediately affect running shells, and tracking down which version of a tool is actually being used can be tedious. Pathcheck bridges this gap by giving you direct access to both the registry state and your current environment.

## Built with Go

## Contributing

If you feel like something is missing either open an issue or fork it, make your changes and create a PR.

## License

[Add your license here]
