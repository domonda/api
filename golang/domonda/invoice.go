package domonda

import (
	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/date"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/domonda/go-types/uu"
	"github.com/domonda/go-types/vat"
)

// Invoice uploaded from a third party system to Domonda.
type Invoice struct {
	// Identifier of the system that produced the invoice and confirmes its values
	ConfirmedBy nullable.TrimmedString `json:"confirmedBy,omitempty"`

	// Name of the partner company
	PartnerName nullable.TrimmedString `json:"partnerName,omitempty"`
	// VAT ID of the partner company
	PartnerVatID vat.NullableID `json:"partnerVatId,omitempty"`
	// Company registration number of the partner company
	PartnerCompRegNo nullable.TrimmedString `json:"partnerCompRegNo,omitempty"`
	// Country of the partner company
	PartnerCountry country.NullableCode `json:"partnerCountry,omitempty"`

	// Number of the invoice
	InvoiceNumber nullable.TrimmedString `json:"invoiceNumber,omitempty"`
	// Internal number of the invoice
	InternalNumber nullable.TrimmedString `json:"internalNumber,omitempty"`
	// Date of the invoice
	InvoiceDate date.NullableDate `json:"invoiceDate,omitempty"`
	// Due date of the invoice
	DueDate date.NullableDate `json:"dueDate,omitempty"`

	// Number of the order
	OrderNumber nullable.TrimmedString `json:"orderNumber,omitempty"`
	// Date of the order
	OrderDate date.NullableDate `json:"orderDate,omitempty"`

	// Whether the invoice is a credit memo
	CreditMemo *bool `json:"creditMemo,omitempty"`

	// Net amount of the invoice (without VAT)
	Net *money.Amount `json:"net,omitempty"`
	// Total or gross amount of the invoice (including VAT)
	Total *money.Amount `json:"total,omitempty"`

	// Single VAT percentage of the invoice
	VATPercent *money.Rate `json:"vatPercent,omitempty"`
	// Multiple VAT percentages of the invoice
	VATPercentages nullable.FloatArray `json:"vatPercentages,omitempty"`
	// Multiple VAT amounts of the invoice, one per VAT percentage
	VATAmounts nullable.FloatArray `json:"vatAmounts,omitempty"`

	// Discount percentage of the invoice
	DiscountPercent *money.Rate `json:"discountPercent,omitempty"`
	// Date until which the discount is valid
	DiscountUntil date.NullableDate `json:"discountUntil,omitempty"`

	// Cost centers of the invoice
	CostCenters map[string]money.Amount `json:"costCenters,omitempty"`

	// Currency of the invoice
	Currency money.NullableCurrency `json:"currency,omitempty"`
	// Conversion rate of the currency
	ConversionRate *money.Rate `json:"conversionRate,omitempty"`
	// Date of the currency conversion rate
	ConversionRateDate date.NullableDate `json:"conversionRateDate,omitempty"`

	// Invoiced goods and services
	GoodsServices nullable.TrimmedString `json:"goodsServices,omitempty"`

	// Date from which the goods and services are delivered.
	// Use same date for from and until if the goods and services are delivered in a single day.
	DeliveredFrom date.NullableDate `json:"deliveredFrom,omitempty"`
	// Date until which the goods and services are delivered.
	// Use same date for from and until if the goods and services are delivered in a single day.
	DeliveredUntil date.NullableDate `json:"deliveredUntil,omitempty"`

	// Delivery note numbers related to the invoice
	DeliveryNoteNumbers []string `json:"deliveryNoteNumbers,omitempty"`

	// Invoice payment IBAN
	IBAN bank.NullableIBAN `json:"iban,omitempty"`
	// Invoice payment BIC
	BIC bank.NullableBIC `json:"bic,omitempty"`

	// Accounting items of the invoice
	AccountingItems []*AccountingItem `json:"accountingItems,omitempty"`
}

type AccountingItem struct {
	Title notnull.TrimmedString `json:"title"`

	GeneralLedgerAccountNumber account.Number `json:"generalLedgerAccountNumber"`

	BookingType string `json:"bookingType" jsonschema:"enum=DEBIT,enum=CREDIT"`
	AmountType  string `json:"amountType"  jsonschema:"enum=NET,enum=TOTAL"`

	Amount money.Amount `json:"amount"`

	ValueAddedTaxID               uu.NullableID `json:"valueAddedTax,omitempty"`
	ValueAddedTaxPercentageAmount *money.Amount `json:"valueAddedTaxPercentageAmount,omitempty"`
}
