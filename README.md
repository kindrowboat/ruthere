# ruthere
Checks that the sites you care about are up, and emails you over local postfix if they're not.

## installation
One a machine using `systemd`, you can install `ruthere` as a user level service.


``` bash
# compile and install
go build ./cmd/ruthere
cp ruthere ~/bin

# configure
mkdir -p ~/.config/ruthere
cp config.example.yml ~/.config/ruthere/config.yml
nano ~/.config/ruthere/config.yml # set the sites you want to monitor

# install, setup, and start service
mkdir -p ~/.config/systemd/user/
cp ruthere.service ~/.config/systemd/user
systemctl --user enable ruthere.service
systemctl --user start ruthere.service
systemctl --user status ruthere.service #make sure it start up properly
```
