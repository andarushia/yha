{
  description = "postgres + go project";

  inputs = {
    nixpkgs = { url = "github:nixos/nixpkgs/nixos-unstable"; };
  };

  outputs = { self, nixpkgs }:
  let
      pkgs = import nixpkgs { inherit system; };
      system = "x86_64-linux";
      go = pkgs.go_1_21;
    in  {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = [
          go
          pkgs.gopls
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
      services.postgres = {
        service.image = "postgres:16";
        service.volumes = [ "${toString ./.}/postgres-data:/var/lib/postgresql/data" ];
        service.environment.POSTGRES_PASSWORD = "sqlxpass";
        service.environment.POSTGRES_USER = "anya";
        service.environment.POSTGRES_DB = "anyadb";
      };
    };
}
