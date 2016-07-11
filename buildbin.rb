#!/usr/bin/env ruby

pwd = `pwd`.strip

dist = './dist'
`mkdir -p #{dist}`

gopkg = 'github.com/CenturyLinkCloud/clc-sdk'

ospairs = [
  ['windows', 'amd64', '.exe'],
  ['linux', 'amd64', ''],
  ['darwin', 'amd64', '.osx'],
]

binaries = [
  ['natip', 'bin/natip/main.go'],
  ['baremetal-info', 'bin/natip/main.go'],
]

binaries.each do |basename, path|
  puts "== building #{basename}"
  ospairs.each do |os, arch, extension|
    cmd = "docker run --rm -it " +
          "-e GOOS=#{os} -e GOARCH=#{arch} " +
          "-v #{pwd}:/go/src/#{gopkg} -w /go/src/#{gopkg} " +
          "golang:1.7 go build -o #{dist}/#{basename}#{extension} #{path}"
    puts "BUILD: #{cmd}"
    puts `#{cmd}`
    raise "ERROR BUILDING #{basename}" unless $?.success?
  end
end
