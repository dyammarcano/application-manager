# Application structure for service management

## Todo

- [x] Add file watcher to reload config on change
- [ ] Add generate inline config
- [x] Add logger
- [ ] Add self update to update service
- [ ] Add service status to check if service is running
- [ ] Add service restart to restart service
- [ ] Add generate script to run service

## Usage

```bash
# Start service
./service start or ./service start -c config.json --logDir /var/log

# Stop service
./service stop

# Restart service
./service restart

# Update service
./service update

# Status service
./service status

# Generate script
./service generate [service name] --config config.json or --config-string 
```

## Config

```json
{
  "name": "service",
  "description": "Service description",
  "version": "1.0.0",
  "author": "Author",
  "license": "MIT"
}
```
