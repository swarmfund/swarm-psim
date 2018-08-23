# Changelog

All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## Deposits with tasks

### Deprecated

- `bitcoin.GetNetParams` in favor of `derive.NetworkParams`

### Changed

- btc_deposit/btc_deposit_verify config update:
    `offchain_currency` and `offchain_blockchain` removed in favor of `network_type`
    `external_system` is now optional, by default `deposit_asset` external system type is used