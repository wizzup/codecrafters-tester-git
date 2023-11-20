{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ../flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
        ];
      }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
}:
buildGoApplication {
  pname = "tester-git";
  version = "0.1";
  pwd = ./..;
  src = ./..;
  modules = ../gomod2nix.toml;
  nativeBuildInputs = with pkgs; [python3 git];
  buildPhase = ''
    go build -o tester-git ./cmd/tester
  '';
  checkPhase = ''
    go test -v ./internal/
  '';
  installPhase = ''
    mkdir -p $out/bin
    cp tester-git $out/bin
  '';
}
