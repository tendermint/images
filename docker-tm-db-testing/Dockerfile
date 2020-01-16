FROM circleci/golang

RUN sudo apt-get install -y --no-install-recommends libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev libzstd-dev

# Install cleveldb
RUN wget https://github.com/google/leveldb/archive/v1.20.tar.gz && tar xvf v1.20.tar.gz && rm -f v1.22.tar.gz
RUN cd leveldb-1.20 && make && sudo scp -r out-static/lib* out-shared/lib* "/usr/local/lib" && cd include && sudo scp -r leveldb /usr/local/include && sudo ldconfig && cd ~

# Install Rocksdb
RUN git clone https://github.com/facebook/rocksdb.git && cd rocksdb/ && sudo make install

# Install boltdb
RUN go get go.etcd.io/bbolt/...
