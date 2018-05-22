# Listener Service

Listener service handles operations from transactions obtained in Horizon and emits events to several predefined targets (stdout is used at the moment as a test target) about:

```
KYC review request submit/resubmit/, approval/rejection handling
Deposit
Invest
Payment/payment-v2 send/receive
Withdraw
Someone referred someone
Referral of someone passed KYC
```

Needs a signer key to be specified in PSIM config.yaml