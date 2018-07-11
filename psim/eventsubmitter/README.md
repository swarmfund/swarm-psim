# Listener Service

Listener service handles operations from transactions obtained in Horizon and emits events to several predefined targets (stdout is used at the moment as a test target) about:

```
KYC review request submited/resubmited, approved/rejected
Funds deposited
Funds invested
Payment/payment-v2 sent/received
Funds withdrawn
User referred
Referred user passed KYC
```

Needs a signer key to be specified in PSIM config.yaml