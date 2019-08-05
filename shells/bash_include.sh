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
        history -s ${cmd}
    fi
}

PS1=${PS1}'$(history -w && shellbuddy -cmd -path -add)'