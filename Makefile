.PHONY : test uats
test: 
	godep go test -v ./clc/...
uats:
	BUNDLE_GEMFILE=spec/Gemfile bundle exec rspec spec/*.rb
