# Install compiler 
```
sudo apt install llvm-12 llvm-12-dev
git clone git@github.com:RainingComputers/ShnooTalk.git
cd ShnooTalk
make build
sudo make install -j 8
sudo codesign -s - /usr/local/bin/shtkc 
cd ..
rm -rf ShnooTalk
```

# Compile
```
make compile
```

# Run benchmark
```
time ./search
```

# Result
```
The best witness is: 10080
./search  44.86s user 0.20s system 98% cpu 45.552 total
```
