wget https://dl.min.io/server/minio/release/linux-amd64/minio -O /usr/local/bin/minio
wget https://dl.min.io/client/mc/release/linux-amd64/mc -O /usr/local/bin/mc

chmod +x /usr/local/bin/minio /usr/local/bin/mc

mkdir -p /mnt/data/minio

sudo bash -c 'echo "[Unit]
Description=MinIO
Documentation=https://min.io/docs/minio/linux/index.html
After=network.target

[Service]
User=minio-user
Group=minio-user
ExecStart=/usr/local/bin/minio server /mnt/data/minio --console-address ":9001"
ExecStartPost=/usr/local/bin/mc alias set local http://172.17.0.1:9000 admin password
ExecStartPost=/usr/local/bin/mc mb --region bj local/dcloud
ExecStartPost=/usr/local/bin/mc admin user svcacct add local admin --name upDownKey --access-key "yDD7FnMBNgh8jqX6ujlT" --secret-key "lnxju2Hjd88fJuqCRxZkdqVjCWi7UaiUvPMmYQ1x"
ExecStartPost=/usr/local/bin/mc anonymous set download local/dcloud
Restart=always
Environment=MINIO_ACCESS_KEY=my-access-key
Environment=MINIO_SECRET_KEY=my-secret-key
Environment=MINIO_ROOT_USER=admin
Environment=MINIO_ROOT_PASSWORD=password

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/minio.service'


useradd -r minio-user -s /sbin/nologin
chown minio-user:minio-user /mnt/data/minio

sudo systemctl daemon-reload
sudo systemctl start minio
sudo systemctl enable minio
