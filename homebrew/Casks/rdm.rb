cask "rdm" do
  name "Remote Development Manager"
  desc "A tool for remote development environments"
  homepage "https://github.com/BlakeWilliams/remote-development-manager"
  version "0.0.6"

  file_to_move = if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/BlakeWilliams/remote-development-manager/releases/download/v#{version}/rdm-darwin-arm64"
      sha256 "c562d6040a2d84e60790f7de7a4bc7e4d9bdad390cc72cc0d402c7eb6e9553b2"
      "rdm-darwin-arm64"
    else
      url "https://github.com/BlakeWilliams/remote-development-manager/releases/download/v#{version}/rdm-darwin-amd64"
      sha256 "617d002120fdfe227aed377a998334ddbfc418758a03b89baa9027ea9f976429"
      "rdm-darwin-amd64"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/BlakeWilliams/remote-development-manager/releases/download/v#{version}/rdm-linux-arm64"
      sha256 "fb42eacfe2ec272d66660569524ad4b732311727f31cb72d9ef8ea36bb852941"
      "rdm-linux-arm64"
    else
      url "https://github.com/BlakeWilliams/remote-development-manager/releases/download/v#{version}/rdm-linux-amd64"
      sha256 "9b79290ef87e0e0f37e71cf9a76ef1a4377472c56907fd01241f9881ecc57d36"
      "rdm-linux-amd64"
    end
  end

  binary file_to_move, target: "rdm"
end
