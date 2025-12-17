# Các tính năng tối thiểu cho ứng dụng Quản lý & Theo dõi Network/Domain (Custos)

Tài liệu này định nghĩa các tính năng cốt lõi (MVP) cho ứng dụng desktop theo dõi và quản lý dữ liệu mạng (network traffic monitoring & management).

## 1. Dashboard (Tổng quan Network)
- **Real-time Graph**: Biểu đồ dòng chảy dữ liệu (Upload/Download) theo thời gian thực.
- **Top Usage**: Hiển thị danh sách các Domain/Process tiêu thụ băng thông nhiều nhất trong phiên hiện tại.
- **Connection Status**: Số lượng kết nối đang mở (Active Connections).

## 2. Giám sát Traffic (Traffic Monitoring)
- **Log DNS Requests**:
    - Ghi lại tất cả các yêu cầu phân giải tên miền (DNS queries).
    - Hiển thị thông tin: Thời gian, Domain, Loại (A, AAAA, CNAME), IP phản hồi, Latency.
- **Log Connection**:
    - Hiển thị kết nối từ Process nào -> IP đích/Domain nào.
    - Protocol (TCP/UDP), Port đích.
    - Trạng thái kết nối (Established, Close_Wait, etc.).
- **Packet Inspection (Cơ bản)**: Hiển thị Header/Host của các kết nối HTTP/HTTPS (nếu giải mã được hoặc qua proxy).

## 3. Kiểm tra & Theo dõi (Tracking & Health Check)
- **Health Check tự động**:
    - Định kỳ kiểm tra trạng thái sống/chết của proxy (Cronjob).
    - Kiểm tra tốc độ phản hồi (Latency/Ping).
    - Tự động đánh dấu "Dead" nếu không phản hồi sau X lần thử.
- **Kiểm tra ẩn danh (Anonymity Check)**: Xác định xem proxy có bị lộ IP thật không (Transparent, Anonymous, Elite).
- **Vị trí địa lý (Geo-location)**: Hiển thị Quốc gia/Thành phố của IP Proxy.

## 4. Quản lý & Chặn (Filtering & Rules)
- **Blocklist/Allowlist**:
    - Chặn kết nối tới các Domain/IP cụ thể (tương tự file hosts hoặc firewall rules).
    - Hỗ trợ Wildcard chặn domain (ví dụ: `*.ads.com`).
- **Quy tắc theo Process**: Cho phép hoặc chặn internet của một ứng dụng cụ thể trên máy.
- **Tùy chỉnh DNS**: Ép buộc sử dụng DNS Server cụ thể (Global hoặc theo quy tắc) để vượt qua chặn hoặc tăng tốc độ.

## 5. Quản lý Proxy & Routing (Chi tiết hóa yêu cầu trước)
- **Route Traffic**: Cấu hình quy tắc để đẩy traffic của Domain X hoặc App Y đi qua một Proxy cụ thể.
- **Proxy Chain**: Quản lý các upstream proxy (SOCKS5/HTTP) để ẩn danh traffic đi ra.

## 6. Báo cáo & Lịch sử (Reports & History)
- **Lưu trữ lịch sử**: Xem lại lịch sử truy cập trong quá khứ (theo ngày/tuần).
- **Thống kê dữ liệu**: Tổng lượng data đã dùng (Total Upload/Download) theo từng Domain.

## 7. Settings (Hệ thống)
- **Chế độ Capture**: Bật/Tắt chế độ bắt gói tin (Data Capture Mode).
- **Net Interface Selection**: Chọn card mạng để lắng nghe (Wi-Fi, Ethernet).
