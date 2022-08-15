# vim-note
vim-note is a memo that can be stored in online storage from the terminal

- storage
  - firebase
- editor
  - vim

## Installation
```bash
go install github.com/dogerescat/vim-note@latest
```
create config.toml(~/vim-note/config.toml)
```bash
[Firebase]
keyPath = "your_firestore_key.json"
storageBucket = "your-fire-store-bucket.com"
```
## Features
create new memo
```bash
vim-note new filename
```

show memo list
```bash
vim-note list
```


edit memo
```bash
vim-note edit [filename]
```
