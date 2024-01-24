NAME := goscript

BASE := $(shell pwd)
BINDIR := $(BASE)/bin

$(NAME): $(BINDIR)
	go build -o $(BINDIR)/$(NAME) $(BASE)/cmd/$(NAME)

$(BINDIR):
	@mkdir -p $(BINDIR)

