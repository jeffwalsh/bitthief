# bitthief
Great way to understand SHA256, hashing, and how NOT to create randomness in code.

# Instructions
No dependencies. Simply run `./bin/bitthief` and it will start looping through every timestamp since the creation of Bitcoin, and the contents of`corncob_lowercase.txt`, as possible seed values. If you want it to read your own txt file instead of `corncob_lowercase.txt`, with no timestamp checking, run it as so:  
  
`./bin/bitthief -timestamp=false -file=file.txt`  
  
Your file first follow the line-delineated method that `corncob_lowercase.txt` does. EG:  
  
`aaron`  
`frank`  
`password`  
   

You can also pass a single word by running `./bin/bitthief -word=password`, and the program will exit after showing you the amount of BTC received and the current balance.