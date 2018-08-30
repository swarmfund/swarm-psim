# Masternode payout module

At each block module pays out fixed SWM amount to one of eligible addresses.

Address is considered eligible if following criteria are met:
    * any amount of SWM was burned to configured address
    * address continuously maintains configured SWM balance
    * configured number of blocks have passed since burning

Addresses are FIFO queued and dropped from the queue if balance falls below threshold.

