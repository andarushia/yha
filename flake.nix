{
  description = "postgres + go project";

  inputs = {
    nixpkgs     = { url = "github:nixos/nixpkgs/nixos-unstable"; };
  };

  outputs = { self, nixpkgs }:
    let
      pkgs = import nixpkgs { inherit system; };
      system = "x86_64-linux";
      go = pkgs.go_1_21;

      podmanSetupScript = let
      registriesConf = pkgs.writeText "registries.conf" ''
        [registries.search]
        registries = ['docker.io']
        [registries.block]
        registries = []
      '';
      in pkgs.writeScript "podman-setup" ''
        #!${pkgs.runtimeShell}
        # Dont overwrite customised configuration
        if ! test -f ~/.config/containers/policy.json; then
          install -Dm555 ${pkgs.skopeo.src}/default-policy.json ~/.config/containers/policy.json
        fi
        if ! test -f ~/.config/containers/registries.conf; then
          install -Dm555 ${registriesConf} ~/.config/containers/registries.conf
        fi
      '';

      # Provides a fake "docker" binary mapping to podman
      dockerCompat = pkgs.runCommandNoCC "docker-podman-compat" {} ''
        mkdir -p $out/bin
        ln -s ${pkgs.podman}/bin/podman $out/bin/docker
      '';
    in  {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = [
          pkgs.pgcli
          
          go
          pkgs.gopls

          dockerCompat
          pkgs.podman  # Docker compat
          pkgs.runc  # Container runtime
          pkgs.conmon  # Container runtime monitor
          pkgs.skopeo  # Interact with container registry
          pkgs.slirp4netns  # User-mode networking for unprivileged namespaces
          pkgs.fuse-overlayfs  # CoW for images, much faster than default vfs
        ];
        shellHook = ''
          if ! test -d .nix-shell; then
            mkdir .nix-shell
          fi

          export NIX_SHELL_DIR=$PWD/.nix-shell
          # Put the PostgreSQL databases in the project directory.
          export PGDATA=$NIX_SHELL_DIR/db

          ${podmanSetupScript}
        '';
      };
    };
}
