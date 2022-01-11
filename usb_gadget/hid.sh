#!/bin/bash -eu

modprobe libcomposite

cd /sys/kernel/config/usb_gadget
rm -rf g1
mkdir -p ./g1
cd ./g1

echo 0x1d6b > idVendor # Linux Foundation
echo 0x0104 > idProduct # Multifunction Composite Gadget
echo 0x0100 > bcdDevice # v1.0.0
echo 0x0200 > bcdUSB # USB2
mkdir -p strings/0x409

echo "HID1234567890" > strings/0x409/serialnumber
echo "Linux" > strings/0x409/manufacturer
echo "Linux USB Gadget/HID" > strings/0x409/product

USBN="usb0"
CONF=2
rm -rf configs/c.${CONF}
mkdir -p configs/c.${CONF}/strings/0x409
echo 250 > configs/c.${CONF}/MaxPower 
echo "Config ${CONF}: ECM network" > configs/c.${CONF}/strings/0x409/configuration 

rm -rf functions/hid.${USBN}
mkdir -p functions/hid.${USBN}
echo 1 > functions/hid.${USBN}/protocol
echo 1 > functions/hid.${USBN}/subclass
echo 8 > functions/hid.${USBN}/report_length
echo -ne \\x05\\x01\\x09\\x06\\xa1\\x01\\x05\\x07\\x19\\xe0\\x29\\xe7\\x15\\x00\\x25\\x01\\x75\\x01\\x95\\x08\\x81\\x02\\x95\\x01\\x75\\x08\\x81\\x01\\x95\\x05\\x75\\x01\\x05\\x08\\x19\\x01\\x29\\x05\\x91\\x02\\x95\\x01\\x75\\x03\\x91\\x01\\x95\\x06\\x75\\x08\\x15\\x00\\x25\\x65\\x05\\x07\\x19\\x00\\x29\\x65\\x81\\x00\\xc0 > functions/hid.${USBN}/report_desc
ln -s functions/hid.${USBN} configs/c.${CONF}/

ls /sys/class/udc > UDC
