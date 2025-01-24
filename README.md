# Go Workflow Management Example

Bu proje, iş akışı yönetim sisteminin örnek bir implementasyonudur.

## Örnek Kullanım

### 1. Sipariş Workflow'unu Oluşturma

```bash
curl -X POST http://localhost:8080/workflows \
-H "Content-Type: application/json" \
-d '{
  "name": "Sipariş İşlemi",
  "type": "order_process",
  "created_by": "65b012345678901234567890",
  "steps": [
    {
      "id": "65b012345678901234567891",
      "type": "task",
      "title": "Sipariş Detayları",
      "assigned_to": "65b012345678901234567890",
      "status": "pending",
      "next_steps": ["65b012345678901234567892"]
    },
    {
      "id": "65b012345678901234567892",
      "type": "approval",
      "title": "Stok Kontrolü",
      "assigned_to": "65b012345678901234567893",
      "status": "pending",
      "next_steps": ["65b012345678901234567894"],
      "required_data": ["order_items", "total_amount"]
    },
    {
      "id": "65b012345678901234567894",
      "type": "process",
      "title": "Fatura Oluşturma",
      "assigned_to": "65b012345678901234567895",
      "status": "pending",
      "result_type": "invoice",
      "required_data": ["order_items", "customer_info", "total_amount", "stock_approval"]
    }
  ]
}'
```

### 2. Sipariş Detaylarını Girme

```bash
curl -X POST http://localhost:8080/workflows/WORKFLOW_ID/steps/65b012345678901234567891/process \
-H "Content-Type: application/json" \
-d '{
  "action": "approve",
  "data": {
    "order_items": [
      {
        "product_id": "PROD001",
        "name": "Laptop",
        "quantity": 1,
        "price": 15000
      }
    ],
    "customer_info": {
      "name": "Ahmet Yılmaz",
      "email": "ahmet@example.com",
      "tax_number": "1234567890"
    },
    "total_amount": 15000
  }
}'
```

### 3. Stok Kontrolü

```bash
curl -X POST http://localhost:8080/workflows/WORKFLOW_ID/steps/65b012345678901234567892/process \
-H "Content-Type: application/json" \
-d '{
  "action": "approve",
  "data": {
    "stock_approval": true,
    "stock_notes": "Stok yeterli",
    "approved_by": "Depo Sorumlusu"
  }
}'
```

### 4. Fatura Oluşturma

```bash
curl -X POST http://localhost:8080/workflows/WORKFLOW_ID/steps/65b012345678901234567894/process \
-H "Content-Type: application/json" \
-d '{
  "action": "approve",
  "data": {
    "invoice_number": "FTR-2024-001",
    "invoice_date": "2024-01-24T15:00:00Z",
    "items": [
      {
        "product_id": "PROD001",
        "name": "Laptop",
        "quantity": 1,
        "unit_price": 15000,
        "total": 15000
      }
    ],
    "subtotal": 15000,
    "tax": 2700,
    "total": 17700
  }
}'
```
