# TWRPwipe bash-completion script
_twrpwipe() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="-c -d -h -v"

    case "${prev}" in
        -c)
            COMPREPLY=()
            return 0
            ;;
        -d)
            COMPREPLY=()
            return 0
            ;;
        *)
            ;;
    esac

    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}

complete -F _twrpwipe twrpwipe
