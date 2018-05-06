# crawler
Web crawler in Go.  Mostly an experiment to learn the language, 
especially its concurrency approach.

Web crawlers are common.  This one has a couple of command line options:

`-maxdepth`: integer - the maximum depth of recursion, taking the provided (set of) URLs as level one
`-domainsticky`: boolean - if true will not follow links to domains external from the starting URL provided on the command line

Full format of the command line:

command \[options\] \[urls\]


