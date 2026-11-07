#!/bin/sh
set -e

if [ "$1" = "configure" ]; then
	getent group plugdev >/dev/null || groupadd --system plugdev

	if ! getent passwd rtldavis >/dev/null; then
		useradd --system --user-group --groups plugdev \
			--home-dir /nonexistent --no-create-home \
			--shell /usr/sbin/nologin rtldavis
	fi

	if [ -d /run/systemd/system ]; then
		systemctl daemon-reload
		# Pick up the new binary on upgrade; a no-op if the service isn't running
		systemctl try-restart rtldavis.service
	fi

	# $2 is the previously configured version, empty on first install
	if [ -z "$2" ]; then
		cat <<-EOF

		rtldavis installed. To finish setup:
		  1. Edit /etc/default/rtldavis with your server URL and API key
		  2. Reboot (or run 'sudo modprobe -r dvb_usb_rtl28xxu') so the
		     DVB-T driver releases the dongle
		  3. sudo systemctl enable --now rtldavis

		EOF
	fi
fi
