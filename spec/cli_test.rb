require 'json'

describe 'clc cli' do
  before :all do
    `go build -o spec/cli ./cli`
  end

  it 'fetches servers' do
    server = 'va1t3bkapi01'
    json = JSON.parse(`./spec/cli server get #{server}`)

    expect(json['id']).to eq(server)
  end
end
