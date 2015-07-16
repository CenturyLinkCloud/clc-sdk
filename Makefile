VERSION=0.1

.PHONY : test uats
test: 
	godep go test -v ./sdk/...
uats:
	BUNDLE_GEMFILE=spec/Gemfile bundle exec rspec spec/*.rb
build:
	godep go build -o clc-$(VERSION) ./clc
deploy: build
	s3cmd -c ~/.s3cfgs/s3cfg_slos -P put clc-$(VERSION) s3://clc-cli/$(VERSION)/clc
	rm clc-$(VERSION)
