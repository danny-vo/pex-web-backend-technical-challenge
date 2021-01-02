# Output directory
OUTDIR = ./out
OBJDIR = $(OUTDIR)/fibonacci-backend
BINARY = fibonacci_server

build:
	mkdir -p $(OBJDIR)
	go build -o $(OBJDIR) ./...

test:
	go test -race ./...

run: build
	$(OBJDIR)/$(BINARY)

clean:
	rm -rf $(OUTDIR)