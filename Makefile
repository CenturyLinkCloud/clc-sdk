.PHONY : test uats
test: 
	godep go test -v ./sdk/...
uats:
	BUNDLE_GEMFILE=spec/Gemfile bundle exec rspec spec/*.rb
