package domonda

import (
	"errors"
	"fmt"
	"slices"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/date"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/domonda/go-types/strutil"
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
	// ISO 3166-1 alpha 2 country code of the partner company
	PartnerCountry country.NullableCode `json:"partnerCountry,omitempty"`
	// Number that identifies the partner company like a vendor or client number
	PartnerNumber nullable.TrimmedString `json:"partnerNumber,omitempty"`

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
	CostCenters map[string]money.Amount `json:"costCenters,omitempty" jsonschema:"oneof_type=object;null"`

	// Cost units of the invoice
	CostUnits map[string]money.Amount `json:"costUnits,omitempty" jsonschema:"oneof_type=object;null"`

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

// AccountingItem represents a single booking line in an invoice's accounting.
// Each item describes how the invoice amount should be posted to the general ledger.
type AccountingItem struct {
	// Title describing this accounting item
	Title notnull.TrimmedString `json:"title"`

	// GeneralLedgerAccountNumber where this item should be posted
	GeneralLedgerAccountNumber account.Number `json:"generalLedgerAccountNumber"`

	// BookingType indicates the side of the booking: "DEBIT" or "CREDIT"
	BookingType string `json:"bookingType" jsonschema:"enum=DEBIT,enum=CREDIT"`

	// AmountType indicates whether Amount is "NET" (without VAT) or "TOTAL" (with VAT)
	AmountType string `json:"amountType" jsonschema:"enum=NET,enum=TOTAL"`

	// Amount to be booked for this item
	Amount money.Amount `json:"amount"`

	// ValueAddedTaxID is the optional UUID of the VAT code to apply
	ValueAddedTaxID uu.NullableID `json:"valueAddedTax,omitempty"`

	// ValueAddedTaxPercentageAmount is the optional VAT percentage as amount (e.g., 20 for 20%)
	ValueAddedTaxPercentageAmount *money.Amount `json:"valueAddedTaxPercentageAmount,omitempty"`
}

// Validate checks if the invoice data is valid according to business rules.
// It validates dates, amounts, VAT percentages, currency codes, IBANs, BICs,
// and ensures consistency between related fields (e.g., total >= net, cost centers sum <= net).
//
// Returns nil if the invoice is valid, or an error describing validation failures.
func (inv *Invoice) Validate() error {
	if inv == nil {
		return errors.New("<nil> Invoice")
	}
	if err := inv.PartnerVatID.Validate(); err != nil {
		return fmt.Errorf("invalid partner VAT-ID: %w", err)
	}
	// if err := inv.PartnerCompanyID.Validate(); err != nil {
	// 	return fmt.Errorf("invalid partner company ID: %w", err)
	// }
	// if err := inv.PartnerCompanyLocationID.Validate(); err != nil {
	// 	return fmt.Errorf("invalid partner company location ID: %w", err)
	// }
	if err := inv.InvoiceDate.Validate(); err != nil {
		return fmt.Errorf("invalid invoice date: %w", err)
	}
	if err := inv.DueDate.Validate(); err != nil {
		return fmt.Errorf("invalid due date: %w", err)
	}
	if err := inv.OrderDate.Validate(); err != nil {
		return fmt.Errorf("invalid order date: %w", err)
	}
	if inv.Net != nil && *inv.Net < 0 {
		return errors.New("net amount must not be negative")
	}
	if inv.Total != nil && *inv.Total < 0 {
		return errors.New("total amount must not be negative")
	}
	if inv.Net != nil && inv.Total != nil && *inv.Total < *inv.Net {
		return fmt.Errorf("total amount %f must not be smaller than net %f", *inv.Total, *inv.Net)
	}
	if inv.VATPercent != nil && (*inv.VATPercent < 0 || *inv.VATPercent > 100) {
		return fmt.Errorf("vat percent %f not in range of [0..100]", *inv.VATPercent)
	}
	for i, percent := range inv.VATPercentages {
		if percent < 0 || percent > 100 {
			return fmt.Errorf("vat percentage[%d] %f not in range of [0..100]", i, percent)
		}
	}
	for i, amount := range inv.VATAmounts {
		if amount < 0 {
			return fmt.Errorf("vat amount[%d] %f must not be negative", i, amount)
		}
	}
	if inv.DiscountPercent != nil && (*inv.DiscountPercent < 0 || *inv.DiscountPercent > 100) {
		return fmt.Errorf("discount percent %f not in range of [0..100]", *inv.DiscountPercent)
	}
	if err := inv.DiscountUntil.Validate(); err != nil {
		return fmt.Errorf("invalid discount until date: %w", err)
	}
	if !inv.Currency.Valid() {
		return fmt.Errorf("invalid currency: %s", inv.Currency)
	}
	if inv.ConversionRate != nil && *inv.ConversionRate <= 0 {
		return fmt.Errorf("conversion rate must be greater zero, but is %f", *inv.ConversionRate)
	}
	if err := inv.ConversionRateDate.Validate(); err != nil {
		return fmt.Errorf("invalid conversion rate date: %w", err)
	}
	if err := inv.DeliveredFrom.Validate(); err != nil {
		return fmt.Errorf("invalid deliveredFrom date: %w", err)
	}
	if err := inv.DeliveredUntil.Validate(); err != nil {
		return fmt.Errorf("invalid deliveredUntil date: %w", err)
	}
	if inv.DeliveredFrom.IsNotNull() && inv.DeliveredUntil.IsNull() {
		return errors.New("deliveredFrom date needs deliveredUntil date to be provided too")
	}
	if inv.DeliveredFrom.IsNotNull() && inv.DeliveredUntil.IsNotNull() && inv.DeliveredFrom.After(inv.DeliveredUntil) {
		return fmt.Errorf("deliveredFrom date %s must not be after deliveredUntil date %s", inv.DeliveredFrom, inv.DeliveredUntil)
	}
	for i := range inv.DeliveryNoteNumbers {
		trimmed := strutil.TrimSpace(inv.DeliveryNoteNumbers[i])
		if trimmed == "" {
			inv.DeliveryNoteNumbers = slices.Delete(inv.DeliveryNoteNumbers, i, i+1)
			i--
		} else {
			inv.DeliveryNoteNumbers[i] = trimmed
		}
	}
	if err := inv.IBAN.Validate(); err != nil {
		return fmt.Errorf("invalid invoice IBAN: %w", err)
	}
	if err := inv.BIC.Validate(); err != nil {
		return fmt.Errorf("invalid invoice BIC: %w", err)
	}
	if err := inv.validateCostAmounts(inv.CostCenters, "costCenter", "cost center"); err != nil {
		return err
	}
	if err := inv.validateCostAmounts(inv.CostUnits, "costUnit", "cost unit"); err != nil {
		return err
	}
	return nil
}

// validateCostAmounts checks the net amounts of a cost center or cost unit
// mapping keyed by number: every number must be non-empty and every amount
// must be positive.
//
// The sum is compared against the net amount of the invoice, converted with
// ConversionRate when one is set. An invoice without a Net amount has nothing
// to compare the sum against, so that check is skipped.
//
// jsonField and label name the validated mapping in error messages,
// jsonField as written in the JSON and label in prose.
func (inv *Invoice) validateCostAmounts(amounts map[string]money.Amount, jsonField, label string) error {
	if len(amounts) == 0 {
		return nil
	}
	var sum money.Amount
	for number, amount := range amounts {
		if number == "" {
			return fmt.Errorf("empty %s string", jsonField)
		}
		if amount == 0 {
			return fmt.Errorf("%s '%s' amount must not be zero", label, number)
		}
		if amount < 0 {
			return fmt.Errorf("%s '%s' amount (%f) must not be negative", label, number, amount)
		}
		sum += amount
	}
	if inv.Net == nil {
		return nil
	}
	net := *inv.Net
	if inv.ConversionRate != nil {
		net = net.MultipliedByRate(*inv.ConversionRate)
	}
	if sum > net {
		return fmt.Errorf("sum of %s amounts %f greater than invoice net sum %f", label, sum, net)
	}
	return nil
}
