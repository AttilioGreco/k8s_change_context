{
  description = "Kubernetes context switcher CLI tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        k8s-change-context = pkgs.buildGoModule {
          pname = "k8s-change-context";
          version = "0.1.0";

          src = ./.;

          # Uses the vendor/ directory — run `go mod vendor` after updating deps.
          vendorHash = null;

          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "Switch Kubernetes contexts from the command line";
            homepage = "https://github.com/AttilioGreco/k8s_change_context";
            license = licenses.mit;
            mainProgram = "k8s-change-context";
          };
        };
      in
      {
        packages = {
          default = k8s-change-context;
          k8s-change-context = k8s-change-context;
        };

        apps.default = flake-utils.lib.mkApp {
          drv = k8s-change-context;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            goreleaser
            gopls
          ];
        };
      }
    );
}
