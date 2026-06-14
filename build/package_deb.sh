#!/bin/bash
set -e

# Thư mục làm việc hiện tại động theo vị trí script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$( dirname "$SCRIPT_DIR" )"
BUILD_DIR="$PROJECT_DIR/build"
DEB_STAGE_DIR="$BUILD_DIR/debian"

# Nhận version từ tham số đầu tiên, mặc định là 1.0
VERSION="${1:-1.0}"
# Bỏ chữ 'v' ở đầu nếu có (ví dụ 'v1.1' -> '1.1')
VERSION="${VERSION#v}"

echo "=== ĐANG CHUẨN BỊ ĐÓNG GÓI DEBIAN PACKAGE (.DEB) VERSION $VERSION ==="

# 1. Dọn dẹp thư mục stagging cũ nếu có
rm -rf "$DEB_STAGE_DIR"

# 2. Tạo cấu trúc thư mục tiêu chuẩn
mkdir -p "$DEB_STAGE_DIR/DEBIAN"
mkdir -p "$DEB_STAGE_DIR/usr/bin"
mkdir -p "$DEB_STAGE_DIR/usr/share/applications"
mkdir -p "$DEB_STAGE_DIR/usr/share/icons/hicolor/512x512/apps"

# 3. Tạo tệp DEBIAN/control
cat << EOF > "$DEB_STAGE_DIR/DEBIAN/control"
Package: google-model-tester
Version: $VERSION
Section: utils
Priority: optional
Architecture: amd64
Maintainer: Nguyen Duy Truong <skul9x@gmail.com>
Description: Trinh kiem thu va so sanh hieu nang cac mo hinh Google Gemini LLMs voi thiet ke Glassmorphism dep mat.
EOF

# 4. Tạo tệp desktop shortcut
cat << 'EOF' > "$DEB_STAGE_DIR/usr/share/applications/google-model-tester.desktop"
[Desktop Entry]
Name=Google Model Tester
Comment=So sánh và kiểm thử hiệu năng tóm tắt các mô hình Gemini LLMs
Exec=google-model-tester
Icon=google-model-tester
Type=Application
Terminal=false
Categories=Utility;Development;
EOF

# 5. Sao chép tệp thực thi chính và icon ứng dụng
cp "$BUILD_DIR/bin/google-model-tester" "$DEB_STAGE_DIR/usr/bin/google-model-tester"
cp "$BUILD_DIR/appicon.png" "$DEB_STAGE_DIR/usr/share/icons/hicolor/512x512/apps/google-model-tester.png"

# 6. Phân quyền chuẩn
chmod 755 "$DEB_STAGE_DIR/usr/bin/google-model-tester"
chmod 644 "$DEB_STAGE_DIR/usr/share/applications/google-model-tester.desktop"
chmod 644 "$DEB_STAGE_DIR/usr/share/icons/hicolor/512x512/apps/google-model-tester.png"

# 7. Đóng gói bằng dpkg-deb
echo "Đóng gói ứng dụng thành file .deb..."
dpkg-deb --build "$DEB_STAGE_DIR" "$BUILD_DIR/bin/google-model-tester_${VERSION}_amd64.deb"

echo "=== ĐÃ ĐÓNG GÓI THÀNH CÔNG: build/bin/google-model-tester_${VERSION}_amd64.deb ==="

# 8. Dọn dẹp thư mục stagging tạm thời
rm -rf "$DEB_STAGE_DIR"
