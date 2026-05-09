{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  packages = with pkgs; [
    fish
    go
    claude-code
  ];

  shellHook = ''
    exec ${pkgs.fish}/bin/fish
  '';
}
