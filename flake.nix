{
  description = "lech-linter - spell it right";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable-small";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
    nur.url = "github:nix-community/NUR";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gomod2nix,
    nur,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [gomod2nix.overlays.default nur.overlays.default];
          config.allowUnfreePredicate = pkg: builtins.elem (pkgs.lib.getName pkg) ["goreleaser-pro"];
        };
        lib = pkgs.lib;
        version =
          if (self ? shortRev)
          then self.shortRev
          else "dev";
      in {
        packages.default = pkgs.buildGoApplication {
          pname = "lech-linter";
          inherit version;

          src = lib.fileset.toSource {
            root = ./.;
            fileset = lib.fileset.unions [
              ./go.mod
              ./go.sum
              ./gomod2nix.toml
              (lib.fileset.fileFilter (file: lib.hasSuffix ".go" file.name) ./.)
            ];
          };

          modules = ./gomod2nix.toml;

          # Only build the main binary
          subPackages = ["."];

          ldflags = [
            "-s"
            "-w"
            "-X github.com/asphaltbuffet/lech-linter/internal/version.Version=${version}"
          ];

          postInstall = ''
            # Generate and install shell completions
            installShellCompletion --cmd lech-linter \
              --bash <($out/bin/lech-linter completion bash) \
              --fish <($out/bin/lech-linter completion fish) \
              --zsh <($out/bin/lech-linter completion zsh)

            # Generate and install man pages
            mkdir -p $out/share/man/man1
            $out/bin/lech-linter man | gzip -c > $out/share/man/man1/lech-linter.1.gz
          '';

          nativeBuildInputs = [pkgs.installShellFiles];

          meta = with lib; {
            description = "spell it right";
            homepage = "https://github.com/asphaltbuffet/lech-linter";
            license = licenses.mit;
            mainProgram = "lech-linter";
          };
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            jujutsu
            jjui
            mise
            vhs
            ripgrep
            fd
            sd
            imagemagick
            gopls
            nixd
            pkgs.nur.repos.goreleaser.goreleaser-pro
            gomod2nix.packages.${system}.default
          ];

          shellHook = ''
            mise trust --all
          '';

          CGO_ENABLED = "1";
        };
      }
    )
    // {
      overlays.default = final: prev: {
        lech-linter = self.packages.${prev.system}.default;
      };
    };
}
