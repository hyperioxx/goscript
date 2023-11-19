echo "Creating modules directory..."
sudo mkdir -p /usr/local/lib/goscript/modules
sudo groupadd goscript
sudo chown -R root:goscript /usr/local/lib/goscript
sudo chmod 774 -R /usr/local/lib/goscript


sudo usermod -aG goscript $USER

echo "Setting GS_MODULES environment variable..."
if ! grep -q 'export GS_MODULES' ~/.bashrc ; then
    echo "export GS_MODULES=/usr/local/lib/goscript/modules" >> ~/.bashrc
fi
source ~/.bashrc

echo "Installation completed."
echo "Please restart your terminal for the changes to take effect."
