function c() {
    if [ -z ${1} ]
    then
        dir=$(shellbuddy -path 3>&1 1>&2 2>&3)
    else
        dir=$(shellbuddy -path -search "${1}" 3>&1 1>&2 2>&3)
    fi

    if [ $? -eq 0 ]
    then
        cd ${dir}
    fi
}

function h() {
    if [ -z ${1} ]
    then
        cmd=$(shellbuddy -cmd 3>&1 1>&2 2>&3)
    else
        cmd=$(shellbuddy -cmd -search "${1}" 3>&1 1>&2 2>&3)
    fi

    if [ $? -eq 0 ]
    then
        eval ${cmd}
        print -s ${cmd}
    fi
}

PROMPT=${PROMPT}'$(shellbuddy -cmd -path -add)'
