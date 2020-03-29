#!/usr/bin/env ruby

# This is the main task and will be executed.
def main
    action = ARGV[0]
    resource = ARGV[1]

    exec_task action, resource
end


# This is the entry point for the actions on a resource.
def exec_task(action = nil, resource = nil)
    if action == nil
        puts 'Falling back to DEFAULT action...'
        default_action
    end


end

# This is what is going to be a default action
def default_action
    build_all
end

def build_all
    puts 'Executing action build_all...'

    system('
    set -e

    cwd=$(pwd)
    mkdir -p "$cwd/bin"

    for a_func in $(ls -d */ | grep -v bin); do
        cd $a_func
        echo "Building \\"$a_func\\"..."
        go build -o "${cwd}/bin"
        cd "$cwd"
    done

    cd $cwd
    ')

    exit_code = $?
    if exit_code != 0
        puts "Error building the functions! Output is #{exit_code}"
        exit exit_code
    end

    puts 'The functions have been built...'

    puts 'Deploying the functions...'
    system('serverless deploy; rm -rf bin/*')
end

main