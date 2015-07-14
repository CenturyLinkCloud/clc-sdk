require 'json'
require 'singleton'

describe 'clc cli' do
  before :all do
    Cli.build
  end

  it 'lists all policies' do
    policies = Cli.get_all_aa_policies
    
    expect(policies['links'][0]['href']).to eq("/v2/antiAffinityPolicies/#{ENV["CLC_ALIAS"]}")
  end

  it 'gets a specific policy' do
    policies = Cli.get_all_aa_policies
    expectedPolicy = policies['items'][0]
    policy = Cli.get_aa_policy(expectedPolicy['id'])

    expect(policy['name']).to eq(expectedPolicy['name'])
  end

  it 'creates an aa policy' do
    name = 'sample'
    policy = Cli.create_policy(name, 'va1')

    expect(policy['name']).to eq(name)
    Cli.delete_policy(policy['id'])
  end

  it 'deletes an aa policy' do
    policy = Cli.create_policy('sample', 'va1') 
    msg = Cli.delete_policy(policy['id'])

    expect(msg).to eq("deleted aa policy: #{policy['id']}\n")
  end

  it 'creates, fetches and deletes a server' do
    server = Cli.create_server('sample')

    expect(server.status).to eq('succeeded')

    json = Cli.get_server(server.uuid)

    expect(json['groupId']).to eq('8aa8153a1ba24542908155da468bb71a')
    expect(json['details']['cpu']).to eq(1)
    expect(json['details']['memoryMB']).to eq(1024)

    delete = Cli.delete_server(json['id'])

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

  def self.get_all_aa_policies
    JSON.parse(`./spec/clc aa get`)
  end

  def self.get_aa_policy(id)
    JSON.parse(`./spec/clc aa get #{id}`)
  end

  def self.create_policy(name, loc)
    JSON.parse(`./spec/clc aa c -n #{name} -l #{loc}`)
  end

  def self.delete_policy(id)
    `./spec/clc aa d #{id}`
  end

  def self.create_server(name)
    json = JSON.parse(`./spec/clc server create -n #{name} -c 1 -m 1 -t standard -g 8aa8153a1ba24542908155da468bb71a -s UBUNTU-14-64-TEMPLATE`)
    id = json['links'].select{ |val| val['rel'] == 'status' }.flat_map{ |val| val['id'] }[0]
    uuid = json['links'].select{ |val| val['rel'] == 'self' }.flat_map{ |val| val['id'] }[0]
    server = Server.new(uuid)

    status = JSON.parse(`./spec/clc status get #{id}`)['status']
    until status == 'succeeded' || status == 'failed' do
      status = JSON.parse(`./spec/clc status get #{id}`)['status']
      sleep 20
    end
    server.status = status
    server
  end

  def self.get_server(id)
    JSON.parse(`./spec/clc server get #{id}`)
  end

  def self.delete_server(name)
    JSON.parse(`./spec/clc server delete #{name}`)
  end
end
