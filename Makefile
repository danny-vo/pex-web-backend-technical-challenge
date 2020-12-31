# Output directory
OUTDIR = ./out
OBJDIR = $(OUTDIR)/fibonacci-backend

build:
	mkdir -p $(OBJDIR)
	go build -o $(OBJDIR) ./...

run: build
	$(OBJDIR)/fibonacci_server

clean:
	rm -rf $(OUTDIR)