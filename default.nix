with import <nixpkgs> {};
stdenv.mkDerivation rec {
  name = "env";
  env = buildEnv {
    name = name;
    paths = buildInputs;
  };
  buildInputs = [
    pkgconfig
    portaudio
    alsaLib
  ];
}
