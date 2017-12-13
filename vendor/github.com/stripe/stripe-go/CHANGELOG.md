# Changelog

## 28.6.1 - 2017-11-02
* [#492](https://github.com/stripe/stripe-go/pull/492) Correct name of user agent header used to send Go version to Stripe's API

## 28.6.0 - 2017-10-31
* [#491](https://github.com/stripe/stripe-go/pull/491) Support for exchange rates APIs

## 28.5.0 - 2017-10-27
* [#488](https://github.com/stripe/stripe-go/pull/488) Support for listing source transactions

## 28.4.2 - 2017-10-25
* [#486](https://github.com/stripe/stripe-go/pull/486) Send the required `object=bank_account` parameter when adding a bank account through an account
* [#487](https://github.com/stripe/stripe-go/pull/487) Make bank account's `account_holder_name` and `account_holder_type` parameters truly optional

## 28.4.1 - 2017-10-24
* [#484](https://github.com/stripe/stripe-go/pull/484) Error early when params not specified for card-related API calls

## 28.4.0 - 2017-10-19
* [#477](https://github.com/stripe/stripe-go/pull/477) Support context on API requests with `Params.Context` and `ListParams.Context`

## 28.3.2 - 2017-10-19
* [#479](https://github.com/stripe/stripe-go/pull/479) Pass token in only one of `external_account` *or* source when appending card

## 28.3.1 - 2017-10-17
* [#476](https://github.com/stripe/stripe-go/pull/476) Make initializing new backends concurrency-safe

## 28.3.0 - 2017-10-10
* [#359](https://github.com/stripe/stripe-go/pull/359) Add support for verify sources (added `Values` on `SourceVerifyParams`)

## 28.2.0 - 2017-10-09
* [#472](https://github.com/stripe/stripe-go/pull/472) Add support for `statement_descriptor` in source objects
* [#473](https://github.com/stripe/stripe-go/pull/473) Add support for detaching sources from customers

## 28.1.0 - 2017-10-05
* [#471](https://github.com/stripe/stripe-go/pull/471) Add support for `RedirectFlow.FailureReason` for sources

## 28.0.1 - 2017-10-03
* [#468](https://github.com/stripe/stripe-go/pull/468) Fix encoding of pointer-based scalars (e.g. `Active *bool` in `Product`)
* [#470](https://github.com/stripe/stripe-go/pull/470) Fix concurrent race in `form` package's encoding caches

## 28.0.0 - 2017-09-27
* [#467](https://github.com/stripe/stripe-go/pull/467) Change `Product.Get` to include `ProductParams` for request metadata
* [#467](https://github.com/stripe/stripe-go/pull/467) Fix sending extra parameters on product and SKU requests

## 27.0.2 - 2017-09-26
* [#465](https://github.com/stripe/stripe-go/pull/465) Fix encoding of `CVC` parameter in `CardParams`

## 27.0.1 - 2017-09-20
* [#461](https://github.com/stripe/stripe-go/pull/461) Fix encoding of `TypeData` under sources

## 27.0.0 - 2017-09-19
* [#458](https://github.com/stripe/stripe-go/pull/458) Remove `ChargeParams.Token` (this seems like it was added accidentally)

## 26.0.0 - 2017-09-17
* Introduce `form` package so it's no longer necessary to build conditional structures to encode parameters -- this may result in parameters that were set but previously not encoded to now be encoded so **PLEASE TEST CAREFULLY WHEN UPGRADING**!
* Alphabetize all struct fields -- this may result in position-based struct initialization to fail if it was being used
* Switch to stripe-mock for testing (test suite now runs completely!)
* Remote Displayer interface and Display implementations
* Add `FraudDetails` to `ChargeParams`
* Remove `FraudReport` from `ChargeParams` (use `FraudDetails` instead)

## 25.2.0 - 2017-09-13
* Add `OnBehalfOf` to charge parameters.
* Add `OnBehalfOf` to subscription parameters.

## 25.1.0 - 2017-09-06
* Use bearer token authentication for API requests

## 25.0.0 - 2017-08-21
* All `Del` methods now take params as second argument (which may be `nil`)
* Product `Delete` has been renamed to `Del` for consistency
* Product `Delete` now returns `(*Product, error)` for consistency
* SKU `Delete` has been renamed to `Del` for consistency
* SKU `Delete` now returns `(*SKU, error)` for consistency

## 24.3.0 - 2017-08-08
* Add `FeeZero` to invoice and `TaxPercentZero` to subscription for zeroing values

## 24.2.0 - 2017-07-25
* Add "range queries" for supported parameters (e.g. `created[gte]=123`)

## 24.1.0 - 2017-07-17
* Add metadata to subscription items

## 24.0.0 - 2017-06-27
	`Pay` on invoice now takes specific pay parameters

## 23.2.1 - 2017-06-26
* Fix bank account retrieval when using a customer ID

## 23.2.0 - 2017-06-26
* Support sharing path while creating a source

## 23.1.0 - 2017-06-26
* Add LoginLinks to client list

## 23.0.0 - 2017-06-23
	plan.Del now takes `stripe.PlanParams` as a second argument

## 22.6.0 - 2017-06-19
* Support for ephemeral keys

## 22.5.0 - 2017-06-15
* Support for checking webhook signatures

## 22.4.1 - 2017-06-15
* Fix returned type of subscription items list
* Note: I meant to release this as 22.3.1, but I'm leaving it as it was released

## 22.3.0 - 2017-06-14
* Fix parameters for subscription items list

## 22.2.0 - 2017-06-13
* Support subscription items when getting upcoming invoice
* Support setting subscription's quantity to zero when getting upcoming invoice

## 22.1.1 - 2017-06-12
* Handle `deleted` parameter when updating subscription items in a subscription

## 22.1.0 - 2017-05-25
* Change `Logger` to a `log.Logger`-like interface so other loggers are usable

## 22.0.0 - 2017-05-25
* Add support for login links
* Add support for new `Type` for accounts
* Make `Event` `Request` (renamed from `Req`) a struct with a new idempotency key
* Rename `Event` `UserID` to `Account`

## 21.5.1 - 2017-05-23
* Fix plan update so `TrialPeriod` parameter is sent

## 21.5.0 - 2017-05-15
* Implement `Get` for `RequestValues`

## 21.4.1 - 2017-05-11
* Pass extra parameters to API calls on bank account deletion

## 21.4.0 - 2017-05-04
* Add `Billing` and `DueDate` filters to invoice listing
* Add `Billing` filter to subscription listing

## 21.3.0 - 2017-05-02
* Add `DetailsCode` to `IdentityVerification`

## 21.2.0 - 2017-04-19
* Send user agent information with `X-Stripe-Client-User-Agent`
* Add `stripe.SetAppInfo` for plugin authors to register app information

## 21.1.0 - 2017-04-12
* Allow coupon to be specified when creating orders
* No longer require that items have descriptions when creating orders

## 21.0.0 - 2017-04-07
* Balances are now retrieved by payout instead of by transfer

## 20.0.0 - 2017-04-06
* Bump API version to 2017-04-06: https://stripe.com/docs/upgrades#2017-04-06
* Add support for payouts and recipient transfers
* Change the transfer resource to support its new format
* Deprecate recipient creation
* Disputes under charges are now expandable and collapsed by default
* Rules under charge outcomes are now expandable and collapsed by default

## 19.17.0 - 2017-04-06
* Please see 20.0.0 (bad release)	

## 19.16.0 - 2017-03-23
* Allow the ID of an identity document to be passed into an account owner update

## 19.15.0 - 2017-03-22
* Add `ShippingCarrier` to dispute evidence

## 19.14.0 - 2017-03-20
* Add `Period`, `Plan`, and `Quantity` to `InvoiceItem`

## 19.13.0 - 2017-03-20
* Add `AdditionalOwnersEmpty` to allow additional owners to be unset

## 19.12.0 - 2017-03-17
* Add new form of file upload using `io.FileReader` and filename

## 19.11.0 - 2017-03-13
* Add `Token` to `SourceObjectParams`

## 19.10.0 - 2017-03-13
* Add `CouponEmpty` (allowing a coupon to be cleared) to customer parameters
* Add `CouponEmpty` (allowing a coupon to be cleared) to subscription parameters

## 19.9.0 - 2017-03-08
* Add missing value "all" to subscription statuses

## 19.8.0 - 2017-03-02
* Add subscription items client to main `client.API` struct

## 19.7.0 - 2017-03-01
* Add `Statement` (statement descriptor) to `CaptureParams`

## 19.6.0 - 2017-02-22
* Add new parameters for invoices and subscriptions

## 19.5.0 - 2017-02-13
* Add new rich `Destination` type to `ChargeParams`

## 19.4.0 - 2017-02-03
* Support Connect account as payment source

## 19.3.0 - 2017-02-02
* Add transfer group to charges and transfers

## 19.2.0 - 2017-01-23
* Add `Rule` to `ChargeOutcome`

## 19.1.0 - 2017-01-18
* Add support for updating sources

## 19.0.2 - 2017-01-04
* Fix subscription `trial_period_days` to be populated by the right value

## 19.0.1 - 2016-12-08
* Include verification document details when persisting `LegalEntity`

## 19.0.0 - 2016-12-07
* Remote `SubProrationDateNow` field from `InvoiceParams`

## 18.14.1 - 2016-12-05
* Truncate `tax_percent` at four decimals (e.g. 3.9750%) instead of two

## 18.14.0 - 2016-11-23
* Add retrieve method for 3-D Secure resources

## 18.13.0 - 2016-11-15
* Add `PaymentSource` to `API`

## 18.12.0 - 2016-11-14
* Allow bank accounts to be created as a customer source

## 18.11.0 - 2016-11-14
* Add `TrialPeriodEnd` to `SubParams`

## 18.10.0 - 2016-11-09
* Add `StatusTransitions` to `Order`

## 18.9.0 - 2016-11-04
* Add `Application` to `Charge`

## 18.8.0 - 2016-10-24
* Add `Review` to `Charge` for the charge reviews

## 18.7.0 - 2016-10-18
* Add `RiskLevel` to `ChargeOutcome`

## 18.6.0 - 2016-10-18
* Support for 403 status codes (permission denied)

## 18.5.0 - 2016-10-18
* Add `Status` to `SubListParams` to allow filtering subscriptions by status

## 18.4.0 - 2016-10-14
* Add `HasEvidence` and `PastDue` to `EvidenceDetails`

## 18.3.0 - 2016-10-10
* Add `NoDiscountable` to `InvoiceItemParams`

## 18.2.0 - 2016-10-10
* Add `BusinessLogo` to `Account`
* Add `ReceiptNumber` to `Charge`
* Add `DestPayment` to `Transfer`

## 18.1.0 - 2016-10-04
* Support for Apple Pay domains

## 18.0.0 - 2016-10-03
* Support for subscription items
* Correct `SourceTx` on `Transfer` to be a `SourceTransaction`
* Change `Charge` on `Resource` to be expandable (now a struct instead of string)

## 17.5.0 - 2016-09-22
* Support customer-related operations for bank accounts

## 17.4.2 - 2016-09-19
* Fix but where some parameters were not being included on order update

## 17.4.1 - 2016-09-15
* Fix bug that required a date of birth to be included on account update

## 17.4.0 - 2016-09-13
* Add missing Kana and Kanji address and name fields to account's legal entity
* Add `ReceiptNumber` and `Status` to `Refund`

## 17.3.0 - 2016-09-07
* Add support for sources endpoint

## 17.2.0 - 2016-08-29
* Add order returns to `API`

## 17.1.0 - 2016-08-22
* Add `DeactiveOn` to `Product`

## 17.0.0 - 2016-08-18
* Allow expansion of destination on transfers
* Allow expansion of sources on balance transactions

## 16.8.0 - 2016-08-17
* Add `OriginatingTransaction` to `Fee`

## 16.7.1 - 2016-08-17
* Allow params to be nil when retrieving a refund

## 16.7.0 - 2016-08-11
* Add support for 3-D Secure

## 16.6.0 - 2016-08-09
* Add `ReceiptNumber` to `Invoice`

## 16.5.0 - 2016-08-08
* Add `Meta` to `Account`

## 16.4.0 - 2016-08-05
* Allow the migration of recipients to accounts
* Add `MigratedTo` to `Recipient`

## 16.3.1 - 2016-07-25
* URL-escape the IDs of coupons and plans when making API requests

## 16.3.0 - 2016-07-19
* Add `NoClosed` to `InvoiceParams` to allow an invoice to be reopened

## 16.2.1 - 2016-07-11
* Consider `SubParams.QuantityZero` when updating a subscription

## 16.2.0 - 2016-07-07
* Upgrade API version to 2016-07-06

## 16.1.0 - 2016-07-07
* Add `Returns` field to `Order`

## 16.0.0 - 2016-06-30
* Remove `Name` field on `SKU`; it's not actually supported
* Support updating `Product` on `SKU`

## 15.6.0 - 2016-06-24
* Allow product and SKU attributes to be updated

## 15.5.0 - 2016-06-24
* Add `TaxPercent` and `TaxPercentZero` to `CustomerParams`

## 15.4.0 - 2016-06-20
* Add `TokenizationMethod` to `Card` struct

## 15.3.0 - 2016-06-15
* Add `BalanceZero` to `CustomerParams` so that balance can be zeroed out

## 15.2.0 - 2016-06-03
* Add `ToValues` to `RequestValues` struct

## 15.1.0 - 2016-05-26
* Add `BusinessVatID` to customer creation parameters

## 15.0.0 - 2016-05-24
* Fix handling of nested objects in arrays in request parameters

## 14.4.0 - 2016-05-24
* Add granular error types in new `Err` field on `stripe.Error`

## 14.3.0 - 2016-05-20
* Allow Relay orders to be returned and add associated types

## 14.2.3 - 2016-05-20
* When creating a bank account token, only send routing number if it's been set

## 14.2.2 - 2016-05-17
* When creating a bank account, only send routing number if it's been set

## 14.2.1 - 2016-05-17
* Add missing SKU clinet to client API type

## 14.2.0 - 2016-05-11
* Add `Reversed` and `AmountReversed` fields to `Transfer`

## 14.1.0 - 2016-05-05
* Allow `default_for_currency` to be set when creating a card

## 14.0.0 - 2016-05-04
* Change the signature for `sub.Delete`. The customer ID is no longer required.

## 13.12.0 - 2016-04-28
* Add `Currency` to `Card`

## 13.11.1 - 2016-04-22
* Fix bug where new external accounts could not be marked default from token

## 13.11.0 - 2016-04-21
* Expose a number of list types that were previously internal (full list below)
* Expose `stripe.AccountList`
* Expose `stripe.TransactionList`
* Expose `stripe.BitcoinReceiverList`
* Expose `stripe.ChargeList`
* Expose `stripe.CountrySpecList`
* Expose `stripe.CouponList`
* Expose `stripe.CustomerList`
* Expose `stripe.DisputeList`
* Expose `stripe.EventList`
* Expose `stripe.FeeList`
* Expose `stripe.FileUploadList`
* Expose `stripe.InvoiceList`
* Expose `stripe.OrderList`
* Expose `stripe.ProductList`
* Expose `stripe.RecipientList`
* Expose `stripe.TransferList`
* Switch to use of `stripe.BitcoinTransactionList`
* Switch to use of `stripe.SKUList`

## 13.10.1 - 2016-04-20
* Add support for `TaxPercentZero` to invoice and subscription updates

## 13.10.0 - 2016-04-19
* Expose `stripe.PlanList` (previously an internal type)

## 13.9.0 - 2016-04-18
* Add `TaxPercentZero` struct to `InvoiceParams`
* Add `TaxPercentZero` to `SubParams`

## 13.8.0 - 2016-04-12
* Add `Outcome` struct to `Charge`

## 13.7.0 - 2016-04-06
* Add `Description`, `IIN`, and `Issuer` to `Card`

## 13.6.0 - 2016-04-05
* Add `SourceType` (and associated constants) to `Transfer`

## 13.5.0 - 2016-03-29
* Add `Meta` (metadata) to `BankAccount`

## 13.4.0 - 2016-03-29
* Add `Meta` (metadata) to `Card`

## 13.3.0 - 2016-03-29
* Add `DefaultCurrency` to `CountrySpec`

## 13.2.0 - 2016-03-18
* Add `SourceTransfer` to `Charge`
* Add `SourceTx` to `Transfer`

## 13.1.0 - 2016-03-15
* Add `Reject` on `Account` to support the new API feature

## 13.0.0 - 2016-03-15
* Upgrade API version to 2016-03-07
* Remove `Account.BankAccounts` in favor of `ExternalAccounts`
* Remove `Account.Currencies` in favor of `CountrySpec`

## 12.1.0 - 2016-02-04
* Add `ListParams.StripeAccount` for making list calls on behalf of connected accounts
* Add `Params.StripeAccount` for symmetry with `ListParams.StripeAccount`
* Deprecate `Params.Account` in favor of `Params.StripeAccount`

## 12.0.0 - 2016-02-02
* Add support for fetching events for managed accounts (`event.Get` now takes `Params`)

## 11.5.0 - 2016-02-26
* Allow a `PII.PersonalIDNumber` number to be used to create a token

## 11.4.0 - 2016-02-24
* Add missing subscription fields to `InvoiceParams` for use with `invoice.GetNext`

## 11.3.0 - 2016-02-19
* Add `AccountHolderName` and `AccountHolderType` to bank accounts

## 11.2.0 - 2016-02-11
* Add support for `CountrySpec`
* Add `SSNProvided`, `PersonalIDProvided` and `BusinessTaxIDProvided` to `LegalEntity`

## 11.1.2 - 2016-02-02
* Fix card update method to correctly take expiration date

## 11.1.1 - 2016-02-01
* Fix recipient update so that it can take a bank token (like create)

## 11.0.1 - 2016-01-11
* Add missing field `country` to shipping details of `Charge` and `Customer`

## 11.0.0 - 2016-01-07
* Add missing field `Default` to `BankAccount`
* Add `OrderParams` parameter to `Order` retrieval
* Fix parameter bug when creating a new `Order`
* Support special value of 'now' for trial end when updating subscriptions

## 10.3.0 - 2015-12-10
* Allow an account to be referenced when creating a card

## 10.2.0 - 2015-12-04
* Add `Update` function on `Coupon` client so that metadata can be set

## 10.1.0 - 2015-12-01
* Add a verification routine for external accounts

## 10.0.0 - 2015-11-30
* Return models along with `error` when deleting resources with `Del`
* Fix bug where country parameter wasn't included for some account creation

## 9.0.0 - 2015-11-13
* Return model (`Sub`) when cancelling a subscription (`sub.Cancel`)

## 8.0.0 - 2015-08-17
* Add ability to list and retrieve refunds without a Charge

## 7.0.0 - 2015-08-03
* Add ability to list and retrieve disputes

## 6.8.0 - 2015-07-29
* Add ability to delete an account

## 6.7.1 - 2015-07-17
* Bug fixes

## 6.7.0 - 2015-07-16
* Expand logging object
* Move proration date to subscription update
* Send country when creating/updating account

## 6.6.0 - 2015-07-06
* Add request ID to errors

## 6.5.0 - 2015-07-06
* Update bank account creation API
* Add destination, application fee, transfer to Charge struct
* Add missing fields to invoice line item
* Rename deprecated customer param value

## 6.4.2 - 2015-06-23
* Add BusinessUrl, BusinessUrl, BusinessPrimaryColor, SupportEmail, and
* SupportUrl to Account.

## 6.4.1 - 2015-06-16
* Change card.dynamic_last_four to card.dynamic_last4

## 6.4.0 - 2015-05-28
* Rename customer.default_card -> default_source

## 6.3.0 - 2015-05-19
* Add shipping address to charges
* Expose card.dynamic_last_four
* Expose account.tos_acceptance
* Bug fixes
* Bump API version to most recent one

## 6.2.0 - 2015-04-09
* Bug fixes
* Add Extra to parameters

## 6.1.0 - 2015-03-17
* Add TaxPercent for subscriptions
* Event bug fixes

## 6.0.0 - 2015-03-15
* Add more operations for /accounts endpoint
* Add /transfers/reversals endpoint
* Add /accounts/bank_accounts endpoint
* Add support for Stripe-Account header

## 5.1.0 - 2015-02-25
* Add new dispute status `warning_closed`
* Add SubParams.TrialEndNow to support `trial_end = "now"`

## 5.0.1 - 2015-02-25
* Fix URL for upcoming invoices

## 5.0.0 - 2015-02-19
* Bump to API version 2014-02-18
* Change Card, DefaultCard, Cards to Source, DefaultSource, Sources in Stripe response objects
* Add paymentsource package for manipulating Customer's sources
* Support Update action for Bitcoin Receivers

## 4.4.3 - 2015-02-08
* Modify NewIdempotencyKey() algorithm to increase likelihood of randomness

## 4.4.2 - 2015-01-24
* Add BankAccountParams.Token
* Add Token.ClientIP
* Add LogLevel

## 4.4.0 - 2015-01-20
* Add Bitcoin support

## 4.3.0 - 2015-01-13
* Added support for listing FileUploads
* Mime parameter on FileUpload has been changed to Type

## 4.2.1 - 2014-12-28
* Handle charges with customer card tokens

## 4.2.0 - 2014-12-18
* Add idempotency support

## 4.1.0 - 2014-12-17
* Bump to API version 2014-12-17.

## 4.0.0 - 2014-12-16
* Add FileUpload resource. This brings in a new endpoint (uploads.stripe.com) and thus makes changes to some of the existing interfaces.
* This also adds support for multipart content.

## 3.1.0 - 2014-12-16
* Add Charge.FraudDetails

## 3.0.1 - 2014-12-15
* Add timeout value to HTTP requests

## 3.0.0 - 2014-12-05
* Add Dispute.EvidenceDetails
* Remove Dispute.DueDate
* Change Dispute.Evidence from string to struct

## 2.0.0 - 2014-11-26
* Change List interface to .Next() and .Resource()
* Better error messages for Get() methods
* EventData.Raw contains the raw event message
* SubParams.QuantityZero can be used for free subscriptions

## 1.0.3 - 2014-10-22
* Add AddMeta method

## 1.0.2 - 2014-09-23
* Minor fixes

## 1.0.1 - 2014-09-23
* Linter-based updates

## 1.0.0 - 2014-09-22
* Initial version
