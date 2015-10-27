export GOPATH=$GOPATH:/Users/jemy/QiniuCloud/Projects/qshell
gox -os="darwin" -arch="386"  -output="../bin/qshell_darwin_386"
gox -os="darwin" -arch="amd64" -output="../bin/qshell_darwin_amd64" 
gox -os="windows" -arch="386" -output="../bin/qshell_windows_386"
gox -os="windows" -arch="amd64" -output="../bin/qshell_windows_amd64"
