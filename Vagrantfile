# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

  # Toran is the VM hosting the src code of this project!
  config.vm.define "toran" do |toran|
    toran.vm.box = "ubuntu/trusty64"
    toran.vm.provision "shell", inline: <<-SHELL
      apt-get update
      apt-get install -y libpcap-dev
    SHELL
    toran.vm.network "public_network"
    toran.vm.network "private_network", ip: "192.168.33.11", virtualbox__intnet: true
    toran.vm.synced_folder ".", "/srv/toran/"
  end

  # Client is a VM without external internet access using the toran VM
  # to get external access
  config.vm.define "client" do |client|
    client.vm.box = "ubuntu/trusty64"
    client.vm.network "private_network", ip: "192.168.33.10", virtualbox__intnet: true
  end
end
