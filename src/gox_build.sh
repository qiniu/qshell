export GOPATH=$GOPATH:/Users/jemy/QiniuCloud/Projects/qshell
gox -os="darwin" -arch="386"  
gox -os="darwin" -arch="amd64"  
gox -os="windows" -arch="386" 
gox -os="windows" -arch="amd64"
gox -os="linux" -arch="arm"
mv src_darwin_386 ../bin/qshell_darwin_386
mv src_darwin_amd64 ../bin/qshell_darwin_amd64
mv src_windows_386.exe ../bin/qshell_windows_386.exe
mv src_windows_amd64.exe ../bin/qshell_windows_amd64.exe
