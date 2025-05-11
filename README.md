<img src="screenshots/prev.jpg" width="100%" height="150px" alt="Preview">

# 💫 About Project:

## 🧸 Yo-Link — Link shortening service with API and web interface

Yo-Link is a simple and fast URL shortening service written in Go. It supports REST API, authorization, and link customization

## 💻 Tech Stack:

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![SQLite](https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white)

### ✨ Key Features

- **URL Shortening Engine:** Shortening long URLs into short ones (for example,https://example.com/very-long-path → http://localhost:8082/yourAlias or http://yourdomain:8082/yourAlias)
- **Smart Storage:** All your links are stored securely in a local **SQLite** database, ensuring quick and reliable access without any third-party dependencies.
- **Link Protection:** HTTP Basic Auth for access protection.
- **Flexible setup:** Configuration via a YAML file.
- **Auto-start:** Deployment via systemd (included in the repository).
- **Integration with services:** REST API for integration with other services.

### 📄 API Documentation

| Заголовок 1 | Заголовок 2     | Заголовок 3                           | Requires Auth |
| ----------- | --------------- | ------------------------------------- | ------------- |
| `POST`      | `/url`          | _Create short URL from original link_ | ✅ **Yes**    |
| `GET`       | `/{your_alias}` | _Redirect to original UR_             | ❌ **No**     |

#### Request Body Format (JSON)

- When making a request to `/url`, you must send JSON in the format

`{
  "url": "https://example.com/very-long-url",
  "alias": "custom-name"
}`

##### Where:

- `/url` (required) - the original long link
- **`alias`** (optional) - the desired short name for the link (if not specified, it will be generated automatically)

### 📦 Deployment

- Automatic deployment via GitHub Actions:

  1. Specify the secrets (`DEPLOY_SSH_KEY`, `AUTH_PASS`) in the repository settings.
  2. Manually launch the workflow **Deploy App** by specifying the version (tag).

- Configuration for systemd: deployment/yo-link.service.

## 📸 Screenshots

- **Command** `/start`

  ![Command /start](screenshots/start.jpg)
  _Welcome to the bot’s main menu — your gateway to effortless link saving and exploration._

- **Command** `/help`

  ![Command /help](screenshots/help.jpg)
  _Helpful instructions at your fingertips — learn how to save links quickly and efficiently._

- **Command** `/random`

  ![Command /random](screenshots/random.jpg)
  _Discover hidden treasures from your saved collection with a simple `/random` command._

  ![Error handling of an existing link](screenshots/double.jpg)
  _Smart duplicate detection — the bot kindly lets you know if you've already saved that link._

![error handling if the user has not uploaded any links](screenshots/noSavedLinks.jpg)
_Friendly reminder when no links are saved yet — encouraging you to start your bookmarking journey!_

---

## 🛠️ Installation

### Prerequisites

- [Go 1.22+](https://golang.org/dl/)

## 🪭 Quick Setup

### Clone repository

- git clone https://github.com/yowie645/Yo-Link.git
- cd Yo-Link

### Run

- go run \path\Yo-link\cmd\yo-link\main.go

Give LinkSaverBot a try and never lose track of your favorite links again!
