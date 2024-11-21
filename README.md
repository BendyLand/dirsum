# Dirsum

A simple CLI tool that recursively counts the number of each file type in a given directory and prints the results!

## Usage

```bash
dirsum <path> <opt args>
dirsum .        # Will output a count of each file type, sorted by name.
dirsum . -n     # This one will be sorted descending by number.
dirsum . -nr    # This will reverse it.
dirsum . -nrt   # This will include the total number of files.
dirsum . -v     # This is equal to -nt.

dirsum . -h     
dirsum . --help # The help args override the rest and will only print the help menu.
```
