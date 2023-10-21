# k8s change context

### Basic tool for rapid switching and show k8s context
```Basic function:```
```
  -h  --help                Print help information
  -c  --showCurrentContext  Show Current Context
  -s  --showListContext     Show Context List
  -n  --newContext          Set new context
```

in combination with rofi, you can quickly switch k8s context

## Example

#### Via command line
```bash
k8s-switch -s | wofi --dmenu --sort-order default | xargs $HOME/.local/bin/k8s-switch -n
```

#### Via Waybar
it will show current context and on click open wofi with context list you can search and select new context
```json
    "custom/kubernetes": {
        "format": "<span color='#5e81ac' font='15' >icon u like </span> {}",
        "interval": 3,
        "exec": "$HOME/.local/bin/k8s-switch -c",
        "on-click": "$HOME/.local/bin/k8s-switch -s | wofi --dmenu --sort-order default | xargs $HOME/.local/bin/k8s-switch -n",
        "tooltip": "false"
    },
```

#### Via sway shortcut
CTRL+SHIFT+m - show wofi with context list you can search and select new context
```bash
bindsym $mod+Shift+m exec $HOME/.local/bin/k8s-switch -s | wofi --dmenu --sort-order default | xargs $HOME/.local/bin/k8s-switch -n
```
