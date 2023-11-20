{
  description = "hello-gol";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  };
  nixConfig.bash-prompt-suffix = "🔨";
  outputs = {
    self,
    nixpkgs,
    gomod2nix,
  }: (let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
    drv = pkgs.callPackage ./nix {
      inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
    };
    dev = pkgs.callPackage ./nix/shell.nix {
      inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
    };
  in {
    formatter.${system} = pkgs.alejandra;
    packages.${system}.default = drv;
    devShells.${system}.default = dev;
  });
}
