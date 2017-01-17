default:
	@export GOPATH=$$(pwd) && go install checker
run: default
	@bin/checker
	@echo ""
