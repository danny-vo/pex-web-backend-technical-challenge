# Output directory
OUTDIR = ./out
OBJDIR = $(OUTDIR)/fibonacci-backend

build:
	mkdir -p $(OBJDIR)
	go build -o $(OBJDIR) ./...

clean:
	rm -rf $(OUTDIR)