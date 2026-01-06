# Hướng dẫn Biến Template Email

Tài liệu này mô tả các biến bạn có thể sử dụng trong nội dung Text hoặc HTML của Email Template để cá nhân hoá chiến dịch.

## Biến Thông Tin Người Nhận
- `{{.FirstName}}` – Tên
- `{{.LastName}}` – Họ
- `{{.Email}}` – Địa chỉ email
- `{{.Position}}` – Chức vụ (nếu có)
- `{{.RId}}` – Mã nhận dạng nội bộ (hiển thị dưới dạng giá trị tham số truy vấn `rid` trong URL)

## Biến URL & Theo Dõi
- `{{.URL}}` – Đường dẫn phishing đầy đủ (bao gồm tham số `rid` gắn với người nhận). Dùng cho nút hoặc thẻ `<a>`.
- `{{.TrackingURL}}` – Link phục vụ theo dõi mở email.
- `{{.Tracker}}` – Thẻ `<img>` ẩn tự động gửi request về server để ghi nhận mở mail.

## Biến QR Code
| Biến | Mô tả | Khi dùng |
|------|-------|---------|
| `{{.QRURL}}` | Thẻ `<img>` QR ở dạng **data URI base64** | Nhanh, không tạo attachment; có thể bị chặn bởi một số webmail. |
| `{{.QRURLCID}}` | Thẻ `<img>` QR sử dụng **CID attachment** (`qr.png`) | Tương thích tốt hơn (Outlook, Thunderbird, Apple Mail, nhiều webmail). |

### Lời khuyên
1. Bắt đầu với `{{.QRURLCID}}` nếu bạn gặp vấn đề QR không hiển thị.
2. Có thể đặt đoạn fallback: _"Nếu không thấy mã QR hãy bấm link: {{.URL}}"_.
3. Không cần tự tạo file `qr.png` trong Template – hệ thống tự đính kèm nếu thấy `cid:qr.png`.

## Biến Khác
- `{{.BaseURL}}` – Domain gốc (không path, không query)
- `{{.From}}` – Tên hiển thị (hoặc địa chỉ) người gửi

## Ví dụ HTML Tích Hợp
```html
<p>Xin chào {{.FirstName}},</p>
<p>Vui lòng xác nhận tại: <a href="{{.URL}}">Truy cập hệ thống</a></p>
<p>Hoặc quét mã QR bên dưới:</p>
{{.QRURLCID}}
<p>Nếu không thấy mã QR hãy dùng link ở trên.</p>
{{.Tracker}}
```

## Ghi chú Kỹ thuật
- Các biến được render bằng Go template engine. Bạn có thể kết hợp hoặc lồng các thẻ HTML bao quanh chúng.
- Không thêm giao thức trước `{{.URL}}` (ví dụ `https://{{.URL}}`) – giá trị đã là absolute URL.
- Tham số truy vấn nhận dạng (trước đây `rid`, hiện là `r_id`) đã thay đổi – mọi phân tích phải dùng `rid`.

## Xử lý Sự cố
| Vấn đề | Nguyên nhân | Giải pháp |
|--------|-------------|-----------|
| QR không hiển thị | Client chặn data URI | Dùng `{{.QRURLCID}}` |
| Tracker không ghi nhận | Người nhận chặn ảnh từ xa | Khuyến khích click / QR, theo dõi qua URL |
| Biến không render | Sai cú pháp (`{{ .FirstName }}` có khoảng thừa) | Dùng đúng: `{{.FirstName}}` |

## Mở rộng Tương Lai
- Tùy chọn kích thước QR (`128`, `256`) khi đủ nhu cầu.
- Template function như `{{ qr 128 }}` để điều chỉnh lỗi sửa và kích thước.

---
*Phiên bản tài liệu: 1.0*
