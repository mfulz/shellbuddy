function c() {
    if [ -z "${1}" ]
    then
        dir=$(shellbuddy -entries dirs 3>&1 1>&2 2>&3)
    else
        dir=$(shellbuddy -entries dirs -search "${1}" 3>&1 1>&2 2>&3)
    fi

    if [ $? -eq 0 ]
    then
        cd ${dir}
    fi
}

function h() {
    if [ -z "${1}" ]
    then
        cmd=$(shellbuddy -entries commands 3>&1 1>&2 2>&3)
    else
        cmd=$(shellbuddy -entries commands -search "${1}" 3>&1 1>&2 2>&3)
    fi

    if [ $? -eq 0 ]
    then
        eval ${cmd}
        print -s ${cmd}
    fi
}

function ci() {
    if [ -z "${1}" ]
    then
        shellbuddy -stdin -entries dirs
    else
        shellbuddy -stdin -entries dirs -search "${1}"
    fi
}

function hi() {
    if [ -z "${1}" ]
    then
        shellbuddy -stdin -entries commands
    else
        shellbuddy -stdin -entries commands -search "${1}"
    fi
}

function cip() {
    if [ -z "${1}" ]
    then
        return
    fi

    if [ -z "${2}" ]
    then
        shellbuddy -stdinpre "${1}" -entries dirs
    else
        shellbuddy -stdinpre "${1}" -entries dirs -search "${2}"
    fi
}

function hip() {
    if [ -z "${1}" ]
    then
        return
    fi

    if [ -z "${2}" ]
    then
        shellbuddy -stdinpre "${1}" -entries commands
    else
        shellbuddy -stdinpre "${1}" -entries commands -search "${2}"
    fi
}

PROMPT=${PROMPT}'$(shellbuddy -add)'
