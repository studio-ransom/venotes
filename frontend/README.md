# 🗒️ Svelte Notes — A Fully Client-Side Notes App

**Svelte Notes** is a fully client-side note organization system — no servers, no databases, no setup required.  
Just open the files in your browser and start organizing your thoughts like a Discord-style workspace for your life, projects, or ideas.

## 🚀 Features

- **100% Client-Side** — No backend, no external services — everything runs in your browser.
- **Svelte + Tailwind CSS** — Built with [Svelte](https://svelte.dev) for reactivity and [Tailwind CSS](https://tailwindcss.com) for clean, responsive UI.
- **Discord-Inspired Structure** — Organize your notes like Discord servers ("guilds") and channels:
  - **Guilds** → high-level categories (e.g., Work, Personal, Projects)
  - **Channels** → subcategories inside each guild (e.g., Coding, Ideas, Questions)
  - **Logs** → timestamped entries inside each channel (like messages)
- **Mobile-Responsive** — Works seamlessly on phones and desktops.
- **Local Storage API** — A simple JavaScript API for loading, saving, and exporting all your data.

## 🧠 Concept Overview

| Discord Concept | Svelte Notes Equivalent |
| --------------- | ----------------------- |
| Server (Guild)  | A workspace or category |
| Channel         | A topic or project area |
| Message         | A note or log entry     |

Each log entry contains a **timestamp**, **text**, and potentially metadata (tags, edited status, etc. in future versions).

## 🏗️ Development Setup

```bash
# Install dependencies
npm install

# Run locally
npm run dev

# Build for production
npm run build
```

Then open `dist/index.html` in any browser — no server needed!

## 🧩 JavaScript API Design

The goal is to separate **UI logic** from **data logic**.
The client-side API manages all read/write operations, data exports, and imports.

### API Functions

#### Guilds
```js
get_guilds() -> [string]
load_guild(guildName: string) -> GuildData
save_guild(guildName: string, data: GuildData)
create_guild(guildName: string)
delete_guild(guildName: string)
```

#### Channels
```js
get_channels(guild: string) -> [string]
create_channel(guild: string, channelName: string)
delete_channel(guild: string, channelName: string)
```

#### Logs
```js
get_logs(guild: string, channel: string) -> [LogData]
channel_post_log(guild: string, channel: string, content: string)
channel_delete_log(guild: string, channel: string, logId: string)
```

#### Storage
```js
load_all()
save_all()
export_data() -> Downloads JSON file
import_data(file) -> Promise
```

## 🧱 UI Layout

**Top Bar**
- Current Guild + Channel name
- Search icon / bar
- Export/Import buttons

**Left Sidebar**
- List of guilds
- Channel list under active guild
- Mobile-responsive with slide-out menu

**Main Panel**
- Scrollable list of logs
- Input box to add new log entries
- Timestamp display per entry

**Responsive Behavior**
- On mobile, sidebar collapses into slide-out menu
- Main area focuses on active channel's logs

## 🔧 Tech Stack

* **Frontend:** Svelte 4
* **Styling:** Tailwind CSS
* **Storage:** LocalStorage
* **Build Tool:** Vite

## 🧪 Future Ideas

* Tagging and searching logs
* Markdown or rich text support
* Export / Import as ZIP
* Theme customization (light/dark)
* Sync through browser file API (optional, still offline)
