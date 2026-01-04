{
  description = "HomeBox External Label Service";

  inputs = {
    naersk.url = "github:nix-community/naersk/master";
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils, naersk }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        name = "homebox-label-service";

        goBuild = pkgs.buildGoModule {
          pname = name;
          version = "0.1.0";
          src = ./.;
          subPackages = [ "." ];
          vendorHash = "sha256-XWKIxTuHjHwiQ2QsLPt3FUI3ovGxImMQbzCXiB7522k=";
        };

        # The actual binary name (Go uses directory/module name)
        binaryName = "HomeBox-External-Label-Service";

        # Create a package that has the binary in /bin
        appPackage = pkgs.runCommand "${name}-app" {} ''
          mkdir -p $out/bin
          cp ${goBuild}/bin/${binaryName} $out/bin/${binaryName}
          chmod +x $out/bin/${binaryName}
        '';

        dockerImage = pkgs.dockerTools.buildImage {
          name = name;
          tag = "latest";
          copyToRoot = pkgs.buildEnv {
            name = "image-root";
            paths = [
              pkgs.cacert
              appPackage
              pkgs.bash
            ];
            pathsToLink = [ "/bin" "/etc" ];
          };
          config = {
            Entrypoint = [ "/bin/${binaryName}" ];
            ExposedPorts = { "8080/tcp" = { }; };
            WorkingDir = "/";
            Env = [
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            ];
          };
        };
      in
      {
        defaultPackage = goBuild;
        packages = {
          default = goBuild;
          homebox-label-service = goBuild;
          dockerImage = dockerImage;
        };

        defaultApp = {
          type = "app";
          program = "${goBuild}/bin/${name}";
        };
      });
}
