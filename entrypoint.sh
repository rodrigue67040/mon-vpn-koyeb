#!/bin/sh
PORT=${PORT:-8080}
UUID=${UUID:-a7543913-8637-4753-be7e-1c4e5639996d}
WS_PATH=${WS_PATH:-/test}
cat > /app/config.json <<EOF
{
"log": { "loglevel": "warning"},
"inbounds": [
{
"listen": "0.0.0.0",
"port": ${PORT},
"protocol": "vless",
"settings": {
"clients": [ { "id": "${UUID}", "flow": ""}],
"decryption": "none"
},
"streamSettings": {
"network": "ws",
"wsSettings": { "path": "${WS_PATH}"}
},
"sniffing": { "enabled": true, "destOverride": ["http", "tls"]}
}
],
"outbounds": [ { "protocol": "freedom"}]
}
EOF
exec /app/xray run -c /app/config.json
