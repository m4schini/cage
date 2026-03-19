{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  packages = with pkgs; [
    zsh
    go
    claude-code
  ];

  shellHook = ''
    exec ${pkgs.zsh}/bin/zsh
  '';
}
