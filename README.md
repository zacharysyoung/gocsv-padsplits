# Pad numbered files

Responding to an issue in GoCSV, [sortable split output files](https://github.com/aotimme/gocsv/issues/62):

> When using gocsv's split, the output file names have suffix such as -1.csv, etc. so if there are more than 9 the filenames don't sort properly. Would be nice to be able to optionally zero pad the number in the suffix with sufficient zeroes that they sort properly.

Given the 11-row input CSV, [input.csv](./input.csv):

```none
I
a
b
c
d
e
f
g
h
i
j
k
```

the following gocsv split command will produce 11 files, named input-1.csv to input-11.csv, and those files will usually sort lexically, like:

```sh
% gocsv split -max-rows=1 input.csv
% ls input-*.csv | sort
input-1.csv
input-10.csv
input-11.csv
input-2.csv
...
```

and not numerically, like:

```none
input-1.csv
input-2.csv
...
input-10.csv
input-11.csv
```

Use gocsv-padsplits to pad those files, like:

```sh
% gocsv-padsplits 'input-'
% ls input-*.csv | sort
input-01.csv
input-02.csv
...
input-10.csv
input-11.csv
```

## Installing

Run go's install subcommand to build and install gocsv-padsplits to GOPATH/bin:

```sh
% go install
% which gocsv-padsplits
/***/***/go/bin/gocsv-padsplits
```
