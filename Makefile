build:
	go build -o repo-installer

install:
	sudo cp repo-installer /usr/local/bin/repo-installer

uninstall:
	sudo rm -f /usr/local/bin/repo-installer

