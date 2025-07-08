#!/bin/bash
cd ./go-ui
go build -o ../config/extras/pcd-iso-ui .
cd ..
sudo rm -rf iso_root.* ubuntu.* ubuntu-custom*
qemu-img create -f qcow2 ubuntu.qcow2 45G
make init
make setup-go-ui
make geniso

qemu-system-x86_64 \
  -m 12288 \
  -smp 3 \
  -boot d \
  -cdrom ./ubuntu-custom* \
  -drive file=ubuntu.qcow2,format=qcow2 \
  -enable-kvm \
  -cpu host,+vmx \
  -net nic -net user \
  -display gtk,grab-on-hover=on \
  -vga virtio \
  -usb \
  -device usb-tablet

 qemu-system-x86_64 \
  -m 12288 \
  -smp 2 \
  -drive file=ubuntu.qcow2,format=qcow2 \
  -enable-kvm \
  -cpu host,+vmx \
  -net nic -net user \
  -display gtk,grab-on-hover=on \
  -vga virtio \
  -usb \
  -device usb-tablet
