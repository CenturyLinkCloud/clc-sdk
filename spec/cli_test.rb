require 'json'
require 'singleton'

describe 'clc cli' do
  before :all do
    Cli.build
  end

  it 'creates, fetches and deletes a server' do
    server = Cli.create('sample')

    expect(server.status).to eq('succeeded')

    json = Cli.get(server.uuid)

    expect(json['groupId']).to eq('8aa8153a1ba24542908155da468bb71a')
    expect(json['details']['cpu']).to eq(1)
    expect(json['details']['memoryMB']).to eq(1024)

    delete = Cli.delete(json['id'])

    expect(delete['server']).to eq(json['id'])
  end
end

class Server
  attr_accessor :name, :uuid, :status

  def initialize(uuid)
    @uuid = uuid
  end
end


class Cli
  def self.build
    `godep go build -o spec/clc ./clc`
  end

  def self.create(name)
    json = JSON.parse(`./spec/clc server create -n #{name} -c 1 -m 1 -t standard -g 8aa8153a1ba24542908155da468bb71a -s UBUNTU-14-64-TEMPLATE`)
    id = json['links'].select{ |val| val['rel'] == 'status' }.flat_map{ |val| val['id'] }[0]
    uuid = json['links'].select{ |val| val['rel'] == 'self' }.flat_map{ |val| val['id'] }[0]
    server = Server.new(uuid)

    status = JSON.parse(`./spec/cli status get #{id}`)['status']
    until status == 'succeeded' || status == 'failed' do
      status = JSON.parse(`./spec/cli status get #{id}`)['status']
      sleep 20
    end
    server.status = status
    server
  end

  def self.get(id)
    JSON.parse(`./spec/clc server get #{id}`)
  end

  def self.delete(name)
    JSON.parse(`./spec/clc server delete #{name}`)
  end
end
