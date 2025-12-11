# 🫂 CliWaifuTamagotchi

Preview:
![Result](screenshots/result.gif)
---
###### You can turn the avatar to **husbando** in `~/.config/cliwaifutamagotchi/settings.json`!
![Husbando](screenshots/husbando-preview.jpg)

![Repo size](https://img.shields.io/github/repo-size/HenryLoM/CliWaifuTamagotchi?color=lightgrey)
![Commits](https://img.shields.io/github/commit-activity/t/HenryLoM/CliWaifuTamagotchi/main?color=blue)
![Last commit](https://img.shields.io/github/last-commit/HenryLoM/CliWaifuTamagotchi?color=informational)
![License](https://img.shields.io/github/license/HenryLoM/CliWaifuTamagotchi?color=orange)

## 📑 Table of Contents
- [✨ Overview](#-overview)
- [🎬 Launching Process](#-launching-process)
- [🎨 Customization](#-customization)
- [📂 Project Structure](#-project-structure)
- [⚙️ Core Scripts](#-core-scripts)
    - [main.go](#maingo)
    - [utils/app-utils.go](#utilsapp-utilsgo)
    - [utils/commands-utils.go](#utilscommands-utilsgo)
    - [utils/happiness-utils.go](#utilshappiness-utilsgo)
    - [utils/palette-handler.go](#utilspalette-handlergo)
    - [utils/settings-handler.go](#utilssettings-handlergo)
    - [utils/encouragements-handler.go](#utilsencouragements-handlergo)
    - [utils/gifts-handler.go](#utilsgifts-handlergo)
- [📜 Notes & Error handling](#-notes--error-handling)
- [🛐 Special thanks](#-special-thanks)

---

## ✨ Overview
CliWaifuTamagotchi is a **terminal-based tamagotchi** that:

- Renders **ASCII expressions and clothes**.
- Provides a small set of **interactions**: Encourage, Gift, Dress Up, Background Mode, Quit.
- Uses a **persistent color palette** stored in `~/.config/cliwaifutamagotchi/palette.json`.
- Uses **persistent detail settings** stored in `~/.config/cliwaifutamagotchi/settings.json`.
- Customize some of the functions editing **`words-of-encouragement.txt` and `gifts.json`** in the same directory.
- Has minimal UI built using **`tview` and `tcell`**.
- Has **Vim-style navigation**: Use `h`, `j`, `k`, `l` keys for intuitive navigation and selection (Must be enabled in **settings.json**).

No tons of loops - only one function that repeats itself every 5 seconds. Everything handles and updates according to it.

---

## 🎬 Launching Process

<details>
  <summary><b>Brew</b> (macOS)</summary>

  1. **Install**

  ```bash
  brew install HenryLoM/CliWaifuTamagotchi/cliwt
  ```

  2. **Run**

  ```bash
  cliwt
  ```

  ---

</details>

<details>
  <summary><b>AUR</b> (Arch)</summary>

  1. **Install**

  ```bash
  yay -S cliwt
  ```
  or
  ```bash
  paru -S cliwt
  ```

  2. **Run**

  ```bash
  cliwt
  ```

  ---

</details>

<details>
  <summary><b>Git</b> (Source code)</summary>

  1. **Clone repository**

  ```bash
  git clone https://github.com/HenryLoM/CliWaifuTamagotchi.git
  cd CliWaifuTamagotchi
  ```

  2. **Build app yourself, then run**

  ```bash
  go build -o cliwt
  ./cliwt
  ```

  - **Or run directly for development**

  ```bash
  go run main.go
  ```

  ---

</details>

> **💡 Notes**
>
> * First run creates `~/.config/cliwaifutamagotchi/` directory and `palette.json`, `settings.json` files in it on its own if missing.
> * On macOS, ensure your terminal supports **true color** for best rendering.

---

## 🎨 Customization

1. **Palette**<br>
JSON file is in `~/.config/cliwaifutamagotchi/` ; Named `palette.json`<br>
JSON file's structure:
```
{
  "background": "#1e1e2e",
  "foreground": "#cdd6f4",
  "border": "#cba6f7",
  "accent": "#eba0ac",
  "title": "#b4befe"
}
```
> Note: default palette is Catppuchin (Mocha).

2. **Settings**<br>
JSON file is in `~/.config/cliwaifutamagotchi/` ; Named `settings.json`<br>
JSON file's structure:
```
{
  "name": "Waifu",
  "defaultMessage": "...",
  "vimNavigation": false,
  "avatarType": "waifu",
  "keys": {
    "encourage": "l",
    "dressup": "2",
    "backgroundMode": "b",
    "quit": "q"
  }
}
```
> Note: try to avoid key overrides when using `"vimNavigation": true`.

3. **Words of encouragement**<br>
TXT file is in `~/.config/cliwaifutamagotchi/` ; Named `words-of-encouragement.txt`<br>
> Note: It's extensible!

4. **Gifts**<br>
JSON file is in `~/.config/cliwaifutamagotchi/` ; Named `gifts.json`<br>
> Note: It's extensible!

---

## 📂 Project Structure

```
CliWaifuTamagotchi/
│
├── README.md
├── LICENSE
├── .gitignore
├── go.mod
├── go.sum
├── main.go                             # Main file that launches the project
│
├── screenshots/
│   ├── result.gif
│   ├── reactions.jpg
│   └── husbando-preview.jpg
│
└── utils/
    │
    ├── ascii-arts/
    │   │
    │   ├── waifu/                      # Arts for waifu avatar
    │   │   ├── clothes/...             # ASCII bodies
    │   │   └── expressions/...         # ASCII heads
    │   │
    │   └── husbando/...                # Arts for husbando avatar
    │       ├── clothes/...             # ASCII bodies
    │       └── expressions/...         # ASCII heads
    │
    ├── assets/
    │   └── words-of-encouragement.txt  # List of lines for Encouragement function
    │
    ├── app-utils.go                    # Main helpers
    ├── commands-utils.go               # Functions for the Action Space
    ├── happiness-utils.go              # Happiness scoring system
    ├── palette-handler.go              # Handling palette out of the file
    ├── settings-handler.go             # Handling settings out of the file
    ├── encouragements-handler.go       # Handling encouragements out of the file
    └── gifts-handler.go                # Handling gifts out of the file
```

---

## ⚙️ Core Scripts

### **main.go**

* Loads ASCII **head, blink frames, and body**.
* Displays **actions menu**: Encourage, Dress Up, Quit.
* Handles **user input** (keys and navigation).
* Queues UI updates safely using `app.QueueUpdateDraw` via `UIEventsChan` that keeps UI changes in order.

### **utils/app-utils.go**

* Helper functions for **loading ASCII files**.
* Manages **UI rendering** and **widget updates**.

### **utils/commands-utils.go**

* Implements **interactions logic**:

  * `Encourage`: random encouraging phrase + happy frame.
  * `GiftMenu`: choose gifts, apply happiness, show reaction.
  * `DressUp`: swaps body/outfit based on selection.
  * `BackgroundMode`: fills the TUI with Waifu, removing all of the odd elements.
* Manages UI state and async updates via UIEventsChan.
* Caches custotmizable files to reduce disk reads.

### **utils/happiness-utils.go**

* Handles the bar and changes emotions of the avatar.
* Handles the happiness scores.

### **utils/palette-handler.go**

* Loads palette from `~/.config/cliwaifutamagotchi/palette.json`.
* Creates **default palette** if missing.
* Provides **color application** helpers.

### **utils/settings-handler.go**

* Loads settings from `~/.config/cliwaifutamagotchi/settings.json`.
* Creates **default settings** if missing.

### **utils/encouragements-handler.go**

* Loads settings from `~/.config/cliwaifutamagotchi/words-of-encouragement.txt`.
* Restores **default encouragements** from relative directory if missing.

### **utils/gifts-handler.go**

* Loads settings from `~/.config/cliwaifutamagotchi/gifts.json`.
* Restores **default gifts** if missing.

---

## 📜 Notes & Error handling

#### **Errors:**
* See errors you can't explain? Try to remove `~/.config/cliwaifutamagotchi/` directory (you can do a backup). If it worked, it means customization was updated and there was a conflict.
* If error remains, leave the issue, we could solve it together!

#### **Warning:**
* Missing/malformed ASCII files may cause a wrong output; handle carefully if modifying assets inside the structure.

#### **Read if you want to contribute:**
* The project lives only because there are people who use it. Let's make sure we build it for people, not to earn another achievement for our profiles.
* Keep the code clean and constructive.
* Keep the project PG.
* Keep the Unix philosophy principle: a program should do one thing and do it well.
* The two main goals of the project:
  * A TUI that is as customizable as we can make it.
  * A tool that is as lightweight as possible, since the project assumes users leave it running in the background.

#### **Future plans you can help with:**
* More interactions (feeding, timed events, stats).
* Save selected outfit and preferences.
* Unit tests and error handling improvements.
* Custom separate font support (because a lot of people meet problems with visuals with their fonts).
* Maybe "Pose Mode" - loop animation or specific pose to select and have on the background.
* Maybe separate module handle stderr so Waifu reacts to the errors you get during your work.
* Maybe separate module to use a ChatBot.

---

## 🛐 Special thanks

- **[sutemo](https://sutemo.itch.io/)** — for the amaazing [female](https://sutemo.itch.io/female-character) & [male](https://sutemo.itch.io/male-character-sprite-for-visual-novel) sprites you can see in the project.
- **[mininit](https://github.com/mininit)** — for embedding all assets and enabling a clean, pure build process.
- **[Ali Medhat](https://github.com/Alimedhat000)** — for adding Vim-style navigation.
- **[Isaac Hesslegrave](https://github.com/HeadedBranch)** — for implementing Arch Linux support via `yay` and `paru`.

---

⤴︎ Return to the [📑 Table of Contents](#-table-of-contents) ⤴︎