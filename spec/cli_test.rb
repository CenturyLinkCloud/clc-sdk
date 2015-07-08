require 'json'
require 'singleton'

describe 'clc cli' do
  before :all do
    `godep go build -o spec/cli ./cli`
    @server = Server.instance
  end

  it 'creates a server' do
    name = 'sample'
    json = JSON.parse(`./spec/cli server create -n #{name} -c 1 -m 1 -t standard -g 8aa8153a1ba24542908155da468bb71a -s UBUNTU-14-64-TEMPLATE`)

    @server.uuid = get_id 'self', json
    id = get_id 'status', json

    status = JSON.parse(`./spec/cli status get #{id}`)['status']
    until status == 'succeeded' || status == 'failed' do
      status = JSON.parse(`./spec/cli status get #{id}`)['status']
      sleep 20
    end
    expect(status).to eq('succeeded')
  end

  def get_id(rel, json)
    json['links'].select{ |val| val['rel'] == rel }.flat_map{ |val| val['id'] }[0]
  end

  it 'fetches a known server' do
    json = JSON.parse(`./spec/cli server get #{@server.uuid}`)

    expect(json['groupId']).to eq('8aa8153a1ba24542908155da468bb71a')
  end

end

class Server
  include Singleton
  attr_accessor :uuid
end
