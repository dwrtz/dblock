# dblock

A productivity tool for blocking websites when you need to focus.

`dblock` is a command-line tool written in Go for managing the contents of your `/etc/hosts` file.

## Features

- **Block/Unblock Domains:** Easily add or remove domains from your `/etc/hosts` file.
- **Scheduling:** Schedule blocking or unblocking operations for a specified duration.
- **Automatic Subdomains:** Automatically handles common subdomains like `www`.
- **Custom Subdomains:** Manually specify additional subdomains to block.
- **Custom Hosts File Location:** Support for custom hosts file paths.

## Table of Contents

- [dblock](#dblock)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Basic Commands](#basic-commands)
  - [Configuration](#configuration)
    - [Sample Configuration](#sample-configuration)
    - [Custom Configuration File](#custom-configuration-file)
  - [Scheduling Operations](#scheduling-operations)
    - [Examples](#examples)
    - [How Scheduling Works](#how-scheduling-works)
  - [The `~/.dblock` Folder](#the-dblock-folder)
  - [Uninstallation](#uninstallation)
  - [License](#license)

## Installation

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/dwrtz/dblock.git
    cd dblock
    ```

2. **Build the Executable:**

    ```bash
    go build
    ```

3. **Move Executable to `/usr/local/bin`:**

    ```bash
    sudo mv dblock /usr/local/bin
    ```

## Usage

### Basic Commands

- Enable Blocking:

    ```bash
    sudo dblock enable
    ```

- Disable Blocking:

    ```bash
    sudo dblock disable
    ```

- Check Status:

    ```bash
    dblock status
    ```

- List Configured Domains:

    ```bash
    dblock list
    ```

## Configuration

On the first run, `dblock` will create a default configuration file at `~/.dblock/default.yaml` if it does not already exist. This file contains a list of sample domains to block.

You can edit this file to specify the domains and subdomains you wish to block.

### Sample Configuration

```yaml
hosts_file: "/etc/hosts" # Optional: specify custom hosts file location
domains:
  - x.com
  - twitter.com
  - youtube.com
  - reddit.com
subdomains:
  - blog.example.com
  - mail.example.org
```

- `hosts_file`: (Optional) Path to the hosts file. Defaults to `/etc/hosts`.
- `domains`: List of domains to block. Common subdomains like `www` are automatically included.
- `subdomains`: (Optional) Manually specify additional subdomains to block.

### Custom Configuration File

You can specify a custom configuration file using the `-c` or `--config` flag.

```bash
sudo dblock enable -c /path/to/custom_config.yaml
```

## Scheduling Operations

You can schedule enable or disable operations for a specified duration using the `-t` or `--timeout` flag.

### Examples

- Block Domains for 60 minutes:

    ```bash
    sudo dblock enable -t 60
    ```

- Unblock Domains for 30 minutes:

    ```bash
    sudo dblock unblock -t 30
    ```

### How Scheduling Works

The tool runs an in-process background goroutine that waits (sleep) for the specified duration and then executes the reverse dblock command. If you close the terminal or the system restarts before the timeout elapses, the scheduled operation will not occur.


## The `~/.dblock` Folder

`dblock` uses the `~/.dblock` directory to store configurations, backups, and logs.
- **`~/.dblock/default.yaml`**: Your editable config file for specifying domains to block. 
- **Backups**: Before modifying the hosts file, a backup is stored in `~/.dblock/backups/`.
- **Logs**: Operation logs are stored in `~/.dblock/logs/`.

## Uninstallation

To uninstall `dblock`:
1. Remove the Executable:
    ```bash
    sudo rm /usr/local/bin/dblock
    ```
2. Remove the `~/.dblock` Directory:
    ```bash
    rm -rf ~/.dblock
    ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

