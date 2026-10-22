for os in linux darwin windows; do
	for arch in amd64 arm64; do
		GOOS="${os}" GOARCH="${arch}" mise run build
	done
done
