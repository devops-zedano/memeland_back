#!/usr/bin/env ruby
# frozen_string_literal: true

# This is the main task and will be executed.
def main
  action = ARGV[0]
  resource = ARGV[1]

  exec_task action, resource
end

# This is the entry point for the actions on a resource.
def exec_task(action = nil, resource = nil)
  if action.nil?
    puts 'Falling back to DEFAULT action...'
    default_action
  end

  case action
  when 'b', 'build'
    if !resource.nil?
      build_function_action resource
    else
      default_action if resource.nil?
    end
  when 'remove', 'r'
    remove_action
  end
end

# This is what is going to be a default action
def default_action
  build_all
end

def build_all
  puts 'Executing action "build_all"...'

  system('
    . ./functions.sh
    build_all
    deploy_all
  ')
end

def build_function_action(func_name)
  puts 'Executing action "build_function"...'

  sanitized_func_name = func_name.gsub('/', '')
  system("
        . ./functions.sh
        build_func '#{sanitized_func_name}' \"$(pwd)\"
        deploy_func '#{sanitized_func_name}'
  ")
end

def remove_action
  system('
    . ./functions.sh
    remove
  ')
end

main
