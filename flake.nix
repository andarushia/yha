{
  description = "postgres + go project";

  inputs = {
    nixpkgs     = { url = "github:nixos/nixpkgs/nixos-unstable"; };
    flake-utils = { url = "github:numtide/flake-utils"; };
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      inherit (pkgs.lib) optional optionals;
      pkgs = import nixpkgs { inherit system; };
      system = "x86_64-linux";
      postgresql = pkgs.postgresql_15;
    in  {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          git

          postgresql
          pgcli

          go
          gopls
        ];
        shellHook = ''
            if ! test -d .nix-shell; then
              mkdir .nix-shell
            fi

            export NIX_SHELL_DIR=$PWD/.nix-shell
            # Put the PostgreSQL databases in the project directory.
            export PGDATA=$NIX_SHELL_DIR/db
            '';
      };
    };
}
