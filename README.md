# PSIM services

## Common configuration

```yaml
discovery:
  host: localhost
  port: 8500
  
horizon:
	addr: http://localhost:8000
	seed: # exchange kp seed

skrill:
	merchant: # skrill's merchant email
	password: # skrill's merchant password
```



## Forfeit listener

```yaml
forfeit_listener:
	disable: false
	leadership_key: service/forfeit_listener/leader
```

Listens on Horizon's `GET forfeit_requests?exchange=...` , executes Skrill send money request to email designated in forfeit request, sends RFR op on success.

Does not voluntarily release leadership once acquired, will spam logs with errors in case of failures.

## Supervisor

```yaml
supervisor:
	disable: false
	cursor: 2017-06-09 # date to start with
	leadership_key: service/supervisor/leader
	
```

Goes through Skrill transaction history and checks if corresponding CER is already processed, otherwise creates new one and ask neighbor to verify and submit it.

Does not voluntarily release leadership once acquired, will spam logs with errors in case of failures.

## Neighbor

```yaml
neighbor:
	disable: false
	host: localhost
	port: 0 # random
	service: neighbor
```

Listens for requests to verify and sign CER, submits if all OK.

## Rate sync

```yaml
rate_sync:
	disable: false
	service: rate_sync
	leadership_key: service/<service>/leader
	host: localhost
	port: 0 # random
```

Uses leadership key to sync other rate sync services.



## Skrill IPN

```yaml
skrill_ipn:
	disable: false
	service: skrill_ipn
	host: localhost
	port: 0 # random
	public_addr: # addr skrill will try to reach us on
```

Listens for Skrill IPN notifications, create CER and asks neighbor to verify and submit it.


# Contribution
1. Install pre-commit-hooks `pip install pre-commit`; Run `pre-commit install` to install `pre-commit` into your git hooks.