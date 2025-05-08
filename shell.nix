{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.gcc
    pkgs.glibc.dev
  ];

  shellHook = ''
    export CGO_ENABLED=1
  '';
}
