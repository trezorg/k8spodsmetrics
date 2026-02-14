for os in linux darwin windows; do
	for arch in amd64 arm64; do
		task build GOOS="${os}" GOARCH="${arch}"
	done
done
