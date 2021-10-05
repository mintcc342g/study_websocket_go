PACKAGE = study_websocket_go
BUILDPATH ?= $(CURDIR)
BASE	= $(BUILDPATH)
BIN		= $(BASE)/bin

UNAME := $(shell uname)
ifeq ($(UNAME), Linux)
	GOENV   ?= CGO_ENABLED=0 GOOS=linux
endif
GOBUILD = ${GOENV} go
GO      = go

BUILDTAG=-tags 'studyWebsocket'
export GO111MODULE=on

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: tidy build ; $(info $(M) building all steps… ) @ ## Build all steps


.PHONY: build
build: ; $(info $(M) building executable… ) @ ## Build program binary
	$Q cd $(BASE)/cmd && $(GOBUILD) build -i \
		$(BUILDTAG) \
		-o $(BIN)/$(PACKAGE)

.PHONY: tidy
tidy: ; $(info $(M) tidy executable… ) @ ## get packages
	$Q cd $(BASE) && go mod tidy
