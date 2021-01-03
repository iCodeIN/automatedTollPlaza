#!/bin/bash
set -e

workdir=.cover
profile="$workdir/cover.out"
mode=count
packages=$(go list ./...)
error=31
minimumCoverage=80

output() {
	color="32"
	if [[ "$2" -gt 0 ]]; then
		color="31"
	fi
	printf "\033[${color}m"
	echo $1
	printf "\033[0m"
}

show_help() {
	output "Requires Go 1.9.2 or later."

	cat <<EOF
Generate test coverage statistics for le-api packages.

    -h | --help     Display the help
    --html          Additionally create HTML report
    -r | --report   Returns coverage report
EOF
	exit 0
}

generate_cover_data() {
	rm -rf "$workdir"
	mkdir "$workdir"

	for pkg in "$@"; do
		f="$workdir/$(echo $pkg | tr / -).cover"
		go test -covermode="$mode" -coverprofile="$f" "$pkg"
	done

	echo "mode: $mode" >"$profile"
	grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"
}

html_report() {
	go tool cover -html="$profile" -o="$workdir"/coverage.html
	open "$workdir"/coverage.html
	exit 0
}

report() {
	go tool cover -func="$profile"
	exit 0
}

case "$1" in
--help | -h)
	show_help
	;;
--html)
	html_report
	;;
--report | -r)
	report
	;;
*)
	# Go unit test with coverage
	generate_cover_data $packages
	percentage=$(go tool cover -func="$profile" |
		tail -n 1 |
		awk '{ print $3 }' |
		sed -e 's/^\([0-9]*\).*$/\1/g')

	if [ "$percentage" -ge 80 ]; then {
		output "Total (statements) covered ${percentage} is excellent test coverage"
		exit 0
	}
	fi
	echo "Total (statements) covered ${percentage}%"
	;;
esac