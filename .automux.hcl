session = "icl"
# config = "./tmux.conf"
# attach_existing = false

window "Editor" {
    exec = "vim"
    focus = true
}

window "Shell" {
    split {
        #     vertical = true
        #     exec = "cmd_to_run_in_split"
        #     size = 30
        #     dir = "sub/"
    }
}
