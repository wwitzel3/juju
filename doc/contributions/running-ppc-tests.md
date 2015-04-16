Running the Juju tests under ppc64el
------------------------------------

You will want to make sure you have the architecrue extensions for qemu.

    sudo apt-get install qemu-system

You will also want to download an Ubuntu ppc64el image.

    wget http://cdimage.ubuntu.com/releases/trusty/release/ubuntu-14.04.1-server-ppc64el.iso

I used the virt-manager GUI tool to create my VM. You will want to make sure you go in to the
Advanced Configuration options and select qemu - ppc64 as the architecture.
