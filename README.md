# 🧾 FPS QR Code Generator API (Hong Kong EMV Format)

This is a Golang-based REST API that generates **Hong Kong Faster Payment System (FPS)** QR Code strings based on the EMVCo standard.

Built using:
- ✅ [Gin](https://github.com/gin-gonic/gin) (lightweight web framework)
- ✅ Custom EMV TLV builder
- ✅ Optional nested `PaymentOperator` extensions (e.g. for Tap & Go)

---

## 🚀 Features

- Generate static and dynamic FPS QR Code strings
- Support for HK mobile numbers and email as FPS IDs
- Auto-detect ID type and format accordingly
- Currency support: `HKD`, `CNY` (mapped to ISO 4217)
- Field 32 extension (e.g., for proprietary tokens or labels)
- CRC16-CCITT validation using lookup table or algorithm
- JSON API with clean structure

---

## 📦 API Request Example

### `POST /api/v1/fps/qr`

**Request Body**:

```json
{
  "fps_id": "+852-91234567",
  "merchant_name": "My Shop",
  "city": "HK",
  "amount": "88.00",
  "currency": "HKD",
  "payment_operator": {
    "global_identifier": "hk.com.hkicl",
    "extra_fields": {
      "01": "txnRef123",
      "02": "sessionTokenXYZ"
    }
  }
}
