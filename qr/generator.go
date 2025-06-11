package qr

import (
	"fmt"
	"log"
	"net/mail"
	"regexp"
	"strings"
)

func formatEMV(id string, value string) string {
	return fmt.Sprintf("%s%02d%s", id, len(value), value)
}

func join(parts []string) string {
	return strings.Join(parts, "")
}

func getCurrencyCode(iso string) string {
	switch strings.ToUpper(iso) {
	case "", "HKD":
		return "344"
	case "CNY":
		return "156"
	case "USD":
		return "840"
	case "EUR":
		return "978"
	case "JPY":
		return "392"
	case "GBP":
		return "826"
	case "SGD":
		return "702"
	default:
		return "344" // fallback to HKD
	}
}

// +852-XXXXXXXX
var hkMobileWithPrefixRegex = regexp.MustCompile(`^\+852\-[5689]\d{7}$`)

func detectFPSIdentifierType(input string) string {
	// 1. Check if valid email
	if _, err := mail.ParseAddress(input); err == nil && strings.Contains(input, "@") {
		return "04" // valid email
	}

	// 2. Check if valid +852-xxxxxxxx mobile
	if hkMobileWithPrefixRegex.MatchString(input) {
		return "03" // valid HK mobile with +852-
	}

	// 3. Unknown or unsupported format
	return "02"
}

func GenerateQRString(req QRRequest) string {
	var payload []string

	// QR Code Conventions
	payload = append(payload, formatEMV("00", "01"))
	if req.Dynamic {
		payload = append(payload, formatEMV("01", "12"))
	} else {
		payload = append(payload, formatEMV("01", "11"))
	}

	// Merchant Account Info (ID 26 for FPS)
	fpsAccount := []string{
		formatEMV("00", "hk.com.hkicl"),
		formatEMV(detectFPSIdentifierType(req.FPSID), req.FPSID),
	}
	if req.MerchantTimeout != "" {
		fpsAccount = append(fpsAccount, formatEMV("05", req.MerchantTimeout))
	}
	payload = append(payload, formatEMV("26", join(fpsAccount)))

	// Additional Merchant Info
	payload = append(payload, formatEMV("52", "0000"))

	payload = append(payload, formatEMV("58", "HK"))
	payload = append(payload, formatEMV("59", req.MerchantName))
	payload = append(payload, formatEMV("60", req.City))
	// Optional Transaction Value
	if req.Amount != "" {
		currencyCode := getCurrencyCode(req.Currency)
		payload = append(payload, formatEMV("53", currencyCode))
		payload = append(payload, formatEMV("54", req.Amount))
	}
	// Additional Data Field (ID 62)
	var additionalFields []string
	if req.BillNumber != "" {
		additionalFields = append(additionalFields, formatEMV("01", req.BillNumber)) // 01 = "Bill Number" (closest spec fit)
	}
	if req.ReferenceLabel != "" {
		additionalFields = append(additionalFields, formatEMV("05", req.ReferenceLabel)) // 05 = "Reference Label" (closest spec fit)
	}
	if len(additionalFields) > 0 {
		payload = append(payload, formatEMV("62", join(additionalFields))) // ID 62 for Additional Data Field
	}

	// Payment Operator Info
	// Use ID 32 for Payment Operator Info
	var paymentOperators []string

	if req.PaymentOperator.GlobalUniqueIdentifier != "" {
		paymentOperators = append(paymentOperators, formatEMV("00", req.PaymentOperator.GlobalUniqueIdentifier))
	}
	if req.PaymentOperator.ExtraFields != nil {
		for key, value := range req.PaymentOperator.ExtraFields {
			if len(value) > 0 {
				paymentOperators = append(paymentOperators, formatEMV(key, value))
			}
		}
	}
	if len(paymentOperators) > 0 {
		payload = append(payload, formatEMV("32", join(paymentOperators)))
	}
	payload = append(payload, "6304") // CRC16 placeholder, will be replaced later
	raw := join(payload)
	crc, err := CRC16(raw)
	if err != nil {
		log.Printf("CRC16 calculation failed: %v", err)
		return ""
	}
	raw += fmt.Sprintf("%04X", crc)

	return raw
}
