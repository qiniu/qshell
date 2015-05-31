export GOPATH=$GOPATH:/Users/jemy/QiniuCloud/Projects/qshell
gox -os="darwin" -arch="386"  
gox -os="darwin" -arch="amd64"  
gox -os="windows" -arch="386" 
gox -os="windows" -arch="amd64"  
gox -os="linux" -arch="386"  
gox -os="linux" -arch="amd64" 
mv src_* ../bin/
