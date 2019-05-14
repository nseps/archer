# archer

Archiving simplified

# Archive Formats

* tar
* cpio
* zip
* rar

# Compression Formats

* gz
* bz2
* xz
* lz4

# File Types

- [x] Directory
- [x] Regular
- [x] Symlink
- [ ] Block Device
- [ ] Char Device
- [ ] Socket
- [ ] Named Pipe (FIFO)

# Combinations

You can combine any archive type (except rar) with any compression. That means 12 compressed and 3 uncompressed output file types.

# What can it do

Only 2 things

* Pack a directory to an archive
* Unpack an archive to a directory

run `archer -h` for more

The main focus of this tools is handling tar and cpio archives as it was born to handle container images and initramfs archives. The file type support will be implemented when needed.

# TODO

* List
* Verify 

# Installing & Configuring

`go get -u github.com/thegrumpylion/archer`

## Generate completion file

`archer -f archer_complete`

You can copy it **/etc/bash_completion.d/** if an option. Otherwise you can just source it.

`. archer_complete`

# Contributing & extending

Contributions are welcome. Just file a PR.

## Hacking Guide

Later. The basic interfaces might change. You can always take a look in the code. Is very minimal.

# Notes on zip & rar

* zip

zip is used as packaging container like tar. No deflate is used to compress any files before packing. You can create archives like name.zip.xz for example and unpack them just fine with this tool but i've never seen zip being used like that and i don't use it anyway.

* rar

Well, it was quick and easy to implement. I haven't tested if works and i did not implement symlink support. Unpack only.

# Limitations

* Cannot handle packing individual files for now. Source can only be a directory.
* No way to pass custom options to archivers and compressors. This keeps the user interface lean but limited.
* No option to output to stdout though seems like a good idea
