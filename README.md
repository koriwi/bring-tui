# bring-tui (MOSTLY VIBE CODED, YOU HAVE BEEN WARNED)

A terminal UI and CLI for the [Bring!](https://www.getbring.com/) shopping list app.

## Installation

```bash
go install github.com/paulleonhardhellweg/bring-tui@latest
```

## Setup

```bash
bring login
bring lists
bring use Einkauf
```

Log in with your Bring! credentials, list your available shopping lists, then set the default one with `use`.

## TUI

Run `bring` without any arguments to open the interactive TUI.

```
↑/k · ↓/j   navigate
a            add item
e            edit selected item
enter        mark as bought / re-add from recently bought
x / delete   remove item
l            switch list
r            reload
q / ctrl+c   quit
```

## CLI

### List items

```bash
bring list
```

### Add an item

```bash
bring add Milch
bring add Milch:1.5%
bring add Butter : gesalzen
bring add "Schlagsahne : 200ml"
```

Everything before the `:` is the item name, everything after is the description. Both are trimmed of surrounding whitespace.

### Mark an item as bought

```bash
bring done Milch
```

### Remove an item

```bash
bring remove Milch
bring remove Schlagsahne
```

### Other

```bash
bring lists              # show all your Bring! lists
bring use Wochenmarkt    # switch default list
bring --refresh          # refresh the stored auth token
```
